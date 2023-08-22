package har

import "github.com/samber/lo"

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
	for _, x := range p {
		if x.Name == n {
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
	Encoding    string `json:"encoding,omitempty"`
	Comment     string `json:"comment,omitempty"`
	File        string `json:"_file,omitempty"`
}
