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
	//"sync"
	//"time"
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
	build_window.logWidget = logger.NewLogView(200)

	build_window.logScroll = container.NewScroll(build_window.logWidget)
	build_window.logScroll.SetMinSize(fyne.NewSize(600, 600))

	build_window.buildManager = buildmanager.NewBuilder()

	go build_window.startLogSubscriber()
	//go build_window.startLogSubscriberWithBatching()
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

/*
	func (self *BuildWindow) startLogSubscriberWithBatching() {
		subscriber := self.buildManager.GetSubscriber()

		const (
			maxBatchSize  = 20
			flushInterval = 3 * time.Second
		)

		var pendingUpdates []buildmanager.LogUpdate
		var mu sync.Mutex
		var lastFlush time.Time = time.Now()

		flushUpdates := func() {
			mu.Lock()
			defer mu.Unlock()

			if len(pendingUpdates) == 0 {
				return
			}

			updates := make([]buildmanager.LogUpdate, len(pendingUpdates))
			copy(updates, pendingUpdates)
			pendingUpdates = pendingUpdates[:0]
			lastFlush = time.Now()

			for _, update := range updates {
				fyne.Do(func() {

				})
				if !update.Append {
					self.logWidget.Clear()
				}
				self.logWidget.AppendLine(update.Message)
			}
		}

		// Periodic flush goroutine
		go func() {
			ticker := time.NewTicker(1 * time.Second) // Check every second
			defer ticker.Stop()

			for range ticker.C {
				mu.Lock()
				shouldFlush := len(pendingUpdates) > 0 && time.Since(lastFlush) >= flushInterval
				mu.Unlock()

				if shouldFlush {
					flushUpdates()
				}
			}
		}()

		// Collect updates
		for update := range subscriber {
			mu.Lock()
			pendingUpdates = append(pendingUpdates, buildmanager.LogUpdate{
				Message: update.Message,
				Append:  update.Append,
			})
			shouldFlush := len(pendingUpdates) >= maxBatchSize
			mu.Unlock()

			if shouldFlush {
				flushUpdates()
			}
		}

		// Final flush
		flushUpdates()
	}
*/
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
			self.buildManager.Build(self.buildPath)
			log.Println("all building done, display final message")

			self.buildStatusLabel.SetText("Building Finished!")

		}()
	})

	hbox := container.NewVBox(buildButton, self.logScroll, self.buildStatusLabel)
	//hbox := container.NewVBox(buildButton, self.logWidget.GetWidget(), self.buildStatusLabel)
	return hbox
}
