package har

import (
	"time"

	"{{{ .Package }}}/app/util"
)

type Cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Path     string `json:"path,omitzero"`
	Domain   string `json:"domain,omitzero"`
	Expires  string `json:"expires,omitzero"`
	HTTPOnly bool   `json:"httpOnly,omitzero"`
	Secure   bool   `json:"secure,omitzero"`
	Comment  bool   `json:"comment,omitzero"`
}

func (c *Cookie) Tags() []string {
	var ret []string
	if c.HTTPOnly {
		ret = append(ret, "http-only")
	}
	if c.Secure {
		ret = append(ret, "secure")
	}
	if c.Comment {
		ret = append(ret, "comment")
	}
	return ret
}

func (c *Cookie) Exp() *time.Time {
	ret, err := time.Parse("2006-01-02T15:04:05.000Z", c.Expires)
	if err != nil {
		return nil
	}
	return &ret
}

func (c *Cookie) ExpRelative() string {
	e := c.Exp()
	return util.TimeRelative(e)
}

type Cookies []*Cookie
