package livebuildconfig

import (
	appstate "LiveBuilder/AppState"
	filesystem "LiveBuilder/Filesystem"
	filelistwidgets "LiveBuilder/frontend/FileListWidgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type fields struct {
	label  string
	getter func() string
	setter func(string)
}

type LBConfigruationTab struct {
	headerGrid *fyne.Container
	fileList   *filelistwidgets.FileListContainer
}

func NewLBConfigurationTab() *LBConfigruationTab {
	cfg := &LBConfigruationTab{}
	cfg.headerGrid = cfg.buildHeaderGrid()
	cfg.fileList = cfg.buildFileList()
	return cfg
}

func (self *LBConfigruationTab) buildHeaderGrid() *fyne.Container {
	appstate := appstate.GetGlobalState()
	fields := []fields{
		{"ISO Volume", appstate.ISOVolumeName, appstate.SetISOVolumeName},
		{"ISO Publisher", appstate.ISOPublisher, appstate.SetISOPublisher},
		{"ISO Application", appstate.ISOApplication, appstate.SetISOApplication},
		{"ISO ImageName", appstate.ISOImageName, appstate.SetISOImageName},
	}

	var headers []fyne.CanvasObject
	var entries []fyne.CanvasObject

	for _, field := range fields {
		entry := widget.NewEntry()
		entry.SetPlaceHolder(field.getter())
		entry.OnChanged = field.setter
		entries = append(entries, entry)
		headers = append(headers, widget.NewLabel(field.label))

	}

	grid := container.NewGridWithColumns(4,
		append(headers, entries...)...,
	)
	return grid
}

func (self *LBConfigruationTab) buildFileList() *filelistwidgets.FileListContainer {
	return filelistwidgets.NewFileListContainer(filesystem.LBCONFIGS_DIR_ID)
}

func (self *LBConfigruationTab) GetContainer() *fyne.Container {
	return container.NewBorder(self.headerGrid, nil, nil, nil, self.fileList.GetContainer())
}
