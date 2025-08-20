package filesystem

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	EMBEDDED_FS_ROOT = "staticfiles"
)

//go:embed all:staticfiles
var embeddedFiles embed.FS

// getAppDataDir returns the appropriate directory based on OS
func GetAppDataDir() (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		// Use APPDATA environment variable, fallback to USERPROFILE
		appData := os.Getenv("APPDATA")
		if appData == "" {
			userProfile := os.Getenv("USERPROFILE")
			if userProfile == "" {
				return "", fmt.Errorf("cannot determine user directory on Windows")
			}
			baseDir = filepath.Join(userProfile, "AppData", "Roaming")
		} else {
			baseDir = appData
		}
	case "darwin":
		// macOS: ~/Library/Application Support
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot get home directory: %v", err)
		}
		baseDir = filepath.Join(homeDir, "Library", "Application Support")
	case "linux":
		// Linux: ~/.config (XDG Base Directory Specification)
		configDir := os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("cannot get home directory: %v", err)
			}
			baseDir = filepath.Join(homeDir, ".config")
		} else {
			baseDir = configDir
		}
	default:
		// Fallback for other Unix-like systems
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot get home directory: %v", err)
		}
		baseDir = filepath.Join(homeDir, "."+APPNAME)
	}
	dir := filepath.Join(baseDir, APPNAME)
	os.MkdirAll(dir, 0777)
	return dir, nil
}

// extractEmbeddedFiles extracts all embedded files to the target directory
func extractEmbeddedFiles(targetDir string) error {
	// Create target directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %v", err)
	}

	// Walk through all embedded files
	return fs.WalkDir(embeddedFiles, EMBEDDED_FS_ROOT, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip the root directory
		if path == EMBEDDED_FS_ROOT {
			return nil
		}

		relativePath, err := filepath.Rel(EMBEDDED_FS_ROOT, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %v", err)
		}

		// Create the full target path
		targetPath := filepath.Join(targetDir, relativePath)

		if _, err := os.Stat(targetPath); err == nil {
			//if file exists continue past this file
			log.Printf("Skipping file %s\n", targetPath)
			return nil
		}

		if d.IsDir() {
			// Create directory
			log.Printf("Creating directory: %s\n", targetPath)
			return os.MkdirAll(targetPath, 0755)
		}

		// Read embedded file content
		content, err := embeddedFiles.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %v", path, err)
		}

		// Create target file
		log.Printf("Extracting file: %s\n", targetPath)
		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %v", targetPath, err)
		}

		return nil
	})
}
