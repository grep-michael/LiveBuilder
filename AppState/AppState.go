package appstate

import (
	filesystem "LiveBuilder/Filesystem"
	"sync"
)

var lock = &sync.Mutex{}

type selectedFileMap map[string]filesystem.DirectoryEntry

type State struct {
	selectedFiles map[string]selectedFileMap
	LBConfigCMD   string
}

var globalState *State

func GetGlobalState() *State {
	if globalState == nil {
		lock.Lock()
		defer lock.Unlock()
		if globalState == nil {
			globalState = &State{
				selectedFiles: make(map[string]selectedFileMap),
			}
		}
	}
	return globalState
}

func (state *State) GetDirectoryEntryMap(identifier string) selectedFileMap {
	fileMap, ok := state.selectedFiles[identifier]
	if !ok {
		fileMap = make(selectedFileMap)
		state.selectedFiles[identifier] = fileMap
	}
	return state.selectedFiles[identifier]
}
