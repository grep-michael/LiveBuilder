package filelistwidgets

import (
	appstate "LiveBuilder/AppState"
	filesystem "LiveBuilder/Filesystem"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type FileListContainer struct {
	selectedFiles    map[string]filesystem.DirectoryEntry
	fileManager      *filesystem.FileManager
	directoryEntries []filesystem.DirectoryEntry
	fileView         *widget.Label
	fileViewHeader   *widget.Label
	list             *widget.List
}

func NewFileListContainer(filesystem_identifier string) *FileListContainer {
	fm := filesystem.GetFileManager()
	selectFileMap := appstate.GetGlobalState().GetDirectoryEntryMap(filesystem_identifier)
	return &FileListContainer{
		selectedFiles:    selectFileMap,
		fileManager:      fm,
		directoryEntries: fm.GetFileSystem(filesystem_identifier),
		fileView:         widget.NewLabel(""),
		fileViewHeader:   widget.NewLabel("Select An Item From The List"),
	}
}
func (self *FileListContainer) isFileSelected(fileEntry filesystem.DirectoryEntry) bool {
	_, ok := self.selectedFiles[fileEntry.Name()]
	return ok
}
func (self *FileListContainer) addSelectedFile(fileEntry filesystem.DirectoryEntry) {
	if !self.isFileSelected(fileEntry) {
		self.selectedFiles[fileEntry.Name()] = fileEntry
	}
}
func (self *FileListContainer) removeSelectedFile(fileEntry filesystem.DirectoryEntry) {
	_, ok := self.selectedFiles[fileEntry.Name()]
	if ok {
		delete(self.selectedFiles, fileEntry.Name())
	}
}
func (self *FileListContainer) toggleFileSelection(fileEntry filesystem.DirectoryEntry) {
	if self.isFileSelected(fileEntry) {
		self.removeSelectedFile(fileEntry)
	} else {
		self.addSelectedFile(fileEntry)
	}
}
func (self *FileListContainer) buildFileList() *widget.List {
	max_str := self.fileManager.GetCompareFileNameLengths(self.directoryEntries, func(a, b string) bool {
		return len(a) > len(b)
	})

	list := widget.NewList(
		func() int {
			return len(self.directoryEntries)
		},
		func() fyne.CanvasObject {
			return NewSelectableListItem(max_str, self)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			fileName := self.directoryEntries[id].Name()
			selectableItem := item.(*FileListItem)

			selectableItem.packageName = fileName
			selectableItem.fileEntry = self.directoryEntries[id]
			selectableItem.label.SetText(fileName)
			selectableItem.updateIcon()

		},
	)

	list.OnUnselected = func(id widget.ListItemID) {
		self.fileView.SetText("Select An Item From The List")
	}
	self.list = list
	return list
}
func (self *FileListContainer) buildFileContentView() *container.Split {
	scroll := container.NewScroll(self.fileView)
	scroll.SetMinSize(fyne.NewSize(200, 400))

	vbox := container.NewVSplit(
		self.fileViewHeader,
		scroll,
	)

	// Set the split to give minimal space to the top (header)
	headerHeight := self.fileViewHeader.MinSize().Height
	totalHeight := headerHeight + 400 // approximate total
	offset := float64(headerHeight) / float64(totalHeight)
	vbox.SetOffset(offset)

	return vbox
}
func (self *FileListContainer) GetContainer() fyne.CanvasObject {
	list := self.buildFileList()
	hsplit := container.NewHSplit(list, self.buildFileContentView())
	hsplit.Refresh()
	return hsplit
}
