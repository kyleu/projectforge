package action

import (
	"bytes"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/exec"
	"projectforge.dev/projectforge/app/util"
)

func runCodeStats(pth string, recursive bool) (*CodeStats, error) {
	cmd := `tokei . -o json --exclude "*.html.go"  --exclude "*.sql.go"`
	ex := exec.NewExec("tokei", 0, cmd, pth, false)
	if err := ex.Run(); err != nil {
		if strings.Contains(err.Error(), "no such file") {
			return nil, nil
		}
		return nil, err
	}
	b := ex.Buffer.Bytes()
	if idx := bytes.Index(b, []byte(" ::: ")); idx > -1 {
		b = b[:idx]
	}
	cs, err := util.FromJSONMap(bytes.TrimSpace(b))
	if err != nil {
		if !recursive {
			return runCodeStats(pth, true)
		}
		return nil, err
	}
	ci := newCodeStats()
	for _, k := range cs.Keys() {
		if k == "Total" {
			continue
		}
		x := cs.GetMapOpt(k)
		rpts, err := x.GetMapArray("reports", true)
		if err != nil {
			return nil, err
		}
		ct := newCodeTypeFromMap(k, x)

		if len(rpts) != 0 {
			for _, rpt := range rpts {
				rk := rpt.GetStringOpt("name")
				s := rpt.GetMapOpt("stats")
				var total int
				for _, sk := range s.Keys() {
					if i := s.GetIntOpt(sk); i > 0 {
						total += i
					}
				}
				ctf := newCodeTypeFromMap(rk, s)
				ct.Files = append(ct.Files, ctf)
			}
		}
		ci.Add(ct)
	}
	return ci, nil
}

type CodeType struct {
	Name     string    `json:"name"`
	Code     int       `json:"code,omitzero"`
	Comments int       `json:"comments,omitzero"`
	Blanks   int       `json:"blanks,omitzero"`
	Files    CodeTypes `json:"files,omitempty"`
}

func newCodeType(name string) *CodeType {
	return &CodeType{Name: name}
}

func newCodeTypeFromMap(k string, m util.ValueMap) *CodeType {
	k = strings.TrimPrefix(k, "./")
	return &CodeType{Name: k, Code: m.GetIntOpt("code"), Comments: m.GetIntOpt("comments"), Blanks: m.GetIntOpt("blanks")}
}

func (c *CodeType) Total() int {
	return c.Code + c.Comments + c.Blanks
}

func (c *CodeType) Add(x *CodeType) {
	c.Code += x.Code
	c.Comments += x.Comments
	c.Blanks += x.Blanks
	c.Files = append(c.Files, x.Files...)
}

type CodeTypes []*CodeType

type CodeStats struct {
	Types CodeTypes `json:"types"`
	Total *CodeType `json:"total"`
}

func newCodeStats() *CodeStats {
	return &CodeStats{Total: newCodeType("Total")}
}

func (c *CodeStats) Add(x *CodeType) {
	c.Types = append(c.Types, x)
	c.Total.Add(x)
}

func (c *CodeStats) ToMaps() []util.ValueMap {
	return lo.Map(c.Types, func(x *CodeType, _ int) util.ValueMap {
		return util.ValueMap{
			"name":     x.Name,
			"code":     x.Code,
			"comments": x.Comments,
			"blanks":   x.Blanks,
			"files":    len(x.Files),
		}
	})
}
