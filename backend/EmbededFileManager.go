package backend

import (
	"embed"
	"fmt"
	"io/fs"
	"sync"
)

var lock = &sync.Mutex{}

//go:embed all:staticfiles
var embeddedFS embed.FS

type EmbeddedFileManager struct {
	packageLists []fs.DirEntry
	scriptsList  []fs.DirEntry
	fs           embed.FS
}

var globalInstance *EmbeddedFileManager

func GetFileManager() *EmbeddedFileManager {
	if globalInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if globalInstance == nil {
			globalInstance = &EmbeddedFileManager{
				fs: embeddedFS,
			}
			globalInstance.RegisterEmbedded()
		}
	}
	return globalInstance
}

func (app *EmbeddedFileManager) RegisterEmbedded() {
	packlists, err := fs.ReadDir(app.fs, "staticfiles/PackageLists")
	if err != nil {
		fmt.Println(err)
		return
	}
	app.packageLists = packlists

	scripts, err := fs.ReadDir(app.fs, "staticfiles/Scripts")
	if err != nil {
		fmt.Println(err)
		return
	}
	app.scriptsList = scripts
}
func (app *EmbeddedFileManager) GetLongestString(list []fs.DirEntry) int {
	var max int
	for _, value := range list {
		length := len(value.Name())
		if length > max {
			//print(value.Name())
			max = length
		}
	}
	fmt.Println(max)
	return max
}

func (app *EmbeddedFileManager) GetTextFromFile(fileName string) string {
	fmt.Println(fileName)
	content, err := fs.ReadFile(embeddedFS, fileName)
	if err != nil {
		return err.Error()
	}
	return string(content)
}

func (app *EmbeddedFileManager) GetPackagelist() []fs.DirEntry {
	return app.packageLists
}
func (app *EmbeddedFileManager) GetScriptList() []fs.DirEntry {
	return app.scriptsList
}
