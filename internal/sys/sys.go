// Package sys reads Linux procfs and sysfs for the monitoring commands.
// It does not execute existing system utilities.
package sys

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const userHZ = 100

// Die prints err to stderr and exits. Command mains use it at the process boundary.
func Die(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Base(os.Args[0]), err)
	os.Exit(1)
}

// Usage prints a usage line and exits 2.
func Usage(args string) {
	fmt.Fprintf(os.Stderr, "usage: %s %s\n", filepath.Base(os.Args[0]), args)
	os.Exit(2)
}

// Human formats a byte count in IEC units.
func Human(n uint64) string {
	if n < 1024 {
		return fmt.Sprintf("%d B", n)
	}
	units := []string{"KiB", "MiB", "GiB", "TiB", "PiB"}
	v := float64(n)
	i := -1
	for v >= 1024 && i < len(units)-1 {
		v /= 1024
		i++
	}
	return fmt.Sprintf("%.1f %s", v, units[i])
}

// Read returns the contents of path.
func Read(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", path, err)
	}
	return string(b), nil
}

// Meminfo holds selected /proc/meminfo values in bytes.
type Meminfo struct {
	Total     uint64
	Free      uint64
	Available uint64
	Buffers   uint64
	Cached    uint64
	SwapTotal uint64
	SwapFree  uint64
}

// ParseMeminfo parses /proc/meminfo text. Values are converted from kB to bytes.
func ParseMeminfo(text string) Meminfo {
	raw := map[string]uint64{}
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		n, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		raw[strings.TrimSuffix(fields[0], ":")] = n * 1024
	}
	return Meminfo{
		Total:     raw["MemTotal"],
		Free:      raw["MemFree"],
		Available: raw["MemAvailable"],
		Buffers:   raw["Buffers"],
		Cached:    raw["Cached"],
		SwapTotal: raw["SwapTotal"],
		SwapFree:  raw["SwapFree"],
	}
}

// LoadMeminfo reads /proc/meminfo.
func LoadMeminfo() (Meminfo, error) {
	text, err := Read("/proc/meminfo")
	if err != nil {
		return Meminfo{}, err
	}
	return ParseMeminfo(text), nil
}

// Loadavg is the triplet from /proc/loadavg.
type Loadavg struct {
	One, Five, Fifteen float64
	Running, Total     int
}

// ParseLoadavg parses one /proc/loadavg line.
func ParseLoadavg(text string) (Loadavg, error) {
	fields := strings.Fields(text)
	if len(fields) < 4 {
		return Loadavg{}, errors.New("loadavg: short record")
	}
	var out Loadavg
	var err error
	if out.One, err = strconv.ParseFloat(fields[0], 64); err != nil {
		return Loadavg{}, fmt.Errorf("loadavg: %w", err)
	}
	if out.Five, err = strconv.ParseFloat(fields[1], 64); err != nil {
		return Loadavg{}, fmt.Errorf("loadavg: %w", err)
	}
	if out.Fifteen, err = strconv.ParseFloat(fields[2], 64); err != nil {
		return Loadavg{}, fmt.Errorf("loadavg: %w", err)
	}
	parts := strings.SplitN(fields[3], "/", 2)
	if len(parts) == 2 {
		out.Running, _ = strconv.Atoi(parts[0])
		out.Total, _ = strconv.Atoi(parts[1])
	}
	return out, nil
}

// LoadLoadavg reads /proc/loadavg.
func LoadLoadavg() (Loadavg, error) {
	text, err := Read("/proc/loadavg")
	if err != nil {
		return Loadavg{}, err
	}
	return ParseLoadavg(text)
}

// Uptime reads seconds since boot from /proc/uptime.
func Uptime() (time.Duration, error) {
	text, err := Read("/proc/uptime")
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return 0, errors.New("uptime: empty")
	}
	sec, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("uptime: %w", err)
	}
	return time.Duration(sec * float64(time.Second)), nil
}

// Process is one task visible in procfs.
type Process struct {
	PID        int
	PPID       int
	Comm       string
	Cmdline    string
	State      byte
	UID        int
	VmRSS      uint64
	Threads    int
	UTime      uint64
	STime      uint64
	StartTicks uint64
}

// CPUTime is user plus system jiffies.
func (p Process) CPUTime() uint64 { return p.UTime + p.STime }

// Processes lists tasks under /proc. Missing tasks are skipped.
func Processes() ([]Process, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("read /proc: %w", err)
	}
	out := make([]Process, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		proc, ok := readProcess(pid)
		if ok {
			out = append(out, proc)
		}
	}
	return out, nil
}

func readProcess(pid int) (Process, bool) {
	base := filepath.Join("/proc", strconv.Itoa(pid))
	statText, err := Read(filepath.Join(base, "stat"))
	if err != nil {
		return Process{}, false
	}
	proc, err := ParseProcStat(statText)
	if err != nil {
		return Process{}, false
	}
	proc.PID = pid
	if cmd, err := os.ReadFile(filepath.Join(base, "cmdline")); err == nil {
		proc.Cmdline = strings.TrimSpace(strings.ReplaceAll(string(cmd), "\x00", " "))
	}
	if proc.Cmdline == "" {
		proc.Cmdline = proc.Comm
	}
	status, err := Read(filepath.Join(base, "status"))
	if err != nil {
		return proc, true
	}
	for _, line := range strings.Split(status, "\n") {
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		val = strings.TrimSpace(val)
		switch key {
		case "Uid":
			fields := strings.Fields(val)
			if len(fields) > 0 {
				proc.UID, _ = strconv.Atoi(fields[0])
			}
		case "VmRSS":
			fields := strings.Fields(val)
			if len(fields) > 0 {
				kb, err := strconv.ParseUint(fields[0], 10, 64)
				if err == nil {
					proc.VmRSS = kb * 1024
				}
			}
		case "Threads":
			proc.Threads, _ = strconv.Atoi(val)
		}
	}
	return proc, true
}

// ParseProcStat parses one /proc/<pid>/stat record. The command field may contain spaces.
func ParseProcStat(text string) (Process, error) {
	start := strings.IndexByte(text, '(')
	end := strings.LastIndexByte(text, ')')
	if start < 0 || end < start+1 {
		return Process{}, errors.New("proc stat: missing command")
	}
	rest := strings.Fields(text[end+1:])
	if len(rest) < 22 {
		return Process{}, errors.New("proc stat: short record")
	}
	ppid, err := strconv.Atoi(rest[1])
	if err != nil {
		return Process{}, fmt.Errorf("proc stat: ppid: %w", err)
	}
	utime, err := strconv.ParseUint(rest[11], 10, 64)
	if err != nil {
		return Process{}, fmt.Errorf("proc stat: utime: %w", err)
	}
	stime, err := strconv.ParseUint(rest[12], 10, 64)
	if err != nil {
		return Process{}, fmt.Errorf("proc stat: stime: %w", err)
	}
	threads, _ := strconv.Atoi(rest[17])
	startTicks, err := strconv.ParseUint(rest[19], 10, 64)
	if err != nil {
		return Process{}, fmt.Errorf("proc stat: starttime: %w", err)
	}
	state := byte('?')
	if rest[0] != "" {
		state = rest[0][0]
	}
	return Process{
		Comm:       text[start+1 : end],
		State:      state,
		PPID:       ppid,
		UTime:      utime,
		STime:      stime,
		Threads:    threads,
		StartTicks: startTicks,
	}, nil
}

// BootTime reads btime from /proc/stat.
func BootTime() (time.Time, error) {
	text, err := Read("/proc/stat")
	if err != nil {
		return time.Time{}, err
	}
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == "btime" {
			sec, err := strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				return time.Time{}, fmt.Errorf("btime: %w", err)
			}
			return time.Unix(sec, 0), nil
		}
	}
	return time.Time{}, errors.New("btime: not found")
}

// StatField returns the integer after a named key in /proc/stat.
func StatField(key string) (uint64, error) {
	text, err := Read("/proc/stat")
	if err != nil {
		return 0, err
	}
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == key {
			n, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return 0, fmt.Errorf("%s: %w", key, err)
			}
			return n, nil
		}
	}
	return 0, fmt.Errorf("%s: not found", key)
}

// CPUJiffies is one aggregated or per-core line from /proc/stat.
type CPUJiffies struct {
	Name   string
	User   uint64
	Nice   uint64
	System uint64
	Idle   uint64
	IOWait uint64
	IRQ    uint64
	Soft   uint64
	Steal  uint64
}

// Total returns the sum of accounted jiffies.
func (c CPUJiffies) Total() uint64 {
	return c.User + c.Nice + c.System + c.Idle + c.IOWait + c.IRQ + c.Soft + c.Steal
}

// ParseCPUStat parses cpu lines from /proc/stat.
func ParseCPUStat(text string) []CPUJiffies {
	var out []CPUJiffies
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 5 || !strings.HasPrefix(fields[0], "cpu") {
			continue
		}
		nums := make([]uint64, 8)
		for i := 0; i < len(nums) && i+1 < len(fields); i++ {
			nums[i], _ = strconv.ParseUint(fields[i+1], 10, 64)
		}
		out = append(out, CPUJiffies{
			Name: fields[0], User: nums[0], Nice: nums[1], System: nums[2], Idle: nums[3],
			IOWait: nums[4], IRQ: nums[5], Soft: nums[6], Steal: nums[7],
		})
	}
	return out
}

// FSSpace is filesystem capacity from statfs.
type FSSpace struct {
	Total, Used, Avail, Files, FilesFree uint64
}

// Statfs reads capacity for path. It does not invoke df.
func Statfs(path string) (FSSpace, error) {
	var st syscall.Statfs_t
	if err := syscall.Statfs(path, &st); err != nil {
		return FSSpace{}, fmt.Errorf("statfs %s: %w", path, err)
	}
	if st.Bsize <= 0 {
		return FSSpace{}, fmt.Errorf("statfs %s: invalid block size", path)
	}
	bsize := uint64(st.Bsize)
	total := st.Blocks * bsize
	var used uint64
	if st.Blocks >= st.Bfree {
		used = (st.Blocks - st.Bfree) * bsize
	}
	return FSSpace{
		Total:     total,
		Used:      used,
		Avail:     st.Bavail * bsize,
		Files:     st.Files,
		FilesFree: st.Ffree,
	}, nil
}

// Mount is one row of /proc/mounts.
type Mount struct {
	Device, Path, Type, Options string
}

// ParseMounts parses /proc/mounts, unescaping octal sequences.
func ParseMounts(text string) []Mount {
	var out []Mount
	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		out = append(out, Mount{
			Device:  unescapeMount(fields[0]),
			Path:    unescapeMount(fields[1]),
			Type:    unescapeMount(fields[2]),
			Options: unescapeMount(fields[3]),
		})
	}
	return out
}

func unescapeMount(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+3 < len(s) {
			n, err := strconv.ParseUint(s[i+1:i+4], 8, 8)
			if err == nil {
				b.WriteByte(byte(n))
				i += 3
				continue
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

// Mounts reads /proc/self/mounts.
func Mounts() ([]Mount, error) {
	text, err := Read("/proc/self/mounts")
	if err != nil {
		return nil, err
	}
	return ParseMounts(text), nil
}

// Diskstat is one /proc/diskstats row.
type Diskstat struct {
	Name                          string
	Reads, ReadSectors, ReadMS    uint64
	Writes, WriteSectors, WriteMS uint64
	InFlight, IOTicks, Weighted   uint64
}

// ParseDiskstats parses /proc/diskstats.
func ParseDiskstats(text string) []Diskstat {
	var out []Diskstat
	for _, line := range strings.Split(text, "\n") {
		f := strings.Fields(line)
		if len(f) < 14 {
			continue
		}
		nums := make([]uint64, 11)
		ok := true
		for i := 0; i < len(nums); i++ {
			n, err := strconv.ParseUint(f[i+3], 10, 64)
			if err != nil {
				ok = false
				break
			}
			nums[i] = n
		}
		if !ok {
			continue
		}
		out = append(out, Diskstat{
			Name: f[2], Reads: nums[0], ReadSectors: nums[2], ReadMS: nums[3],
			Writes: nums[4], WriteSectors: nums[6], WriteMS: nums[7],
			InFlight: nums[8], IOTicks: nums[9], Weighted: nums[10],
		})
	}
	return out
}

// NetDev is one /proc/net/dev interface.
type NetDev struct {
	Name                       string
	RxBytes, RxPackets, RxErrs uint64
	RxDrop                     uint64
	TxBytes, TxPackets, TxErrs uint64
	TxDrop                     uint64
}

// ParseNetDev parses /proc/net/dev.
func ParseNetDev(text string) []NetDev {
	var out []NetDev
	for _, line := range strings.Split(text, "\n") {
		if !strings.Contains(line, ":") {
			continue
		}
		name, rest, _ := strings.Cut(line, ":")
		f := strings.Fields(rest)
		if len(f) < 12 {
			continue
		}
		nums := make([]uint64, 12)
		ok := true
		for i := range nums {
			n, err := strconv.ParseUint(f[i], 10, 64)
			if err != nil {
				ok = false
				break
			}
			nums[i] = n
		}
		if !ok {
			continue
		}
		out = append(out, NetDev{
			Name:    strings.TrimSpace(name),
			RxBytes: nums[0], RxPackets: nums[1], RxErrs: nums[2], RxDrop: nums[3],
			TxBytes: nums[8], TxPackets: nums[9], TxErrs: nums[10], TxDrop: nums[11],
		})
	}
	return out
}

// TCPSock is one row from /proc/net/tcp or tcp6.
type TCPSock struct {
	Local, Remote string
	State         string
	Inode         string
}

var tcpStates = map[string]string{
	"01": "ESTABLISHED", "02": "SYN_SENT", "03": "SYN_RECV", "04": "FIN_WAIT1",
	"05": "FIN_WAIT2", "06": "TIME_WAIT", "07": "CLOSE", "08": "CLOSE_WAIT",
	"09": "LAST_ACK", "0A": "LISTEN", "0B": "CLOSING",
}

// ParseTCPTable parses a /proc/net/tcp* table.
func ParseTCPTable(text string, v6 bool) []TCPSock {
	var out []TCPSock
	lines := strings.Split(text, "\n")
	for _, line := range lines[1:] {
		f := strings.Fields(line)
		if len(f) < 10 {
			continue
		}
		state := tcpStates[strings.ToUpper(f[3])]
		if state == "" {
			state = f[3]
		}
		out = append(out, TCPSock{
			Local:  decodeSockAddr(f[1], v6),
			Remote: decodeSockAddr(f[2], v6),
			State:  state,
			Inode:  f[9],
		})
	}
	return out
}

func decodeSockAddr(field string, v6 bool) string {
	ip, port, ok := strings.Cut(field, ":")
	if !ok {
		return field
	}
	portN, err := strconv.ParseUint(port, 16, 16)
	if err != nil {
		return field
	}
	if v6 {
		return fmt.Sprintf("[%s]:%d", decodeIPv6(ip), portN)
	}
	return fmt.Sprintf("%s:%d", DecodeIPv4Hex(ip), portN)
}

// DecodeIPv4Hex decodes an 8-digit little-endian hex address from procfs.
func DecodeIPv4Hex(hex string) string {
	n, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return hex
	}
	return fmt.Sprintf("%d.%d.%d.%d", byte(n), byte(n>>8), byte(n>>16), byte(n>>24))
}

func decodeIPv6(hex string) string {
	if len(hex) != 32 {
		return hex
	}
	b := make([]byte, 16)
	for i := 0; i < 4; i++ {
		n, err := strconv.ParseUint(hex[i*8:(i+1)*8], 16, 32)
		if err != nil {
			return hex
		}
		// Each 32-bit word in the procfs table is little-endian.
		b[i*4] = byte(n)
		b[i*4+1] = byte(n >> 8)
		b[i*4+2] = byte(n >> 16)
		b[i*4+3] = byte(n >> 24)
	}
	parts := make([]string, 8)
	for i := 0; i < 8; i++ {
		parts[i] = fmt.Sprintf("%x", uint16(b[i*2])<<8|uint16(b[i*2+1]))
	}
	return strings.Join(parts, ":")
}

// TCPSockets reads IPv4 and IPv6 TCP tables. A missing table is skipped.
func TCPSockets() ([]TCPSock, error) {
	var out []TCPSock
	var saw bool
	for _, item := range []struct {
		path string
		v6   bool
	}{{"/proc/net/tcp", false}, {"/proc/net/tcp6", true}} {
		text, err := Read(item.path)
		if err != nil {
			continue
		}
		saw = true
		out = append(out, ParseTCPTable(text, item.v6)...)
	}
	if !saw {
		return nil, errors.New("tcp tables unavailable")
	}
	return out, nil
}

// UDPCount counts rows in a UDP table, excluding the header.
func UDPCount(path string) (int, error) {
	text, err := Read(path)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "sl") || strings.Contains(line, "local_address") {
			continue
		}
		n++
	}
	return n, nil
}

// Route is one IPv4 route.
type Route struct {
	Iface, Destination, Gateway, Mask string
	Metric                            int
}

// ParseRoutes parses /proc/net/route.
func ParseRoutes(text string) []Route {
	var out []Route
	lines := strings.Split(text, "\n")
	for _, line := range lines[1:] {
		f := strings.Fields(line)
		if len(f) < 8 {
			continue
		}
		metric, _ := strconv.ParseInt(f[6], 16, 32)
		out = append(out, Route{
			Iface:       f[0],
			Destination: DecodeIPv4Hex(f[1]),
			Gateway:     DecodeIPv4Hex(f[2]),
			Mask:        DecodeIPv4Hex(f[7]),
			Metric:      int(metric),
		})
	}
	return out
}

// UtmpRecord is one parsed utmp/wtmp entry.
type UtmpRecord struct {
	Type             int
	PID              int
	Line, User, Host string
	When             time.Time
}

const utmpSize = 384

// ParseUtmp parses Linux utmp records. Incomplete trailing bytes are ignored.
func ParseUtmp(b []byte) []UtmpRecord {
	var out []UtmpRecord
	for len(b) >= utmpSize {
		rec := b[:utmpSize]
		b = b[utmpSize:]
		typ := int(binary.LittleEndian.Uint16(rec[0:2]))
		pid := int(int32(binary.LittleEndian.Uint32(rec[4:8])))
		line := cString(rec[8:40])
		user := cString(rec[44:76])
		host := cString(rec[76:332])
		sec := int64(int32(binary.LittleEndian.Uint32(rec[340:344])))
		out = append(out, UtmpRecord{Type: typ, PID: pid, Line: line, User: user, Host: host, When: time.Unix(sec, 0)})
	}
	return out
}

func cString(b []byte) string {
	if i := indexByte(b, 0); i >= 0 {
		b = b[:i]
	}
	return string(b)
}

func indexByte(b []byte, c byte) int {
	for i, v := range b {
		if v == c {
			return i
		}
	}
	return -1
}

// ReadUtmp reads a utmp-format file.
func ReadUtmp(path string) ([]UtmpRecord, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return ParseUtmp(b), nil
}

// KmsgLine is one kernel ring-buffer record.
type KmsgLine struct {
	Priority int
	Text     string
}

// ReadKmsg reads up to limit records from /dev/kmsg without blocking.
func ReadKmsg(limit int) ([]KmsgLine, error) {
	f, err := os.OpenFile("/dev/kmsg", os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, fmt.Errorf("open /dev/kmsg: %w", err)
	}
	defer f.Close()
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek /dev/kmsg: %w", err)
	}
	r := bufio.NewReader(f)
	out := make([]KmsgLine, 0, limit)
	for len(out) < limit {
		line, err := r.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, syscall.EAGAIN) {
				break
			}
			if len(out) > 0 {
				break
			}
			return out, fmt.Errorf("read /dev/kmsg: %w", err)
		}
		pri := 6
		text := strings.TrimRight(line, "\n")
		if head, msg, ok := strings.Cut(text, ";"); ok {
			if n, err := strconv.Atoi(strings.SplitN(head, ",", 2)[0]); err == nil {
				pri = n
			}
			text = msg
		}
		out = append(out, KmsgLine{Priority: pri, Text: text})
	}
	return out, nil
}

// CountFDs returns open and deleted-open descriptor counts for a pid.
func CountFDs(pid int) (open, deleted int, err error) {
	dir := filepath.Join("/proc", strconv.Itoa(pid), "fd")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0, err
	}
	for _, entry := range entries {
		target, err := os.Readlink(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		open++
		if strings.HasSuffix(target, "(deleted)") {
			deleted++
		}
	}
	return open, deleted, nil
}

// FileSize is a walked file.
type FileSize struct {
	Path string
	Size int64
}

// WalkSizes sums regular files under root, up to maxFiles. It does not follow directory symlinks.
func WalkSizes(root string, maxFiles int) (total int64, files []FileSize, truncated bool, err error) {
	count := 0
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() || d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		info, err := d.Info()
		if err != nil || !info.Mode().IsRegular() {
			return nil
		}
		count++
		if count > maxFiles {
			truncated = true
			return filepath.SkipAll
		}
		total += info.Size()
		files = append(files, FileSize{Path: path, Size: info.Size()})
		return nil
	})
	sort.Slice(files, func(i, j int) bool { return files[i].Size > files[j].Size })
	return total, files, truncated, err
}

// ReadTrimmed reads a sysfs file and trims whitespace. Missing files return an error.
func ReadTrimmed(path string) (string, error) {
	text, err := Read(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

// Sorted walks names and returns them in order.
func Sorted(names []string) []string {
	out := append([]string(nil), names...)
	sort.Strings(out)
	return out
}

// Percent returns used/total as a percentage. total of 0 yields 0.
func Percent(used, total uint64) float64 {
	if total == 0 {
		return 0
	}
	return float64(used) * 100 / float64(total)
}

// Ticks converts user_hz jiffies to a duration. Linux user_hz is 100.
func Ticks(ticks uint64) time.Duration {
	return time.Duration(ticks) * time.Second / userHZ
}

// ScanLines returns lines from the first readable path for which match is true, up to limit.
func ScanLines(paths []string, match func(string) bool, limit int) ([]string, string, error) {
	var last error
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			last = err
			continue
		}
		var hits []string
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for sc.Scan() {
			line := sc.Text()
			if match(line) {
				hits = append(hits, line)
			}
		}
		closeErr := f.Close()
		if err := sc.Err(); err != nil {
			return nil, path, fmt.Errorf("scan %s: %w", path, err)
		}
		if closeErr != nil {
			return nil, path, closeErr
		}
		if len(hits) > limit {
			hits = hits[len(hits)-limit:]
		}
		return hits, path, nil
	}
	if last == nil {
		last = errors.New("no log path configured")
	}
	return nil, "", last
}

// ChildDirs lists directory names under path.
func ChildDirs(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}
