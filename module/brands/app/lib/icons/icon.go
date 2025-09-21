package icons

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	svgPrefix   = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 24 24">`
	symbolStart = `<symbol id="svg-%s" viewBox="0 0 24 24">`
	symbolEnd   = `</symbol>`
	symbolLink  = `<use xlink:href="#svg-%s" />`
	svgSuffix   = `</svg>`
)

type Icon struct {
	Key        string   `json:"key"`
	Title      string   `json:"title"`
	Color      string   `json:"color,omitzero"`
	Source     string   `json:"source,omitzero"`
	Aliases    []string `json:"aliases,omitempty"`
	Guidelines string   `json:"guidelines,omitzero"`
	License    string   `json:"license,omitzero"`
	Path       string   `json:"path,omitzero"`
}

func (i *Icon) SetContent(s string) error {
	if !strings.Contains(s, "0 0 24 24") {
		return errors.Errorf("content for icon [%s] is not in the correct format", i.Key)
	}
	if strings.Count(s, "<path") != 1 {
		return errors.Errorf("content for icon [%s] does not contain a single path", i.Key)
	}
	pathIdx := strings.Index(s, "<path")
	i.Path = strings.TrimSuffix(s[pathIdx:], svgSuffix)
	if idx := strings.Index(i.Path, "class="); idx == -1 {
		i.Path = strings.Replace(i.Path, "<path", `<path class="svg-fill"`, 1)
	}
	return nil
}

func (i *Icon) HTML(prefix string) string {
	return svgPrefix + fmt.Sprintf(symbolStart, prefix+i.Key) + i.Path + symbolEnd + fmt.Sprintf(symbolLink, prefix+i.Key) + svgSuffix
}

type Icons []*Icon
