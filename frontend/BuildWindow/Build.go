package buildwindow

import (
	execution "LiveBuilder/Execution"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
)

type BuildWindow struct {
	window            fyne.Window
	buildPath         string
	selectedPathLabel *widget.Label
	buildLogLabel     *widget.Label
	livebuilder       *execution.LiveBuilder
}

func NewBuildWindow(window fyne.Window) *fyne.Container {
	build_window := BuildWindow{
		window:            window,
		selectedPathLabel: widget.NewLabel("Select folder"),
		buildLogLabel:     widget.NewLabel("Build Log will be here"),
	}
	build_window.livebuilder = execution.NewLiveBuilder(build_window.buildLogLabel)

	filesectionHeader := build_window.buildFolderSelectionHeader()
	buildArea := build_window.buildMainBuildArea()
	return container.NewBorder(filesectionHeader, nil, nil, nil, buildArea)
}

func (self *BuildWindow) buildFolderSelectionHeader() *fyne.Container {
	choose_folder_btn := widget.NewButton("Choose Build Location", func() {
		dialog.ShowFolderOpen(func(folder fyne.ListableURI, err error) {
			if err != nil {
				log.Println("Error selecting folder:", err)
				return
			}
			if folder == nil {
				return
			}
			folderPath := folder.Path()
			self.selectedPathLabel.SetText("Selected: " + folderPath)
			self.livebuilder.SetWorkingDir(folderPath)
		}, self.window)
	})

	hbox := container.NewVBox(choose_folder_btn, self.selectedPathLabel)

	return hbox
}
func (self *BuildWindow) buildMainBuildArea() *fyne.Container {
	buildButton := widget.NewButton("Execute Live Build", func() {
		self.buildLogLabel.SetText("Building...")
		go func() {
			err := self.livebuilder.ConfigureLB()
			if err != nil {
				fmt.Println("build Error")
				fmt.Println(err)
				self.buildLogLabel.SetText("Error: " + err.Error())
			} else {
				self.buildLogLabel.SetText("Configure finished")
			}

		}()
	})

	scroll := container.NewScroll(self.buildLogLabel)
	scroll.SetMinSize(fyne.NewSize(200, 400))
	hbox := container.NewVBox(buildButton, scroll)

	return hbox
}
