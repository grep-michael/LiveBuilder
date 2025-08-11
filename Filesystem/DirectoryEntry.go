package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
)

type DirectoryEntry struct {
	name     string
	fullPath string
	fileInfo fs.FileInfo
}

func (c *DirectoryEntry) Name() string {
	return c.name
}
func (c *DirectoryEntry) IsDir() bool {
	return c.fileInfo.IsDir()
}
func (c *DirectoryEntry) FullPath() string {
	return c.fullPath
}

func ScanDirectory(dirPath string) ([]DirectoryEntry, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var customEntries []DirectoryEntry

	for _, entry := range entries {
		customEntry, err := NewCustomDirEntryFromEntry(entry, dirPath)
		if err != nil {
			continue
		}
		customEntries = append(customEntries, customEntry)
	}

	return customEntries, nil
}

func NewCustomDirEntryFromEntry(entry fs.DirEntry, basePath string) (DirectoryEntry, error) {
	fullPath := filepath.Join(basePath, entry.Name())
	info, err := entry.Info()
	if err != nil {
		return DirectoryEntry{}, err
	}

	return DirectoryEntry{
		name:     entry.Name(),
		fullPath: fullPath,
		fileInfo: info,
	}, nil
}
