package doctor

var AllDependencies = Dependencies{rsvg}

var rsvg = &Dependency{
	Key:     "rsvg",
	Title:   "rsvg",
	Summary: "Renders SVGs to PNG for the icon pipeline",
	Cmd:     "rsvg-convert",
	Args:    []string{"-v"},
}
