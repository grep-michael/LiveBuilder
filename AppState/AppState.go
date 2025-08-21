package appstate

import (
	filesystem "LiveBuilder/Filesystem"
	"strings"
	"sync"
)

type LBConfig struct {
	ISOVolume      string
	ISOPublisher   string
	ISOApplication string
	ISOImageName   string
}

func initalLBconfig() *LBConfig {
	return &LBConfig{
		ISOVolume:      "DefaultVolume",
		ISOPublisher:   "DefaultPublisher",
		ISOApplication: "DefaultApplication",
		ISOImageName:   "DefaultImage",
	}
}

var lock = &sync.Mutex{}

type selectedFileMap map[string]filesystem.DirectoryEntry

type State struct {
	selectedFiles map[string]selectedFileMap
	WriteLock     sync.Mutex
	LBcfg         *LBConfig
}

var globalState *State

func GetGlobalState() *State {
	if globalState == nil {
		lock.Lock()
		defer lock.Unlock()
		if globalState == nil {
			globalState = &State{
				selectedFiles: make(map[string]selectedFileMap),
				LBcfg:         initalLBconfig(),
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
func (state *State) setISOStringField(field *string, value string) {
	state.WriteLock.Lock()
	defer state.WriteLock.Unlock()
	value = strings.ReplaceAll(value, " ", "_")
	*field = value
}
func (state *State) SetISOVolumeName(name string) {
	state.setISOStringField(&state.LBcfg.ISOVolume, name)
}
func (state *State) SetISOPublisher(name string) {
	state.setISOStringField(&state.LBcfg.ISOPublisher, name)
}
func (state *State) SetISOApplication(name string) {
	state.setISOStringField(&state.LBcfg.ISOApplication, name)
}
func (state *State) SetISOImageName(name string) {
	state.setISOStringField(&state.LBcfg.ISOImageName, name)
}
func (state *State) ISOVolumeName() string {
	return state.LBcfg.ISOVolume
}
func (state *State) ISOPublisher() string {
	return state.LBcfg.ISOPublisher
}
func (state *State) ISOApplication() string {
	return state.LBcfg.ISOApplication
}
func (state *State) ISOImageName() string {
	return state.LBcfg.ISOImageName
}
