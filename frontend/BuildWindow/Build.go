package buildwindow

import (
	buildmanager "LiveBuilder/BuildManager"
	logger "LiveBuilder/BuildManager/Logger"
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
	buildStatusLabel  *widget.Label
	logWidget         *logger.LogView
	logContent        strings.Builder
	buildManager      *buildmanager.BuildManager
	logScroll         *container.Scroll
	settings          *SettingsWidget
}

func NewBuildWindow(window fyne.Window) *fyne.Container {
	build_window := BuildWindow{
		window:            window,
		selectedPathLabel: widget.NewLabel("Select folder"),
		buildStatusLabel:  widget.NewLabel("Statuses"),
	}

	build_window.logWidget = logger.NewLogView(500)
	build_window.logScroll = container.NewScroll(build_window.logWidget)
	build_window.logScroll.SetMinSize(fyne.NewSize(600, 500))

	build_window.buildManager = buildmanager.NewBuilder()

	go build_window.startLogSubscriber()
	build_window.settings = NewRightClickWidget(window, build_window.buildManager)

	headers := build_window.buildFolderSelectionHeader()
	buildArea := build_window.buildMainBuildArea()
	return container.NewBorder(headers, nil, nil, nil, buildArea)
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
	vbox := container.NewVBox(choose_folder_btn, self.selectedPathLabel)

	hbox := container.NewBorder(nil, nil, nil, self.settings, vbox)
	return hbox
}

func (self *BuildWindow) startLogSubscriber() {
	subscriber := self.buildManager.GetSubscriber()
	for update := range subscriber {
		fyne.Do(func() {
			if !update.Append {
				self.logWidget.Clear()
			}
			self.logWidget.AppendLine(update.Message)
			self.logScroll.ScrollToBottom()
		})

	}
}

func (self *BuildWindow) buildMainBuildArea() *fyne.Container {
	buildButton := widget.NewButton("Execute Live Build", func() {
		self.logContent.Reset()
		self.buildStatusLabel.SetText("Building...")

		go func() {
			self.buildManager.InitializeBuild(self.buildPath)
			self.buildManager.BuildConditional(self.settings.selectedOption)
			self.buildStatusLabel.SetText("Building Finished!")
		}()
	})

	hbox := container.NewBorder(buildButton, self.buildStatusLabel, nil, nil, self.logScroll)
	return hbox
}
