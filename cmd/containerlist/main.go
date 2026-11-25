// containerlist finds container runtime sockets and container-like cgroups.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

)

func main() {
	for _, sock := range []string{"/run/docker.sock", "/run/containerd/containerd.sock", "/run/podman/podman.sock"} {
		if _, err := os.Stat(sock); err != nil {
			fmt.Printf("%s absent\n", sock)
			continue
		}
		fmt.Printf("%s present\n", sock)
	}
	found := 0
	_ = filepath.WalkDir("/sys/fs/cgroup", func(path string, d os.DirEntry, err error) error {
		if err != nil || !d.IsDir() || found >= 50 {
			return nil
		}
		base := filepath.Base(path)
		if strings.Contains(base, "docker-") || strings.Contains(base, "libpod-") || strings.Contains(base, "cri-containerd-") {
			fmt.Printf("cgroup %s\n", path)
			found++
		}
		return nil
	})
	if found == 0 {
		fmt.Println("no container cgroups")
	}
}
