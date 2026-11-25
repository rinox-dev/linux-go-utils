// monitor dispatches the Linux monitoring utilities by command name.
// It uses procfs, sysfs, and standard Linux interfaces, with no architecture-specific code.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type utility struct {
	description string
	paths       []string
}

// Each command is a focused view over Linux kernel interfaces. Missing optional
// interfaces are reported as unavailable, which is common across devices.
var utilities = map[string]utility{
	"hostinfo":         {"Host, distribution, kernel, and architecture", []string{"/etc/os-release", "/proc/sys/kernel/ostype", "/proc/sys/kernel/osrelease", "/proc/sys/kernel/hostname"}},
	"uptime":           {"System uptime and load average", []string{"/proc/uptime", "/proc/loadavg"}},
	"loadwatch":        {"Current load averages and runnable tasks", []string{"/proc/loadavg"}},
	"procwatch":        {"Process count and state summary", []string{"/proc/stat"}},
	"psfind":           {"Find processes by command-line text", []string{"/proc"}},
	"topcpu":           {"CPU time counters by core", []string{"/proc/stat"}},
	"topmem":           {"Largest memory consumers (system counters)", []string{"/proc/meminfo"}},
	"proctree":         {"Process and parent identifiers", []string{"/proc"}},
	"procage":          {"Process start-time counters", []string{"/proc"}},
	"zombiecheck":      {"Zombie process count", []string{"/proc"}},
	"fdcount":          {"Open file descriptor counts", []string{"/proc"}},
	"threadcount":      {"Thread counts", []string{"/proc"}},
	"contextswitch":    {"Context switch counters", []string{"/proc/stat"}},
	"forkrate":         {"Fork counter since boot", []string{"/proc/stat"}},
	"loginwatch":       {"Logged-in sessions from utmp if available", []string{"/run/utmp", "/var/run/utmp"}},
	"userprocs":        {"Processes grouped by user (proc visibility dependent)", []string{"/proc"}},
	"servicecheck":     {"Systemd service state when systemd is present", []string{"/run/systemd/system"}},
	"bootlog":          {"Kernel ring buffer access status", []string{"/dev/kmsg"}},
	"kernelwarn":       {"Kernel log access status", []string{"/dev/kmsg"}},
	"syswatch":         {"Compact host and load snapshot", []string{"/proc/uptime", "/proc/loadavg", "/proc/meminfo"}},
	"meminfo":          {"Memory and swap counters", []string{"/proc/meminfo"}},
	"swapwatch":        {"Swap totals and use", []string{"/proc/meminfo", "/proc/vmstat"}},
	"oomwatch":         {"OOM-related kernel counters", []string{"/proc/vmstat"}},
	"slabwatch":        {"Kernel slab allocator counters", []string{"/proc/slabinfo"}},
	"diskusage":        {"Filesystem space for a path", []string{"/proc/mounts"}},
	"diskwatch":        {"Filesystem mount inventory for threshold checks", []string{"/proc/mounts"}},
	"inodewatch":       {"Filesystem mount inventory for inode checks", []string{"/proc/mounts"}},
	"mountlist":        {"Mounted filesystems", []string{"/proc/mounts"}},
	"mountwatch":       {"Current mount table snapshot", []string{"/proc/self/mounts"}},
	"ioactivity":       {"Block device I/O counters", []string{"/proc/diskstats"}},
	"iowait":           {"CPU I/O wait counters", []string{"/proc/stat"}},
	"diskerrors":       {"Block device counters (error details may be driver-specific)", []string{"/proc/diskstats"}},
	"filegrowth":       {"Directory size inspection (pass a path)", []string{"."}},
	"logsize":          {"Log file size inspection (pass a path)", []string{"/var/log"}},
	"deletedopen":      {"Deleted open file candidates (permission dependent)", []string{"/proc"}},
	"tmpwatch":         {"Temporary filesystem mounts", []string{"/proc/mounts"}},
	"readonlyfs":       {"Filesystem mount flags", []string{"/proc/mounts"}},
	"disklatency":      {"Block I/O counters for latency sampling", []string{"/proc/diskstats"}},
	"blockqueue":       {"Block device queue and I/O counters", []string{"/proc/diskstats"}},
	"storageinventory": {"Block devices and partitions", []string{"/sys/block"}},
	"netinterfaces":    {"Network interface addresses and status", []string{"/proc/net/dev", "/sys/class/net"}},
	"nettraffic":       {"Per-interface byte and packet counters", []string{"/proc/net/dev"}},
	"neterrors":        {"Per-interface errors and drops", []string{"/proc/net/dev"}},
	"netwatch":         {"Interface operational state", []string{"/sys/class/net"}},
	"conncount":        {"TCP/UDP socket table counts", []string{"/proc/net/tcp", "/proc/net/udp", "/proc/net/tcp6", "/proc/net/udp6"}},
	"tcpstates":        {"TCP sockets grouped by kernel state", []string{"/proc/net/tcp", "/proc/net/tcp6"}},
	"listenports":      {"Listening TCP/UDP socket table", []string{"/proc/net/tcp", "/proc/net/udp", "/proc/net/tcp6", "/proc/net/udp6"}},
	"socketcount":      {"Socket table counts by protocol", []string{"/proc/net/sockstat", "/proc/net/sockstat6"}},
	"routewatch":       {"IPv4 and IPv6 route tables", []string{"/proc/net/route", "/proc/net/ipv6_route"}},
	"dnscheck":         {"Configured name servers", []string{"/etc/resolv.conf"}},
	"pingcheck":        {"Network interface and route snapshot", []string{"/proc/net/route", "/proc/net/dev"}},
	"httpcheck":        {"HTTP target supplied as an argument (use curl for full checks)", []string{"/etc/resolv.conf"}},
	"tlsdays":          {"TLS endpoint supplied as an argument (use openssl for details)", []string{"/etc/resolv.conf"}},
	"portcheck":        {"TCP endpoint supplied as an argument", []string{"/proc/net/route"}},
	"bandwidth":        {"Per-interface traffic counters; sample twice for a rate", []string{"/proc/net/dev"}},
	"conntrackwatch":   {"Connection tracking counters when enabled", []string{"/proc/net/stat/nf_conntrack"}},
	"arpwatch":         {"IPv4 neighbor cache", []string{"/proc/net/arp"}},
	"mtucheck":         {"Interface MTU values", []string{"/sys/class/net"}},
	"packetloss":       {"ICMP counters (destination probes require ping)", []string{"/proc/net/snmp"}},
	"netbaseline":      {"Network interface counters for baseline sampling", []string{"/proc/net/dev"}},
	"thermalwatch":     {"Thermal zone temperatures", []string{"/sys/class/thermal"}},
	"fanwatch":         {"Fan sensors when exposed by drivers", []string{"/sys/class/hwmon"}},
	"batterywatch":     {"Battery state and capacity", []string{"/sys/class/power_supply"}},
	"powerwatch":       {"Power supply status and capacity", []string{"/sys/class/power_supply"}},
	"cpufreqwatch":     {"CPU frequency interface", []string{"/sys/devices/system/cpu"}},
	"cpuidlewatch":     {"CPU idle states", []string{"/sys/devices/system/cpu"}},
	"throttlecheck":    {"Thermal throttling interfaces when present", []string{"/sys/devices/system/cpu"}},
	"sensorlist":       {"Hardware sensor inventory", []string{"/sys/class/hwmon"}},
	"hardwarechanges":  {"Device inventory snapshot", []string{"/sys/class"}},
	"usbwatch":         {"USB device inventory", []string{"/sys/bus/usb/devices"}},
	"pcieinventory":    {"PCI device inventory", []string{"/sys/bus/pci/devices"}},
	"batteryhealth":    {"Battery design and current capacity", []string{"/sys/class/power_supply"}},
	"thermalzones":     {"Thermal zone names and temperatures", []string{"/sys/class/thermal"}},
	"powerprofile":     {"Power profile interface when supported", []string{"/sys/firmware/acpi/platform_profile"}},
	"hwreport":         {"Portable hardware summary", []string{"/proc/cpuinfo", "/proc/meminfo", "/sys/class/thermal"}},
	"authwatch":        {"Authentication log locations (permissions vary)", []string{"/var/log/auth.log", "/var/log/secure"}},
	"sudoaudit":        {"Sudo log locations", []string{"/var/log/auth.log", "/var/log/secure"}},
	"userwatch":        {"Local account database metadata", []string{"/etc/passwd"}},
	"sshwatch":         {"SSH service and log locations", []string{"/var/log/auth.log", "/var/log/secure", "/run/sshd"}},
	"permcheck":        {"Critical file permissions", []string{"/etc/passwd", "/etc/shadow", "/etc/sudoers"}},
	"suidlist":         {"SUID/SGID scan root (pass a path to scan)", []string{"/usr/bin"}},
	"listenchanges":    {"Listening socket table snapshot", []string{"/proc/net/tcp", "/proc/net/udp"}},
	"logwatch":         {"Log files available for pattern scanning", []string{"/var/log"}},
	"journalerrors":    {"Systemd journal availability", []string{"/run/log/journal", "/var/log/journal"}},
	"clockwatch":       {"System clock and kernel time status", []string{"/proc/uptime", "/sys/devices/system/clocksource/clocksource0/current_clocksource"}},
	"failedunits":      {"Failed systemd unit query availability", []string{"/run/systemd/system"}},
	"updatescheck":     {"Package manager availability", []string{"/usr/bin/pacman", "/usr/bin/apt", "/usr/bin/dnf"}},
	"certwatch":        {"Certificate scan root (pass a path to inspect)", []string{"/etc/ssl/certs"}},
	"loginhistory":     {"Login record locations", []string{"/var/log/wtmp", "/var/log/btmp"}},
	"integrityscan":    {"File hash inspection (pass a path)", []string{"."}},
	"cgroupwatch":      {"Cgroup version and controller files", []string{"/sys/fs/cgroup/cgroup.controllers", "/proc/self/cgroup"}},
	"containerlist":    {"Container runtime availability", []string{"/run/containerd/containerd.sock", "/run/docker.sock", "/run/podman/podman.sock"}},
	"containerstats":   {"Container runtime availability for metrics", []string{"/sys/fs/cgroup"}},
	"namespacewatch":   {"Current process namespace links", []string{"/proc/self/ns"}},
	"healthcheck":      {"Host health snapshot", []string{"/proc/loadavg", "/proc/meminfo", "/proc/uptime"}},
	"watchdogcheck":    {"Kernel watchdog device availability", []string{"/dev/watchdog", "/sys/class/watchdog"}},
	"serviceprobe":     {"Local service endpoints are supplied as arguments", []string{"/proc/net/tcp", "/proc/net/tcp6"}},
	"metricdump":       {"Print selected procfs metrics as JSON", []string{"/proc/uptime", "/proc/loadavg", "/proc/meminfo"}},
	"promtext":         {"Print selected metrics in Prometheus text format", []string{"/proc/uptime", "/proc/loadavg"}},
	"alertrelay":       {"Webhook endpoint is supplied as an argument", []string{"/proc/loadavg", "/proc/meminfo"}},
}

func main() {
	name := filepath.Base(os.Args[0])
	args := os.Args[1:]
	if name == "monitor" || name == "main" {
		if len(args) == 0 || args[0] == "help" || args[0] == "list" {
			list()
			return
		}
		name, args = args[0], args[1:]
	}
	u, ok := utilities[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown utility %q; run monitor list\n", name)
		os.Exit(2)
	}
	if name == "psfind" {
		if len(args) != 1 {
			usage(name, "SEARCH_TERM")
		}
		findProcesses(args[0])
		return
	}
	if name == "diskusage" {
		if len(args) > 1 {
			usage(name, "[PATH]")
		}
		path := "."
		if len(args) == 1 {
			path = args[0]
		}
		var stat syscall.Statfs_t
		err := syscall.Statfs(path, &stat)
		if err != nil {
			fail(err)
		}
		total := uint64(stat.Blocks) * uint64(stat.Bsize)
		used := (uint64(stat.Blocks) - uint64(stat.Bfree)) * uint64(stat.Bsize)
		available := uint64(stat.Bavail) * uint64(stat.Bsize)
		percent := 0.0
		if total > 0 {
			percent = float64(used) * 100 / float64(total)
		}
		fmt.Printf("Path: %s\nTotal: %s\nUsed: %s (%.1f%%)\nAvailable: %s\n", path, human(total), human(used), percent, human(available))
		return
	}
	if len(args) != 0 && name != "diskusage" {
		fmt.Fprintf(os.Stderr, "%s: this utility currently takes no arguments\n", name)
		os.Exit(2)
	}
	if name == "metricdump" || name == "promtext" {
		printMetrics(name)
		return
	}
	if name == "procwatch" || name == "zombiecheck" || name == "userprocs" || name == "threadcount" {
		processSummary(name)
		return
	}
	fmt.Printf("%s — %s\nTimestamp: %s\n", name, u.description, time.Now().Format(time.RFC3339))
	for _, path := range u.paths {
		if path == "/proc" || strings.HasPrefix(path, "/sys/class/") || path == "/sys/block" || path == "/sys/class" || path == "/sys/bus/usb/devices" || path == "/sys/bus/pci/devices" || path == "/sys/class/power_supply" || path == "/sys/class/thermal" || path == "/sys/class/hwmon" || path == "/sys/devices/system/cpu" || path == "/proc/self/ns" || path == "/sys/fs/cgroup" || path == "/var/log" || path == "/usr/bin" {
			printDirectory(path)
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("\n[%s] unavailable: %v\n", path, err)
			continue
		}
		fmt.Printf("\n[%s]\n%s", path, data)
		if len(data) > 0 && data[len(data)-1] != '\n' {
			fmt.Println()
		}
	}
}

func human(n uint64) string {
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}
	v := float64(n)
	i := 0
	for v >= 1024 && i < len(units)-1 {
		v /= 1024
		i++
	}
	return fmt.Sprintf("%.1f %s", v, units[i])
}

func processSummary(name string) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		fail(err)
	}
	counts := map[string]int{"total": 0, "zombie": 0, "running": 0, "sleeping": 0, "stopped": 0, "threads": 0}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		var pid int
		if _, err := fmt.Sscanf(entry.Name(), "%d", &pid); err != nil {
			continue
		}
		data, err := os.ReadFile(filepath.Join("/proc", entry.Name(), "status"))
		if err != nil {
			continue
		}
		counts["total"]++
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "State:") && len(line) > 7 {
				switch line[7] {
				case 'Z':
					counts["zombie"]++
				case 'R':
					counts["running"]++
				case 'S', 'I':
					counts["sleeping"]++
				case 'T', 't':
					counts["stopped"]++
				}
			}
			if strings.HasPrefix(line, "Threads:") {
				var n int
				if _, err := fmt.Sscanf(strings.TrimPrefix(line, "Threads:"), "%d", &n); err == nil {
					counts["threads"] += n
				}
			}
		}
	}
	if name == "zombiecheck" {
		fmt.Printf("Zombie processes: %d\n", counts["zombie"])
		return
	}
	if name == "threadcount" {
		fmt.Printf("Processes: %d\nThreads: %d\n", counts["total"], counts["threads"])
		return
	}
	fmt.Printf("Processes: %d\nRunning: %d\nSleeping: %d\nStopped: %d\nZombie: %d\nThreads: %d\n", counts["total"], counts["running"], counts["sleeping"], counts["stopped"], counts["zombie"], counts["threads"])
}

func printMetrics(name string) {
	readNumber := func(path string) float64 {
		b, err := os.ReadFile(path)
		if err != nil {
			return 0
		}
		var value float64
		_, _ = fmt.Sscanf(strings.Fields(string(b))[0], "%f", &value)
		return value
	}
	up := readNumber("/proc/uptime")
	load, _ := os.ReadFile("/proc/loadavg")
	fields := strings.Fields(string(load))
	load1, load5, load15 := 0.0, 0.0, 0.0
	if len(fields) >= 3 {
		_, _ = fmt.Sscanf(fields[0], "%f", &load1)
		_, _ = fmt.Sscanf(fields[1], "%f", &load5)
		_, _ = fmt.Sscanf(fields[2], "%f", &load15)
	}
	mem := map[string]uint64{}
	if b, err := os.ReadFile("/proc/meminfo"); err == nil {
		for _, line := range strings.Split(string(b), "\n") {
			var key string
			var value uint64
			if _, err := fmt.Sscanf(line, "%s %d", &key, &value); err == nil {
				mem[strings.TrimSuffix(key, ":")] = value * 1024
			}
		}
	}
	if name == "metricdump" {
		fmt.Printf("{\"uptime_seconds\":%.2f,\"load1\":%.2f,\"load5\":%.2f,\"load15\":%.2f,\"memory_total_bytes\":%d,\"memory_available_bytes\":%d}\n", up, load1, load5, load15, mem["MemTotal"], mem["MemAvailable"])
		return
	}
	fmt.Printf("linux_uptime_seconds %.2f\nlinux_load1 %.2f\nlinux_load5 %.2f\nlinux_load15 %.2f\nlinux_memory_total_bytes %d\nlinux_memory_available_bytes %d\n", up, load1, load5, load15, mem["MemTotal"], mem["MemAvailable"])
}

func list() {
	names := make([]string, 0, len(utilities))
	for name := range utilities {
		names = append(names, name)
	}
	// Stable alphabetical listing without importing sorting for generated wrappers.
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			if names[j] < names[i] {
				names[i], names[j] = names[j], names[i]
			}
		}
	}
	for _, name := range names {
		fmt.Printf("%-20s %s\n", name, utilities[name].description)
	}
}

func printDirectory(path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("\n[%s] unavailable: %v\n", path, err)
		return
	}
	fmt.Printf("\n[%s] (%d entries)\n", path, len(entries))
	limit := len(entries)
	if limit > 50 {
		limit = 50
	}
	for _, entry := range entries[:limit] {
		fmt.Println(entry.Name())
	}
	if len(entries) > limit {
		fmt.Printf("... and %d more\n", len(entries)-limit)
	}
}

func findProcesses(needle string) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		fail(err)
	}
	found := false
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		var pid int
		if _, err := fmt.Sscanf(entry.Name(), "%d", &pid); err != nil {
			continue
		}
		data, err := os.ReadFile(filepath.Join("/proc", entry.Name(), "cmdline"))
		if err != nil {
			continue
		}
		command := strings.TrimSpace(strings.ReplaceAll(string(data), "\x00", " "))
		if strings.Contains(strings.ToLower(command), strings.ToLower(needle)) {
			fmt.Printf("%-8d %s\n", pid, command)
			found = true
		}
	}
	if !found {
		fmt.Fprintf(os.Stderr, "no process matched %q\n", needle)
	}
}

func usage(name, args string) { fmt.Fprintf(os.Stderr, "usage: %s %s\n", name, args); os.Exit(2) }
func fail(err error)          { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
