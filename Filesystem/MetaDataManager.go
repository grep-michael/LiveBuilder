package filesystem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileMetadata struct {
	InstallPath string   `json:"install_path"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	FileType    string   `json:"file_type"`
}

// LoadFileMetadata loads metadata from a sidecar .meta.json file
// If the file doesn't exist, creates a default one
// If it exists but has missing fields, fills them with defaults
func LoadFileMetadata(filePath string) (FileMetadata, error) {
	metaPath := filePath + ".meta.json"

	// Get default metadata
	defaults := getDefaultMetadata()

	// Try to load existing metadata
	if data, err := os.ReadFile(metaPath); err == nil {
		// File exists, try to unmarshal
		var existingMeta FileMetadata
		if err := json.Unmarshal(data, &existingMeta); err == nil {
			// Fill missing fields with defaults
			if existingMeta.InstallPath == "" {
				existingMeta.InstallPath = defaults.InstallPath
			}
			if len(existingMeta.Tags) == 0 {
				existingMeta.Tags = defaults.Tags
			}
			if existingMeta.Description == "" {
				existingMeta.Description = defaults.Description
			}
			if existingMeta.FileType == "" {
				existingMeta.FileType = defaults.FileType
			}

			// Save only if we added missing fields
			needsSave := false
			data2, _ := json.Marshal(existingMeta)
			if string(data) != string(data2) {
				needsSave = true
			}

			if needsSave {
				if err := SaveFileMetadata(filePath, existingMeta); err != nil {
					return existingMeta, fmt.Errorf("failed to save updated metadata: %v", err)
				}
			}

			return existingMeta, nil
		}
	}

	// File doesn't exist or couldn't be parsed, create default
	if err := SaveFileMetadata(filePath, defaults); err != nil {
		return defaults, fmt.Errorf("failed to create metadata file: %v", err)
	}

	return defaults, nil
}

// SaveFileMetadata saves metadata to a sidecar .meta.json file
func SaveFileMetadata(filePath string, meta FileMetadata) error {
	metaPath := filePath + ".meta.json"

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(metaPath, data, 0644)
}

// getDefaultMetadata returns a standard set of default values
func getDefaultMetadata() FileMetadata {
	return FileMetadata{
		InstallPath: ".",
		Tags:        []string{""},
		Description: "",
		FileType:    "",
	}
}

// GetAllFilesWithMetadata scans a directory and returns files with their metadata
func GetAllFilesWithMetadata(dir string) (map[string]FileMetadata, error) {
	files := make(map[string]FileMetadata)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || strings.HasSuffix(entry.Name(), ".meta.json") {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		meta, err := LoadFileMetadata(filePath)
		if err != nil {
			continue // Skip files with metadata errors
		}

		files[filePath] = meta
	}

	return files, nil
}

// FilterFilesByTag returns files that have any of the specified tags
func FilterFilesByTag(filesWithMeta map[string]FileMetadata, tags ...string) map[string]FileMetadata {
	filtered := make(map[string]FileMetadata)

	for filePath, meta := range filesWithMeta {
		for _, tag := range tags {
			for _, fileTag := range meta.Tags {
				if strings.EqualFold(tag, fileTag) {
					filtered[filePath] = meta
					goto nextFile
				}
			}
		}
	nextFile:
	}

	return filtered
}
