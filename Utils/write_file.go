package Utils

import (
	"log"
	"os"
	"path/filepath"
)

func WriteFile(path, content string, perm os.FileMode) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), perm); err != nil {
		log.Fatalf("write %s: %v", path, err)
	}
}
