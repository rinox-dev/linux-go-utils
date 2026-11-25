// Package catalog is the list of monitoring commands. Implementations live in cmd/<name>.
package catalog

// Command is one monitoring utility.
type Command struct {
	Name        string
	Description string
}

// All returns the commands in stable alphabetical order.
func All() []Command {
	out := append([]Command(nil), commands...)
	return out
}

// Has reports whether name is a known command.
func Has(name string) bool {
	for _, c := range commands {
		if c.Name == name {
			return true
		}
	}
	return false
}

// commands is kept alphabetical so monitor list stays stable.
var commands = []Command{
	{"alertrelay", "Post a host snapshot to a webhook URL"},
	{"arpwatch", "IPv4 neighbor cache"},
	{"authwatch", "Authentication failures and successes from local logs"},
	{"bandwidth", "Per-interface traffic rate from a one-second sample"},
	{"batteryhealth", "Battery design capacity versus current full capacity"},
	{"batterywatch", "Battery state and capacity"},
	{"blockqueue", "Block device in-flight I/O and queue settings"},
	{"bootlog", "Recent kernel ring-buffer lines"},
	{"certwatch", "Certificate expiry under a directory"},
	{"cgroupwatch", "Cgroup version and controllers"},
	{"clockwatch", "Clock source and uptime"},
	{"conncount", "TCP and UDP socket counts"},
	{"conntrackwatch", "Connection-tracking counters"},
	{"containerlist", "Container-like cgroups and runtime sockets"},
	{"containerstats", "CPU and memory of container cgroups"},
	{"contextswitch", "Context-switch counter since boot"},
	{"cpufreqwatch", "Current CPU frequency"},
	{"cpuidlewatch", "CPU idle-state time"},
	{"deletedopen", "Open file descriptors pointing at deleted files"},
	{"diskerrors", "Block device error counters when the driver exposes them"},
	{"disklatency", "Block I/O time counters"},
	{"diskusage", "Filesystem space for a path"},
	{"diskwatch", "Filesystem use for every mount"},
	{"dnscheck", "Configured name servers"},
	{"failedunits", "Systemd units with an invocation record and no active cgroup"},
	{"fanwatch", "Fan sensor readings"},
	{"fdcount", "Open file descriptors by process"},
	{"filegrowth", "Largest files under a path"},
	{"forkrate", "Process creation counter since boot"},
	{"hardwarechanges", "Device class inventory snapshot"},
	{"healthcheck", "Load, memory, disk, and zombie health"},
	{"hostinfo", "Host, distribution, kernel, and architecture"},
	{"httpcheck", "HTTP status of a URL"},
	{"hwreport", "Short hardware summary"},
	{"inodewatch", "Inode use for every mount"},
	{"integrityscan", "SHA-256 of files under a path"},
	{"ioactivity", "Block device read and write counters"},
	{"iowait", "CPU I/O wait share"},
	{"journalerrors", "Error lines in text logs when the journal is present"},
	{"kernelwarn", "Kernel warning and error lines"},
	{"listenchanges", "Listening TCP and UDP sockets"},
	{"listenports", "Listening TCP sockets"},
	{"loadwatch", "Load averages and runnable tasks"},
	{"loginhistory", "Recent login records from wtmp"},
	{"loginwatch", "Logged-in sessions from utmp"},
	{"logsize", "Largest files under a log directory"},
	{"logwatch", "Recent error-like lines in a log directory"},
	{"meminfo", "Memory and swap counters"},
	{"metricdump", "Selected host metrics as JSON"},
	{"mountlist", "Mounted filesystems"},
	{"mountwatch", "Mount table with read-only and network flags"},
	{"mtucheck", "Interface MTU values"},
	{"namespacewatch", "Namespaces of the current process"},
	{"netbaseline", "Network counters for a later comparison"},
	{"neterrors", "Interface errors and drops"},
	{"netinterfaces", "Interface addresses and operational state"},
	{"nettraffic", "Interface byte and packet counters"},
	{"netwatch", "Interface operational state"},
	{"oomwatch", "Out-of-memory kill counter"},
	{"packetloss", "ICMP and IP error counters"},
	{"pcieinventory", "PCI device vendor and device IDs"},
	{"permcheck", "Mode of critical account files"},
	{"pingcheck", "ICMP echo to a host"},
	{"portcheck", "TCP connect to a host and port"},
	{"powerprofile", "ACPI platform profile"},
	{"powerwatch", "Power supply status"},
	{"procage", "Process age since start"},
	{"proctree", "Process parent identifiers"},
	{"procwatch", "Process count by state"},
	{"promtext", "Selected host metrics in Prometheus text format"},
	{"psfind", "Find processes by command-line text"},
	{"readonlyfs", "Filesystems mounted read-only"},
	{"routewatch", "IPv4 route table"},
	{"sensorlist", "Hardware sensor names"},
	{"servicecheck", "Systemd services with an invocation record"},
	{"serviceprobe", "TCP connect to one or more host:port endpoints"},
	{"slabwatch", "Kernel slab cache counters"},
	{"socketcount", "Socket summary from sockstat"},
	{"sshwatch", "sshd processes and failed SSH log lines"},
	{"storageinventory", "Block devices and sizes"},
	{"sudoaudit", "Sudo lines in authentication logs"},
	{"suidlist", "Setuid and setgid files under a path"},
	{"swapwatch", "Swap use and paging counters"},
	{"syswatch", "Host, load, and memory snapshot"},
	{"tcpstates", "TCP sockets grouped by state"},
	{"thermalwatch", "Thermal zone temperatures"},
	{"thermalzones", "Thermal zone names and temperatures"},
	{"threadcount", "Thread counts by process"},
	{"throttlecheck", "CPU thermal throttling counters"},
	{"tlsdays", "Days until a TLS certificate expires"},
	{"tmpwatch", "Temporary and tmpfs mounts"},
	{"topcpu", "Processes with the most CPU time"},
	{"topmem", "Processes with the most resident memory"},
	{"updatescheck", "Installed package database size, without calling package managers"},
	{"uptime", "Uptime and load average"},
	{"usbwatch", "USB device names"},
	{"userprocs", "Process counts by user id"},
	{"userwatch", "Local accounts from the passwd database"},
	{"watchdogcheck", "Kernel watchdog devices"},
	{"zombiecheck", "Zombie processes"},
}
