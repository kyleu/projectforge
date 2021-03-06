package cutil

type Location struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Icon  string `json:"icon,omitempty"`
}

type Locations []*Location
