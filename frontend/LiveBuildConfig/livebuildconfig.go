package livebuildconfig

import (
	appstate "LiveBuilder/AppState"
	cmds "LiveBuilder/Commands"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewlbTextEdit() *fyne.Container {
	cfg := cmds.NewlbConfig()
	cfgtext := cfg.BuildTemplate()
	textEntry := widget.NewMultiLineEntry()
	state := appstate.GetGlobalState()
	state.LBConfigCMD = cfgtext
	textEntry.OnChanged = func(text string) {
		state.LBConfigCMD = text
	}
	textEntry.Text = cfgtext
	textEntry.Wrapping = fyne.TextWrapWord
	return container.NewBorder(nil, nil, nil, nil, textEntry)

}
