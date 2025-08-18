package logger

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
	"sync"
)

// LogWidget provides an efficient log display with minimal padding
type LogWidget struct {
	lines      []string
	maxLines   int
	container  *container.Scroll
	content    *fyne.Container
	mutex      sync.RWMutex
	autoScroll bool
	minSize    fyne.Size
	labels     []*widget.Label
}

// NewLogWidget creates a new log widget with variable height items
func NewLogWidget(maxLines int) *LogWidget {
	if maxLines <= 0 {
		maxLines = 1000
	}

	lw := &LogWidget{
		lines:      make([]string, 0),
		maxLines:   maxLines,
		autoScroll: true,
		minSize:    fyne.NewSize(600, 400),
		labels:     make([]*widget.Label, 0),
	}

	lw.setupContainer()
	return lw
}

// setupContainer initializes the container for compact log items
func (lw *LogWidget) setupContainer() {
	lw.content = container.NewVBox()
	lw.container = container.NewScroll(lw.content)
	lw.container.SetMinSize(lw.minSize)
}

func (lw *LogWidget) GetWidget() fyne.CanvasObject {
	return lw.container
}

func (lw *LogWidget) createCompactLabel(text string) *fyne.Container {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	label.TextStyle = fyne.TextStyle{}
	label.Alignment = fyne.TextAlignLeading

	// Wrap in a border container with minimal padding
	container := container.NewBorder(nil, nil, nil, nil, label)

	return container
}

// AppendLine adds a new line to the log with minimal padding
func (lw *LogWidget) AppendLine(line string) {
	lw.mutex.Lock()

	// Handle multi-line strings
	splitLines := strings.Split(line, "\n")
	for _, singleLine := range splitLines {
		// Skip empty lines
		if strings.TrimSpace(singleLine) != "" {
			lw.lines = append(lw.lines, singleLine)

			// Create label container with minimal padding
			labelContainer := lw.createCompactLabel(singleLine)

			// Add to container
			lw.content.Add(labelContainer)
		}
	}

	// Trim old lines if we exceed max capacity
	if len(lw.lines) > lw.maxLines {
		excess := len(lw.lines) - lw.maxLines

		// Remove excess lines
		lw.lines = lw.lines[excess:]

		// Remove excess containers from the beginning
		objects := lw.content.Objects
		if len(objects) >= excess {
			for i := 0; i < excess; i++ {
				lw.content.Remove(objects[i])
			}
		}
	}

	lw.mutex.Unlock()

	lw.content.Refresh()
	// Auto-scroll to bottom if enabled
	if lw.autoScroll {
		lw.ScrollToBottom()
	}

}

// AppendLines adds multiple lines efficiently with minimal padding
func (lw *LogWidget) AppendLines(lines []string) {
	lw.mutex.Lock()

	for _, line := range lines {
		splitLines := strings.Split(line, "\n")
		for _, singleLine := range splitLines {
			// Skip empty lines
			if strings.TrimSpace(singleLine) != "" {
				lw.lines = append(lw.lines, singleLine)

				// Create label container with minimal padding
				labelContainer := lw.createCompactLabel(singleLine)

				// Add to container
				lw.content.Add(labelContainer)
			}
		}
	}

	// Trim old lines if we exceed max capacity
	if len(lw.lines) > lw.maxLines {
		excess := len(lw.lines) - lw.maxLines

		// Remove excess lines
		lw.lines = lw.lines[excess:]

		// Remove excess containers from the beginning
		objects := lw.content.Objects
		if len(objects) >= excess {
			for i := 0; i < excess; i++ {
				lw.content.Remove(objects[i])
			}
		}
	}

	lw.mutex.Unlock()

	// Refresh for all lines
	lw.content.Refresh()

	if lw.autoScroll {
		lw.ScrollToBottom()
	}
}

// Clear removes all log lines and containers
func (lw *LogWidget) Clear() {
	lw.mutex.Lock()
	lw.lines = make([]string, 0)

	// Remove all containers
	lw.content.RemoveAll()
	lw.labels = make([]*widget.Label, 0)

	lw.mutex.Unlock()
	lw.content.Refresh()

}

// SetAutoScroll enables or disables automatic scrolling to bottom
func (lw *LogWidget) SetAutoScroll(enabled bool) {
	lw.mutex.Lock()
	lw.autoScroll = enabled
	lw.mutex.Unlock()
}

// ScrollToBottom scrolls to the bottom of the log
func (lw *LogWidget) ScrollToBottom() {
	lw.mutex.RLock()
	lineCount := len(lw.lines)
	lw.mutex.RUnlock()

	if lineCount > 0 {
		lw.container.ScrollToBottom()
	}
}

// ScrollToTop scrolls to the top of the log
func (lw *LogWidget) ScrollToTop() {
	lw.container.ScrollToTop()
}

// GetLineCount returns the current number of lines
func (lw *LogWidget) GetLineCount() int {
	lw.mutex.RLock()
	defer lw.mutex.RUnlock()
	return len(lw.lines)
}

// GetMaxLines returns the maximum number of lines
func (lw *LogWidget) GetMaxLines() int {
	lw.mutex.RLock()
	defer lw.mutex.RUnlock()
	return lw.maxLines
}

// SetMaxLines updates the maximum number of lines
func (lw *LogWidget) SetMaxLines(maxLines int) {
	if maxLines <= 0 {
		maxLines = 1000
	}

	lw.mutex.Lock()
	lw.maxLines = maxLines

	// Trim existing lines if needed
	if len(lw.lines) > lw.maxLines {
		excess := len(lw.lines) - lw.maxLines
		lw.lines = lw.lines[excess:]
	}
	lw.mutex.Unlock()
}

// SetMinSize sets the minimum size of the log widget
func (lw *LogWidget) SetMinSize(size fyne.Size) {
	lw.minSize = size
	if lw.container != nil {
		lw.container.SetMinSize(size)
	}
}

// SetSize sets the current size of the log widget
func (lw *LogWidget) SetSize(size fyne.Size) {
	if lw.container != nil {
		lw.container.Resize(size)
	}
}

// GetMinSize returns the current minimum size
func (lw *LogWidget) GetMinSize() fyne.Size {
	return lw.minSize
}

// GetAllLines returns a copy of all current lines (for export/save functionality)
func (lw *LogWidget) GetAllLines() []string {
	lw.mutex.RLock()
	defer lw.mutex.RUnlock()

	// Return a copy to prevent external modification
	result := make([]string, len(lw.lines))
	copy(result, lw.lines)
	return result
}

// SearchLines returns line numbers containing the search term (useful for your ITAD logs)
func (lw *LogWidget) SearchLines(searchTerm string) []int {
	lw.mutex.RLock()
	defer lw.mutex.RUnlock()

	var matches []int
	searchLower := strings.ToLower(searchTerm)

	for i, line := range lw.lines {
		if strings.Contains(strings.ToLower(line), searchLower) {
			matches = append(matches, i)
		}
	}

	return matches
}

// ScrollToLine scrolls to a specific line number (best effort)
func (lw *LogWidget) ScrollToLine(lineNum int) {
	lw.mutex.RLock()
	lineCount := len(lw.lines)
	lw.mutex.RUnlock()

	if lineNum >= 0 && lineNum < lineCount {
		// Calculate approximate scroll position
		ratio := float32(lineNum) / float32(lineCount)

		// Get content height and scroll accordingly
		contentHeight := lw.content.Size().Height
		containerHeight := lw.container.Size().Height

		if contentHeight > containerHeight {
			scrollOffset := ratio * (contentHeight - containerHeight)
			lw.container.Offset = fyne.NewPos(0, scrollOffset)
			lw.container.Refresh()
		}
	}
}
