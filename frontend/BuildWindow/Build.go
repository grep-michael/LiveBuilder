package buildwindow

import (
	buildmanager "LiveBuilder/BuildManager"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
	"strings"
)

type BuildWindow struct {
	window            fyne.Window
	buildPath         string
	selectedPathLabel *widget.Label
	buildLogText      *widget.RichText
	logContent        strings.Builder
	logScroll         *container.Scroll
	//livebuilder       *execution.LiveBuilder
	buildManager *buildmanager.BuildManager
}

func NewBuildWindow(window fyne.Window) *fyne.Container {
	build_window := BuildWindow{
		window:            window,
		selectedPathLabel: widget.NewLabel("Select folder"),
	}
	build_window.buildLogText = widget.NewRichTextFromMarkdown("Build Log will appear here...")
	build_window.buildLogText.Wrapping = fyne.TextWrapWord

	build_window.logScroll = container.NewScroll(build_window.buildLogText)
	build_window.logScroll.SetMinSize(fyne.NewSize(600, 400))

	build_window.buildManager = buildmanager.NewBuilder()

	go build_window.startLogSubscriber()

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
			self.buildPath = folderPath
		}, self.window)
	})

	hbox := container.NewVBox(choose_folder_btn, self.selectedPathLabel)

	return hbox
}
func (self *BuildWindow) startLogSubscriber() {
	subscriber := self.buildManager.GetSubscriber()

	for update := range subscriber {
		if update.Append {
			self.logContent.WriteString(update.Message)
			self.logContent.WriteString("\n")
		} else {
			self.logContent.Reset()
			self.logContent.WriteString(update.Message)
		}
		fyne.DoAndWait(func() {
			self.buildLogText.ParseMarkdown(self.logContent.String())
			self.logScroll.ScrollToBottom()
		})
	}
}
func (self *BuildWindow) buildMainBuildArea() *fyne.Container {
	buildButton := widget.NewButton("Execute Live Build", func() {
		self.logContent.Reset()
		self.buildLogText.ParseMarkdown("Building...")

		go func() {
			self.buildManager.Build(self.buildPath)
			log.Println("all building done, display final message")

			fyne.DoAndWait(func() {
				self.buildLogText.ParseMarkdown("Build Completed")
			})
			self.logContent.Reset()
		}()
	})

	hbox := container.NewVBox(buildButton, self.logScroll)

	return hbox
}
