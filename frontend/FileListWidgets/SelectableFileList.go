package filelistwidgets

import (
	appstate "LiveBuilder/AppState"
	filesystem "LiveBuilder/Filesystem"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ListItem struct {
	IsCategory bool
	IsExpanded bool
	Category   string
	FileEntry  *filesystem.DirectoryEntry
	Depth      int
}

type FileListContainer struct {
	selectedFiles    map[string]filesystem.DirectoryEntry
	fileManager      *filesystem.FileManager
	directoryEntries []filesystem.DirectoryEntry
	fileView         *widget.Label
	fileViewHeader   *widget.Label
	list             *widget.List
	listItems        []ListItem
	categoryFiles    map[string][]filesystem.DirectoryEntry
}

func NewFileListContainer(filesystem_identifier string) *FileListContainer {
	fm := filesystem.GetFileManager()
	selectFileMap := appstate.GetGlobalState().GetDirectoryEntryMap(filesystem_identifier)
	flc := &FileListContainer{
		selectedFiles:    selectFileMap,
		fileManager:      fm,
		directoryEntries: fm.GetFileSystem(filesystem_identifier),
		fileView:         widget.NewLabel(""),
		fileViewHeader:   widget.NewLabel("Select An Item From The List"),
		categoryFiles:    make(map[string][]filesystem.DirectoryEntry),
	}
	flc.organizeByCategoriesAndTags()
	flc.buildListItems()
	return flc
}

func (self *FileListContainer) organizeByCategoriesAndTags() {
	// Clear existing categories
	self.categoryFiles = make(map[string][]filesystem.DirectoryEntry)

	// Organize files by their tags
	for _, entry := range self.directoryEntries {
		if len(entry.MetaData.Tags) == 0 {
			// Files without tags go in "Uncategorized"
			self.categoryFiles["Uncategorized"] = append(self.categoryFiles["Uncategorized"], entry)
		} else {
			// Add file to each of its tag categories
			for _, tag := range entry.MetaData.Tags {
				self.categoryFiles[tag] = append(self.categoryFiles[tag], entry)
			}
		}
	}
}

func (self *FileListContainer) buildListItems() {
	self.listItems = []ListItem{}

	// Add categories and their files
	for category, files := range self.categoryFiles {
		// Add category item
		categoryItem := ListItem{
			IsCategory: true,
			IsExpanded: false,
			Category:   category,
			FileEntry:  nil,
			Depth:      0,
		}
		self.listItems = append(self.listItems, categoryItem)

		// Add files under this category (initially hidden)
		for _, file := range files {
			fileItem := ListItem{
				IsCategory: false,
				IsExpanded: false,
				Category:   category,
				FileEntry:  &file,
				Depth:      1,
			}
			self.listItems = append(self.listItems, fileItem)
		}
	}
}

func (self *FileListContainer) toggleCategory(categoryName string) {
	// Find the category and toggle its expansion state
	for i := range self.listItems {
		if self.listItems[i].IsCategory && self.listItems[i].Category == categoryName {
			self.listItems[i].IsExpanded = !self.listItems[i].IsExpanded
			break
		}
	}

	// Refresh the list to show/hide files
	if self.list != nil {
		self.list.Refresh()
	}
}

func (self *FileListContainer) getVisibleItems() []ListItem {
	var visible []ListItem

	for _, item := range self.listItems {
		if item.IsCategory {
			// Always show categories
			visible = append(visible, item)
		} else {
			// Only show files if their category is expanded
			for _, categoryItem := range self.listItems {
				if categoryItem.IsCategory &&
					categoryItem.Category == item.Category &&
					categoryItem.IsExpanded {
					visible = append(visible, item)
					break
				}
			}
		}
	}

	return visible
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
	list := widget.NewList(
		func() int {
			return len(self.getVisibleItems())
		},
		func() fyne.CanvasObject {
			return NewSelectableListItem("", self)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			visibleItems := self.getVisibleItems()
			if id >= len(visibleItems) {
				return
			}

			listItem := visibleItems[id]
			selectableItem := item.(*FileListItem)

			if listItem.IsCategory {
				selectableItem.SetAsCategory(listItem.Category, listItem.IsExpanded)
			} else {
				selectableItem.SetAsFile(*listItem.FileEntry, listItem.Depth)
			}
		},
	)

	list.OnUnselected = func(id widget.ListItemID) {
		self.fileViewHeader.SetText("Select An Item From The List")
		self.fileView.SetText("")
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
