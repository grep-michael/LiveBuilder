package buildwindow

import (
	//"fmt"
	buildmanager "LiveBuilder/BuildManager"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// SettingsWidget is a custom widget that shows a context menu on right-click
type SettingsWidget struct {
	widget.BaseWidget
	content      fyne.CanvasObject
	window       fyne.Window
	buildManager *buildmanager.BuildManager
	// Checkbox data
	checkboxes []*widget.Check
	checkOrder []string
	checkData  map[string]bool
}

func NewRightClickWidget(window fyne.Window, buildManager *buildmanager.BuildManager) *SettingsWidget {

	checkData := map[string]bool{
		"lb bootstrap":          true,
		"lb chroot":             true,
		"lb binary_chroot":      false,
		"lb binary_rootfs":      false,
		"lb binary_linux-image": false,
		"lb binary":             true,
		"lb source":             true,
	}
	checkOrder := []string{
		"lb bootstrap",
		"lb chroot",
		"lb binary_chroot",
		"lb binary_rootfs",
		"lb binary_linux-image",
		"lb binary",
		"lb source",
	}

	w := &SettingsWidget{
		content:      widget.NewIcon(theme.SettingsIcon()),
		window:       window,
		checkData:    checkData,
		checkOrder:   checkOrder,
		buildManager: buildManager,
	}
	w.ExtendBaseWidget(w)
	w.Resize(fyne.NewSize(100, 100))
	return w
}

func (w *SettingsWidget) CreateRenderer() fyne.WidgetRenderer {
	return &settingsRender{
		widget:  w,
		content: w.content,
	}
}

func (w *SettingsWidget) MouseDown(event *desktop.MouseEvent) {
	if event.Button == desktop.MouseButtonSecondary {
		w.showContextMenu(event.AbsolutePosition)
	}
}

func (w *SettingsWidget) TappedSecondary(event *fyne.PointEvent) {
	canvasPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(w)

	// Position the menu near the widget
	menuPos := fyne.NewPos(
		canvasPos.X+event.Position.X,
		canvasPos.Y+event.Position.Y+25, // Offset slightly below the icon
	)

	w.showContextMenu(menuPos)
}

func (w *SettingsWidget) showContextMenu(pos fyne.Position) {
	// Create checkboxes
	checklist := w.buildCheckList()
	buttons := w.buildActionButtons()

	cont := container.NewHBox(checklist, buttons)

	popup := widget.NewPopUp(cont, w.window.Canvas())
	popup.ShowAtPosition(pos)
}

func (w *SettingsWidget) buildCheckList() *fyne.Container {
	w.checkboxes = make([]*widget.Check, len(w.checkData))
	var menuItems []fyne.CanvasObject

	for _, key := range w.checkOrder {

		check := widget.NewCheck(key, func(checked bool) {
			w.checkData[key] = checked
		})
		check.SetChecked(w.checkData[key])
		w.checkboxes = append(w.checkboxes, check)
		menuItems = append(menuItems, check)
	}

	return container.NewVBox(menuItems...)
}

func (w *SettingsWidget) buildActionButtons() *fyne.Container {

	cleanBtn := widget.NewButton("Purge Build", func() {
		w.buildManager.NukeBuild()
	})

	return container.NewVBox(cleanBtn)
}

func (w *SettingsWidget) GetCheckboxStates() map[string]bool {
	return w.checkData
}

type settingsRender struct {
	widget  *SettingsWidget
	content fyne.CanvasObject
}

func (r *settingsRender) Layout(size fyne.Size) {
	if r.content != nil {
		r.content.Resize(size)
	}
}

func (r *settingsRender) MinSize() fyne.Size {
	if r.content != nil {
		return r.content.MinSize()
	}
	return fyne.NewSize(100, 50)
}

func (r *settingsRender) Refresh() {
	if r.content != nil {
		r.content.Refresh()
	}
}

func (r *settingsRender) Objects() []fyne.CanvasObject {
	if r.content != nil {
		return []fyne.CanvasObject{r.content}
	}
	return []fyne.CanvasObject{}
}

func (r *settingsRender) Destroy() {}
