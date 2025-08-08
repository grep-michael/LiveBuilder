package frontend

import (
	"LiveBuilder/backend"
	//"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strings"
)

type PackageListContainer struct {
	selectedPackages []string
	fm               *backend.EmbeddedFileManager
	fileView         *widget.Label
}

func NewPackageListContainer() *PackageListContainer {
	fm := backend.GetFileManager()
	return &PackageListContainer{
		fm:       fm,
		fileView: widget.NewLabel("Select An Item From The List"),
	}
}

func (plc *PackageListContainer) buildPackageListWidget() *widget.List {
	maxLen := plc.fm.GetLongestString(plc.fm.GetPackagelist())
	max_str := strings.Repeat("a", maxLen)
	list := widget.NewList(
		func() int {
			return len(plc.fm.GetPackagelist())
		},
		func() fyne.CanvasObject {
			icon := widget.NewIcon(theme.DocumentIcon())
			label := widget.NewLabel(max_str)
			container := container.NewHBox(icon, label)
			return container
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			text := plc.fm.GetPackagelist()[id].Name()
			container := item.(*fyne.Container)
			label := container.Objects[1].(*widget.Label)
			label.SetText(text)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		fileName := plc.fm.GetPackagelist()[id].Name()
		text := plc.fm.GetTextFromFile(fileName)
		plc.fileView.SetText(text)
	}
	list.OnUnselected = func(id widget.ListItemID) {
		plc.fileView.SetText("Select An Item From The List")
	}
	return list
}

func (plc *PackageListContainer) GetContainer() fyne.CanvasObject {
	list := plc.buildPackageListWidget()

	hsplit := container.NewHSplit(list, container.NewCenter(plc.fileView))
	hsplit.Refresh()
	return hsplit
}
