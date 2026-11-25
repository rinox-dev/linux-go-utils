package sys

import (
	"testing"
	"time"
)

func TestParseMeminfo(t *testing.T) {
	in := "MemTotal:       2048 kB\nMemAvailable:   1024 kB\nSwapTotal:         0 kB\n"
	got := ParseMeminfo(in)
	if got.Total != 2048*1024 || got.Available != 1024*1024 || got.SwapTotal != 0 {
		t.Fatalf("ParseMeminfo() = %+v", got)
	}
}

func TestParseLoadavg(t *testing.T) {
	got, err := ParseLoadavg("0.10 0.20 0.30 1/100 9\n")
	if err != nil {
		t.Fatal(err)
	}
	if got.One != 0.10 || got.Running != 1 || got.Total != 100 {
		t.Fatalf("ParseLoadavg() = %+v", got)
	}
}

func TestParseProcStat(t *testing.T) {
	in := "12 (foo bar) S 1 12 12 0 0 0 0 0 0 0 10 20 0 0 20 0 3 0 500 0 0 0\n"
	got, err := ParseProcStat(in)
	if err != nil {
		t.Fatal(err)
	}
	if got.Comm != "foo bar" || got.State != 'S' || got.PPID != 1 || got.UTime != 10 || got.STime != 20 || got.Threads != 3 || got.StartTicks != 500 {
		t.Fatalf("ParseProcStat() = %+v", got)
	}
}

func TestDecodeIPv4Hex(t *testing.T) {
	if got := DecodeIPv4Hex("0100000A"); got != "10.0.0.1" {
		t.Fatalf("DecodeIPv4Hex() = %s", got)
	}
}

func TestParseMountsUnescapes(t *testing.T) {
	got := ParseMounts(`sysfs /sys sysfs rw,nosuid 0 0` + "\n" + `/dev/sda1 /mnt/my\040disk ext4 rw 0 0` + "\n")
	if len(got) != 2 || got[1].Path != "/mnt/my disk" {
		t.Fatalf("ParseMounts() = %+v", got)
	}
}

func TestParseNetDev(t *testing.T) {
	in := "Inter-|   Receive                                                |  Transmit\n face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed\n  lo: 10 1 0 0 0 0 0 0 20 2 0 0 0 0 0 0\n"
	got := ParseNetDev(in)
	if len(got) != 1 || got[0].Name != "lo" || got[0].RxBytes != 10 || got[0].TxBytes != 20 {
		t.Fatalf("ParseNetDev() = %+v", got)
	}
}

func TestParseTCPTable(t *testing.T) {
	in := "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n   0: 0100007F:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 1234 1 0000000000000000 100 0 0 10 0\n"
	got := ParseTCPTable(in, false)
	if len(got) != 1 || got[0].State != "LISTEN" || got[0].Local != "127.0.0.1:8080" {
		t.Fatalf("ParseTCPTable() = %+v", got)
	}
}

func TestHuman(t *testing.T) {
	if got := Human(512); got != "512 B" {
		t.Fatalf("Human(512) = %s", got)
	}
	if got := Human(1536); got != "1.5 KiB" {
		t.Fatalf("Human(1536) = %s", got)
	}
}

func TestParseUtmpSize(t *testing.T) {
	b := make([]byte, utmpSize)
	b[0] = 7
	copy(b[44:48], []byte("root"))
	b[340] = 1
	got := ParseUtmp(b)
	if len(got) != 1 || got[0].Type != 7 || got[0].User != "root" || !got[0].When.Equal(time.Unix(1, 0)) {
		t.Fatalf("ParseUtmp() = %+v", got)
	}
}

func TestPercentZero(t *testing.T) {
	if Percent(1, 0) != 0 {
		t.Fatal("expected 0")
	}
}
