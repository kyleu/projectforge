package har

import (
	"fmt"

	"{{{ .Package }}}/app/util"
)

type Page struct {
	StartedDateTime string      `json:"startedDateTime"`
	ID              string      `json:"id"`
	Title           string      `json:"title"`
	PageTiming      *PageTiming `json:"pageTiming"`
	Comment         string      `json:"comment,omitempty"`
}

type Pages []*Page

type PageTiming struct {
	OnContentLoad int    `json:"onContentLoad"`
	OnLoad        int    `json:"onLoad"`
	Comment       string `json:"comment"`
}

type PageTimings struct {
	IsMicros bool    `json:"isMicros,omitempty"`
	Total    float64 `json:"total,omitempty"`
	Blocked  float64 `json:"blocked,omitempty"`
	DNS      float64 `json:"dns,omitempty"`
	Connect  float64 `json:"connect,omitempty"`
	Send     float64 `json:"send,omitempty"`
	Wait     float64 `json:"wait,omitempty"`
	Receive  float64 `json:"receive,omitempty"`
	SSL      float64 `json:"ssl,omitempty"`
	Comment  string  `json:"comment,omitempty"`
}

func (p *PageTimings) Elapsed() int {
	if p.Total > 0 {
		return int(p.Total)
	}
	return int(p.Blocked + p.DNS + p.Connect + p.SSL + p.Wait + p.Receive + p.SSL)
}

func (p *PageTimings) Map() map[string]string {
	mult := 1000
	if p.IsMicros {
		mult = 1
	}
	ret := make(map[string]string, 10)
	if p.Blocked > 0 {
		ret["Blocked"] = util.MicrosToMillis(int(p.Blocked) * mult)
	}
	if p.DNS > 0 {
		ret["DNS"] = util.MicrosToMillis(int(p.DNS) * mult)
	}
	if p.Connect > 0 {
		ret["Connect"] = util.MicrosToMillis(int(p.Connect) * mult)
	}
	if p.Send > 0 {
		ret["Send"] = util.MicrosToMillis(int(p.Send) * mult)
	}
	if p.Wait > 0 {
		ret["Wait"] = util.MicrosToMillis(int(p.Wait) * mult)
	}
	// if p.Receive > 0 {
	// 	ret["Receive"] = util.MicrosToMillis(int(p.Receive) * mult)
	// }
	if p.SSL > 0 {
		ret["SSL"] = util.MicrosToMillis(int(p.SSL) * mult)
	}
	return ret
}

func (p *PageTimings) Strings() []string {
	m := p.Map()
	ret := make([]string, 0, len(m))
	for k, v := range m {
		ret = append(ret, fmt.Sprintf("%s: %s", k, v))
	}
	return ret
}
