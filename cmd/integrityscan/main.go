// integrityscan prints SHA-256 digests for files under a path.
package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	if len(os.Args) != 2 {
		sys.Usage("PATH")
	}
	root := os.Args[1]
	n := 0
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		info, err := d.Info()
		if err != nil || !info.Mode().IsRegular() {
			return nil
		}
		if n >= 200 {
			fmt.Println("truncated: yes")
			return filepath.SkipAll
		}
		sum, err := hashFile(path)
		if err != nil {
			fmt.Printf("skip %s: %v\n", path, err)
			return nil
		}
		fmt.Printf("%x  %s\n", sum, path)
		n++
		return nil
	})
	sys.Die(err)
	if n == 0 {
		fmt.Println("no files")
	}
}

func hashFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
