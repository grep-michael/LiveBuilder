package logger

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type LogView struct {
	widget.BaseWidget
	max     int
	lineH   float32
	objects []*canvas.Text
}

func NewLogView(maxLines int) *LogView {
	l := &LogView{
		objects: []*canvas.Text{},
		max:     maxLines,
		lineH:   float32(theme.TextSize()) + 4, // spacing per line
	}
	l.ExtendBaseWidget(l)
	return l
}

func (l *LogView) Clear() {
	l.objects = l.objects[:0]
	l.Refresh()
}

// Append a new log line
func (l *LogView) AppendLine(line string) {
	line = strings.ReplaceAll(line, "\n", "")
	txt := canvas.NewText(line, theme.Color(theme.ColorNameForeground))
	txt.TextSize = theme.TextSize()
	txt.TextStyle = fyne.TextStyle{Monospace: true}
	if len(l.objects) >= l.max {
		l.objects = append(l.objects[1:], txt)
	} else {
		l.objects = append(l.objects, txt)
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

func (r *logRenderer) MinSize() fyne.Size {
	// Width grows with longest line, height grows with number of lines
	width := float32(200)
	for _, line := range r.log.objects {
		w := fyne.MeasureText(line.Text, theme.TextSize(), fyne.TextStyle{}).Width
		if w > width {
			width = w
		}
	}
	height := float32(len(r.log.objects)) * r.log.lineH
	return fyne.NewSize(width+8, height)
}

func (r *logRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)

	y := float32(0)
	for _, obj := range r.log.objects {
		obj.Move(fyne.NewPos(4, y))
		y += r.log.lineH
	}
}

func (r *logRenderer) Refresh() {
	r.bg.FillColor = theme.Color(theme.ColorNameBackground)
	r.bg.Refresh()

	for _, t := range r.log.objects {
		t.Color = theme.Color(theme.ColorNameForeground)
		t.Refresh()
	}
	canvas.Refresh(r.log)
}

func (r *logRenderer) BackgroundColor() color.Color { return theme.Color(theme.ColorNameBackground) }
func (r *logRenderer) Objects() []fyne.CanvasObject {
	objs := make([]fyne.CanvasObject, 1+len(r.log.objects))
	objs[0] = r.bg
	for i, t := range r.log.objects {
		objs[i+1] = t
	}
	return objs
}
func (r *logRenderer) Destroy() {}
