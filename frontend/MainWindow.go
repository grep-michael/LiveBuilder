package frontend

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

// MainWindow wraps fyne.Window with custom functionality
type MainWindow struct {
	app    fyne.App
	window fyne.Window
	title  string
}

func NewMainWindow(title string) *MainWindow {
	myApp := app.New()
	myWindow := myApp.NewWindow(title)

	mw := &MainWindow{
		app:    myApp,
		window: myWindow,
		title:  title,
	}
	mw.BuildMainContent()
	return mw
}

func (self *MainWindow) BuildMainContent() {
	file_selector := buildFileSelectionView()
	tabs := container.NewAppTabs(
		container.NewTabItem("lb config Editor", buildLBConfigView()),
		container.NewTabItem("File Selection", file_selector),
	)
	self.SetContent(tabs)
}

func (mw *MainWindow) SetContent(content fyne.CanvasObject) {
	mw.window.SetContent(content)
	mw.autoResize(content)
}

func (mw *MainWindow) SetContentWithPadding(content fyne.CanvasObject, padding fyne.Size) {
	mw.window.SetContent(content)
	mw.resizeWithPadding(content, padding)
}

func (mw *MainWindow) autoResize(content fyne.CanvasObject) {
	defaultPadding := fyne.NewSize(30, 60)
	mw.resizeWithPadding(content, defaultPadding)
}

func (mw *MainWindow) resizeWithPadding(content fyne.CanvasObject, padding fyne.Size) {
	contentSize := content.Size()
	windowSize := fyne.NewSize(
		contentSize.Width+padding.Width,
		contentSize.Height+padding.Height,
	)
	mw.window.Resize(windowSize)
}

func (mw *MainWindow) AddWidget(widget fyne.CanvasObject) {
	currentContent := mw.window.Content()
	if currentContent == nil {
		mw.SetContent(widget)
		return
	}

	// If content exists, create a vertical container
	vbox := container.NewVBox(currentContent, widget)
	mw.SetContent(vbox)
}

func (mw *MainWindow) SetFixedSize(width, height float32) {
	mw.window.Resize(fyne.NewSize(width, height))
	mw.window.SetFixedSize(true)
}

func (mw *MainWindow) SetResizable(resizable bool) {
	mw.window.SetFixedSize(!resizable)
}

func (mw *MainWindow) ShowAndRun() {
	mw.window.ShowAndRun()
}

func (mw *MainWindow) Show() {
	mw.window.Show()
}

func (mw *MainWindow) Close() {
	mw.window.Close()
}

func (mw *MainWindow) GetWindow() fyne.Window {
	return mw.window
}

func (mw *MainWindow) GetApp() fyne.App {
	return mw.app
}

func (mw *MainWindow) SetTitle(title string) {
	mw.title = title
	mw.window.SetTitle(title)
}

func (mw *MainWindow) GetTitle() string {
	return mw.title
}
