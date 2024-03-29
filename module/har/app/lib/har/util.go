package har

import (
	"strings"

	"github.com/samber/lo"
)

const Ext = ".har"

type Wrapper struct {
	Log *Log `json:"log"`
}

type NVP struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Comment string `json:"comment,omitempty"`
}

type NVPs []*NVP

func (p NVPs) GetValue(n string) string {
	n = strings.ToLower(n)
	for _, x := range p {
		if strings.EqualFold(x.Name, n) {
			return x.Value
		}
	}
	return ""
}

func (p NVPs) WithReplacements(repl func(s string) string) NVPs {
	return lo.Map(p, func(n *NVP, _ int) *NVP {
		return &NVP{Name: repl(n.Name), Value: repl(n.Value), Comment: repl(n.Comment)}
	})
}

type Content struct {
	Size        int    `json:"size"`
	Compression int    `json:"compression,omitempty"`
	MimeType    string `json:"mimeType"`
	Text        string `json:"text,omitempty"`
	JSON        any    `json:"json,omitempty"`
	Encoding    string `json:"encoding,omitempty"`
	Comment     string `json:"comment,omitempty"`
	File        string `json:"_file,omitempty"`
}
