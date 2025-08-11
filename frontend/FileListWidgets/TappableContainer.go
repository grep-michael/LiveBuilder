package filelistwidgets

import (
	filesystem "LiveBuilder/Filesystem"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Custom tappable container for list items
type FileListItem struct {
	widget.BaseWidget
	packageName       string
	fileListContainer *FileListContainer
	icon              *widget.Icon
	label             *widget.Label
	container         *fyne.Container
	fileEntry         filesystem.DirectoryEntry
}

func NewSelectableListItem(fileName string, flc *FileListContainer) *FileListItem {
	icon := widget.NewIcon(theme.DocumentIcon())
	label := widget.NewLabel(fileName)
	container := container.NewHBox(icon, label)

	item := &FileListItem{
		packageName:       fileName,
		fileListContainer: flc,
		icon:              icon,
		label:             label,
		container:         container,
	}

	item.ExtendBaseWidget(item)
	item.updateIcon()
	return item
}

// Implement the required CreateRenderer method
func (self *FileListItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(self.container)
}

func (self *FileListItem) updateIcon() {
	if self.fileListContainer.isFileSelected(self.fileEntry) {
		self.icon.SetResource(theme.ConfirmIcon())
	} else {
		self.icon.SetResource(theme.DocumentIcon())
	}
}

func (self *FileListItem) getFileContents() string {
	bytes, err := os.ReadFile(self.fileEntry.FullPath())
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return string(bytes)
}

func (self *FileListItem) Tapped(_ *fyne.PointEvent) {
	// Regular tap - show file content
	text := self.getFileContents()
	self.fileListContainer.fileView.SetText(text)
}

func (self *FileListItem) TappedSecondary(_ *fyne.PointEvent) {
	// Right-click or Ctrl+Click - toggle selection
	self.fileListContainer.toggleFileSelection(self.fileEntry)
	self.updateIcon()
}
func (self *FileListItem) SetPackageName(packageName string) {
	self.packageName = packageName
	self.label.SetText(packageName)
	self.updateIcon()
}
