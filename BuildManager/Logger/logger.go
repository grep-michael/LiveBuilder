package logger

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"

	"fyne.io/fyne/v2/widget"
)

// LogView is a custom widget for displaying lots of logs
type LogView struct {
	widget.BaseWidget
	lines   []string
	max     int
	lineH   float32
	bgColor color.Color
	fgColor color.Color
}

func NewLogView(maxLines int) *LogView {
	l := &LogView{
		lines:   []string{},
		max:     maxLines,
		lineH:   16, // px per line (tweak based on font size)
		bgColor: theme.Color(theme.ColorNameInputBackground),
		fgColor: theme.Color(theme.ColorNameInputBorder),
	}
	l.ExtendBaseWidget(l)
	return l
}

func (l *LogView) Clear() {
	l.lines = []string{}
	l.Refresh()
}

// Append a new log line
func (l *LogView) AddLine(line string) {
	if len(l.lines) >= l.max {
		copy(l.lines, l.lines[1:])
		l.lines[len(l.lines)-1] = line
	} else {
		l.lines = append(l.lines, line)
	}
	l.Refresh()
}

// Implement fyne.WidgetRenderer
func (l *LogView) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(l.bgColor)
	objects := []fyne.CanvasObject{bg}
	return &logRenderer{log: l, bg: bg, objects: objects}
}

type logRenderer struct {
	log     *LogView
	bg      *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (r *logRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
}

func (r *logRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, 200)
}

func (r *logRenderer) Refresh() {
	r.objects = r.objects[:1] // keep bg only

	// draw visible lines
	y := float32(0)
	for _, line := range r.log.lines {
		txt := canvas.NewText(line, r.log.fgColor)
		txt.TextSize = 14
		txt.Move(fyne.NewPos(4, y))
		y += r.log.lineH
		r.objects = append(r.objects, txt)
	}
	canvas.Refresh(r.log)
}

func (r *logRenderer) BackgroundColor() color.Color { return r.log.bgColor }
func (r *logRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *logRenderer) Destroy()                     {}
