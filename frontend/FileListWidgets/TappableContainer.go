package filelistwidgets

import (
	filesystem "LiveBuilder/Filesystem"
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Custom tappable container for list items
type FileListItem struct {
	widget.BaseWidget
	fileListContainer *FileListContainer
	icon              *widget.Icon
	label             *widget.Label
	container         *fyne.Container
	fileEntry         *filesystem.DirectoryEntry
	isCategory        bool
	categoryName      string
	isExpanded        bool
	depth             int
}

func NewSelectableListItem(fileName string, flc *FileListContainer) *FileListItem {
	icon := widget.NewIcon(theme.DocumentIcon())
	label := widget.NewLabel(fileName)
	container := container.NewHBox(icon, label)

	item := &FileListItem{
		fileListContainer: flc,
		icon:              icon,
		label:             label,
		container:         container,
		isCategory:        false,
		depth:             0,
	}

	item.ExtendBaseWidget(item)
	return item
}

// Implement the required CreateRenderer method
func (self *FileListItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(self.container)
}

func (self *FileListItem) SetAsCategory(categoryName string, isExpanded bool) {
	self.isCategory = true
	self.categoryName = categoryName
	self.isExpanded = isExpanded
	self.fileEntry = nil
	self.depth = 0

	// Set category icon based on expansion state
	if isExpanded {
		self.icon.SetResource(theme.MenuExpandIcon())
	} else {
		self.icon.SetResource(theme.MenuDropDownIcon())
	}

	// Set category label with bold formatting
	self.label.SetText(fmt.Sprintf("%s", categoryName))
	self.label.TextStyle = fyne.TextStyle{Bold: true}
}

func (self *FileListItem) SetAsFile(fileEntry filesystem.DirectoryEntry, depth int) {
	self.isCategory = false
	self.fileEntry = &fileEntry
	self.depth = depth
	self.categoryName = ""

	// Add indentation for nested files
	indent := strings.Repeat("    ", depth)
	self.label.SetText(fmt.Sprintf("%s%s", indent, fileEntry.Name()))
	self.label.TextStyle = fyne.TextStyle{Bold: false}

	self.updateFileIcon()
}

func (self *FileListItem) updateFileIcon() {
	if self.fileEntry == nil {
		return
	}

	if self.fileListContainer.isFileSelected(*self.fileEntry) {
		self.icon.SetResource(theme.ConfirmIcon())
	} else {
		// Set icon based on file type
		switch self.fileEntry.MetaData.FileType {
		case "script", "python", "go", "javascript":
			self.icon.SetResource(theme.ComputerIcon())
		case "config", "json", "xml", "yaml":
			self.icon.SetResource(theme.SettingsIcon())
		case "log", "txt":
			self.icon.SetResource(theme.DocumentIcon())
		default:
			self.icon.SetResource(theme.FileIcon())
		}
	}
}

func (self *FileListItem) getFileContents() string {
	if self.fileEntry == nil {
		return ""
	}

	bytes, err := os.ReadFile(self.fileEntry.FullPath())
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}
	return string(bytes)
}

func (self *FileListItem) Tapped(_ *fyne.PointEvent) {
	if self.isCategory {
		// Toggle category expansion
		self.fileListContainer.toggleCategory(self.categoryName)
	} else if self.fileEntry != nil {
		// Show file content
		text := self.getFileContents()

		// Create a nice header with file info
		header := fmt.Sprintf("%s\n", self.fileEntry.FullPath())
		if self.fileEntry.MetaData.Description != "" {
			header += fmt.Sprintf("Description: %s\n", self.fileEntry.MetaData.Description)
		}
		if len(self.fileEntry.MetaData.Tags) > 0 {
			header += fmt.Sprintf("Tags: %s\n", strings.Join(self.fileEntry.MetaData.Tags, ", "))
		}
		if self.fileEntry.MetaData.FileType != "" {
			header += fmt.Sprintf("Type: %s\n", self.fileEntry.MetaData.FileType)
		}
		header += strings.Repeat("-", 50)

		self.fileListContainer.fileViewHeader.SetText(header)
		self.fileListContainer.fileView.SetText(text)
	}
}

func (self *FileListItem) TappedSecondary(_ *fyne.PointEvent) {
	if !self.isCategory && self.fileEntry != nil {
		// Right-click or Ctrl+Click - toggle file selection
		self.fileListContainer.toggleFileSelection(*self.fileEntry)
		self.updateFileIcon()
	}
}

func (self *FileListItem) SetPackageName(packageName string) {
	if !self.isCategory {
		self.label.SetText(packageName)
		self.updateFileIcon()
	}
}
