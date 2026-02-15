package layout

const (
	compactBreakpointWidth  = 100
	compactBreakpointHeight = 24
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
	statusHeight := 2
	editorHeight := 3
	bodyHeight := max(1, height-statusHeight-editorHeight)

	ret := Rects{Compact: compact}
	ret.Status = Rect{X: 0, Y: height - statusHeight, W: width, H: statusHeight}
	ret.Editor = Rect{X: 0, Y: height - statusHeight - editorHeight, W: width, H: editorHeight}

	if compact {
		headerHeight := 3
		mainHeight := max(1, bodyHeight-headerHeight)
		ret.Header = Rect{X: 0, Y: 0, W: width, H: headerHeight}
		ret.Main = Rect{X: 0, Y: headerHeight, W: width, H: mainHeight}
		ret.Overlay = Rect{X: 0, Y: headerHeight, W: width, H: mainHeight}
		return ret
	}

	sidebarWidth := 32
	if sidebarWidth > width/2 {
		sidebarWidth = width / 2
	}
	mainWidth := max(1, width-sidebarWidth)
	ret.Main = Rect{X: 0, Y: 0, W: mainWidth, H: bodyHeight}
	ret.Sidebar = Rect{X: mainWidth, Y: 0, W: sidebarWidth, H: bodyHeight}
	ret.Overlay = Rect{X: 0, Y: 0, W: width, H: bodyHeight}
	return ret
}
