package catalog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCommandsHaveMains(t *testing.T) {
	for _, c := range All() {
		path := filepath.Join("..", "..", "cmd", c.Name, "main.go")
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing %s", path)
		}
	}
}

func TestAlphabetical(t *testing.T) {
	all := All()
	if len(all) != 100 {
		t.Fatalf("len(All()) = %d, want 100", len(all))
	}
	for i := 1; i < len(all); i++ {
		if all[i].Name <= all[i-1].Name {
			t.Fatalf("commands not alphabetical at %s", all[i].Name)
		}
	}
}
