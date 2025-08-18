package logger

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

type LogView struct {
	widget.BaseWidget
	lines []string
	max   int
	lineH float32
}

func NewLogView(maxLines int) *LogView {
	l := &LogView{
		lines: []string{},
		max:   maxLines,
		lineH: float32(theme.TextSize()) + 4, // spacing per line
	}
	l.ExtendBaseWidget(l)
	return l
}

func (l *LogView) Clear() {
	l.lines = l.lines[:0]
	l.Refresh()
}

// Append a new log line
func (l *LogView) AppendLine(line string) {
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
	bg := canvas.NewRectangle(theme.Color(theme.ColorNameBackground))
	return &logRenderer{log: l, bg: bg}
}

type logRenderer struct {
	log *LogView
	bg  *canvas.Rectangle
}

func (r *logRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
}

func (r *logRenderer) MinSize() fyne.Size {
	// Width grows with longest line, height grows with number of lines
	width := float32(200)
	for _, line := range r.log.lines {
		w := fyne.MeasureText(line, theme.TextSize(), fyne.TextStyle{}).Width
		if w > width {
			width = w
		}
	}
	height := float32(len(r.log.lines)) * r.log.lineH
	return fyne.NewSize(width+8, height)
}

func (r *logRenderer) Refresh() {
	// Rebuild all objects: background + text lines
	objects := []fyne.CanvasObject{r.bg}
	y := float32(0)
	for _, line := range r.log.lines {
		txt := canvas.NewText(line, theme.Color(theme.ColorNameForeground))
		txt.TextSize = theme.TextSize()
		txt.TextStyle = fyne.TextStyle{Monospace: true}
		txt.Move(fyne.NewPos(4, y))
		y += r.log.lineH
		objects = append(objects, txt)
	}
	r.bg.FillColor = theme.Color(theme.ColorNameBackground)
	r.bg.Refresh()
	canvas.Refresh(r.log)
	// Replace renderer’s object list
	rObjs := make([]fyne.CanvasObject, len(objects))
	copy(rObjs, objects)
	rObjects[r.log] = rObjs
}

var rObjects = make(map[*LogView][]fyne.CanvasObject)

func (r *logRenderer) BackgroundColor() color.Color { return theme.Color(theme.ColorNameBackground) }
func (r *logRenderer) Objects() []fyne.CanvasObject { return rObjects[r.log] }
func (r *logRenderer) Destroy()                     {}
