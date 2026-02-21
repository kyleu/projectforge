package layout

const (
	compactBreakpointWidth  = 100
	compactBreakpointHeight = 24
	statusBarHeight         = 2
	editorBarHeight         = 1
	compactHeaderHeight     = 3
	defaultHeaderHeight     = 1
	nonCompactSidebarWidth  = 34
	// Bordered horizontal joins can render short in some terminals; this keeps right edges visible.
	nonCompactWidthCompensation = 2
)

type Rect struct {
	X int
	Y int
	W int
	H int
}

type Rects struct {
	Compact bool
	Header  Rect
	Main    Rect
	Sidebar Rect
	Editor  Rect
	Status  Rect
	Overlay Rect
}

func Solve(width int, height int) Rects {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	compact := width < compactBreakpointWidth || height < compactBreakpointHeight
	statusHeight := statusBarHeight
	editorHeight := editorBarHeight
	headerHeight := defaultHeaderHeight
	if compact {
		headerHeight = compactHeaderHeight
	}
	bodyHeight := max(1, height-statusHeight-editorHeight-headerHeight)

	ret := Rects{Compact: compact}
	ret.Header = Rect{X: 0, Y: 0, W: width, H: headerHeight}
	ret.Status = Rect{X: 0, Y: height - statusHeight, W: width, H: statusHeight}
	ret.Editor = Rect{X: 0, Y: height - statusHeight - editorHeight, W: width, H: editorHeight}

	if compact {
		ret.Main = Rect{X: 0, Y: headerHeight, W: width, H: bodyHeight}
		ret.Overlay = Rect{X: 0, Y: headerHeight, W: width, H: bodyHeight}
		return ret
	}

	sidebarWidth := nonCompactSidebarWidth
	contentWidth := max(1, width+nonCompactWidthCompensation)
	if sidebarWidth > contentWidth/2 {
		sidebarWidth = contentWidth / 2
	}
	mainWidth := max(1, contentWidth-sidebarWidth)
	ret.Main = Rect{X: 0, Y: headerHeight, W: mainWidth, H: bodyHeight}
	ret.Sidebar = Rect{X: mainWidth, Y: headerHeight, W: sidebarWidth, H: bodyHeight}
	ret.Overlay = Rect{X: 0, Y: headerHeight, W: contentWidth, H: bodyHeight}
	return ret
}
