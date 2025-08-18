package buildwindow

import (
	buildmanager "LiveBuilder/BuildManager"
	logger "LiveBuilder/BuildManager/Logger"

	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type BuildWindow struct {
	window            fyne.Window
	buildPath         string
	selectedPathLabel *widget.Label
	buildStatusLabel  *widget.Label
	logWidget         *logger.LogWidget
	logContent        strings.Builder
	buildManager      *buildmanager.BuildManager
	logScroll         *container.Scroll
	//buildLogText      *widget.RichText
	//livebuilder       *execution.LiveBuilder
}

func NewBuildWindow(window fyne.Window) *fyne.Container {
	build_window := BuildWindow{
		window:            window,
		selectedPathLabel: widget.NewLabel("Select folder"),
		buildStatusLabel:  widget.NewLabel("Statuses"),
	}
	//build_window.buildLogText = widget.NewRichTextFromMarkdown("Build Log will appear here...")
	//build_window.buildLogText.Wrapping = fyne.TextWrapWord
	build_window.logWidget = logger.NewLogWidget(200)

	//build_window.logScroll = container.NewScroll(build_window.logWidget)
	//build_window.logScroll.SetMinSize(fyne.NewSize(600, 400))

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
		if !update.Append {
			self.logWidget.Clear()
		}
		fyne.Do(func() {
			self.logWidget.AppendLine(update.Message)
		})

	}
	//for update := range subscriber {
	//	if update.Append {
	//		self.logContent.WriteString(update.Message)
	//		self.logContent.WriteString("\n")
	//	} else {
	//		self.logContent.Reset()
	//		self.logContent.WriteString(update.Message)
	//	}
	//	fyne.DoAndWait(func() {
	//		self.buildLogText.ParseMarkdown(self.logContent.String())
	//		self.logScroll.ScrollToBottom()
	//	})
	//}
}
func (self *BuildWindow) buildMainBuildArea() *fyne.Container {
	buildButton := widget.NewButton("Execute Live Build", func() {
		self.logContent.Reset()
		self.buildStatusLabel.SetText("Building...")

		go func() {
			self.buildManager.Build(self.buildPath)
			log.Println("all building done, display final message")
			fyne.Do(func() {
				self.buildStatusLabel.SetText("Building Finished!")
			})
		}()
	})

	//hbox := container.NewVBox(buildButton, self.logScroll, self.buildStatusLabel)
	hbox := container.NewVBox(buildButton, self.logWidget.GetWidget(), self.buildStatusLabel)
	return hbox
}
