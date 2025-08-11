package appstate

import (
	filesystem "LiveBuilder/Filesystem"
	"sync"
)

var lock = &sync.Mutex{}

type selectedFileMap map[string]filesystem.DirectoryEntry

type State struct {
	selectedPackages map[string]selectedFileMap
	LBConfigCMD      string
}

var globalState *State

func GetGlobalState() *State {
	if globalState == nil {
		lock.Lock()
		defer lock.Unlock()
		if globalState == nil {
			globalState = &State{
				selectedPackages: make(map[string]selectedFileMap),
			}
		}
	}
	return globalState
}

func (state *State) GetDirectoryEntryMap(identifier string) selectedFileMap {
	fileMap, ok := state.selectedPackages[identifier]
	if !ok {
		fileMap = make(selectedFileMap)
		state.selectedPackages[identifier] = fileMap
	}
	return state.selectedPackages[identifier]
}
