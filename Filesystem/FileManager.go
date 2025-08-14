package filesystem

import (
	"log"
	"path/filepath"
	"sync"
)

const (
	APPNAME              = "LiveBuidler"
	PACKAGE_DIR_ID       = "PackageLists"
	CUSTOMFILES_DIR_ID   = "CustomFiles"
	LBCONFIG_TEMPLATE_ID = "lbconfig_template.template"
	SPLASH_SCREENS_ID    = "SplashScreens"
)

var lock = &sync.Mutex{}

type FileManager struct {
	appDriectory string
	fileSystems  map[string][]DirectoryEntry
}

var globalInstance *FileManager

func GetFileManager() *FileManager {
	if globalInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if globalInstance == nil {
			globalInstance = &FileManager{
				fileSystems: make(map[string][]DirectoryEntry),
			}
			globalInstance.InializeFilesystem()
		}
	}
	return globalInstance
}

func (self *FileManager) InializeFilesystem() {
	appDir, err := GetAppDataDir()
	if err != nil {
		log.Fatalf("Error getting app data directory: %v", err)
	}
	self.appDriectory = appDir
	if err := extractEmbeddedFiles(appDir); err != nil {
		log.Fatalf("Error extracting embedded files: %v", err)
	}
	self.buildFilesystemMap()
}

func (self *FileManager) buildFilesystemMap() {
	//var err error

	package_path := filepath.Join(self.GetAppDataDir(), PACKAGE_DIR_ID)
	log.Printf("Packages path: %s\n", package_path)
	self.fileSystems[PACKAGE_DIR_ID], _ = ScanDirectory(package_path)

	scripts_path := filepath.Join(self.GetAppDataDir(), CUSTOMFILES_DIR_ID)
	log.Printf("Scripts path: %s\n", scripts_path)
	self.fileSystems[CUSTOMFILES_DIR_ID], _ = ScanDirectory(scripts_path)

	splash_path := filepath.Join(self.GetAppDataDir(), SPLASH_SCREENS_ID)
	log.Printf("splash path: %s\n", scripts_path)
	self.fileSystems[SPLASH_SCREENS_ID], _ = ScanDirectory(splash_path)

}
func (self *FileManager) GetFileSystem(fs_identifier string) []DirectoryEntry {
	return self.fileSystems[fs_identifier]
}

func (self *FileManager) GetAppDataDir() string {
	return self.appDriectory
}

func (self *FileManager) GetCompareFileNameLengths(_filesystem []DirectoryEntry, compare func(string, string) bool) string {

	if len(_filesystem) == 0 {
		return ""
	}
	result := _filesystem[0].Name()

	for _, filename := range _filesystem[1:] {
		if compare(filename.Name(), result) {
			result = filename.Name()
		}
	}

	return result

}
