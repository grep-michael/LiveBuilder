package filesystem

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type DirectoryEntry struct {
	name     string
	fullPath string
	fileInfo fs.FileInfo
	MetaData FileMetadata
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

		if strings.HasSuffix(entry.Name(), ".meta.json") {
			continue
		}

		customEntry, err := NewCustomDirEntryFromEntry(entry, dirPath)
		if err != nil {
			continue
		}
		metaData, err := LoadFileMetadata(customEntry.fullPath)
		if err != nil {
			log.Printf("Failed to load meta data for file %s, with error %v\n", customEntry.fullPath, err)
		}
		customEntry.MetaData = metaData
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
