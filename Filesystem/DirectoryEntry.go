package filesystem

import (
	"io/fs"
	"log"
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
	var customEntries []DirectoryEntry

	err := filepath.WalkDir(dirPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v\n", path, err)
			return nil // Continue walking despite errors
		}

		// Skip .meta.json files
		if strings.HasSuffix(entry.Name(), ".meta.json") {
			return nil
		}

		if entry.IsDir() {
			return nil
		}

		// Skip the root directory itself
		if path == dirPath {
			return nil
		}

		customEntry, err := NewCustomDirEntryFromPath(path, entry)
		if err != nil {
			log.Printf("Failed to create directory entry for %s: %v\n", path, err)
			return nil // Continue walking
		}

		metaData, err := LoadFileMetadata(customEntry.fullPath)
		if err != nil {
			log.Printf("Failed to load meta data for file %s, with error %v\n", customEntry.fullPath, err)
		}
		customEntry.MetaData = metaData

		customEntries = append(customEntries, customEntry)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return customEntries, nil
}
func NewCustomDirEntryFromPath(fullPath string, entry fs.DirEntry) (DirectoryEntry, error) {
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
