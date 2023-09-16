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
	IsMicros bool   `json:"isMicros,omitempty"`
	Total    int    `json:"total,omitempty"`
	Blocked  int    `json:"blocked,omitempty"`
	DNS      int    `json:"dns,omitempty"`
	Connect  int    `json:"connect,omitempty"`
	Send     int    `json:"send,omitempty"`
	Wait     int    `json:"wait,omitempty"`
	Receive  int    `json:"receive,omitempty"`
	SSL      int    `json:"ssl,omitempty"`
	Comment  string `json:"comment,omitempty"`
}

func (p *PageTimings) Map() map[string]string {
	mult := 1000
	if p.IsMicros {
		mult = 1
	}
	ret := make(map[string]string, 10)
	if p.Blocked > 0 {
		ret["Blocked"] = util.MicrosToMillis(p.Blocked * mult)
	}
	if p.DNS > 0 {
		ret["DNS"] = util.MicrosToMillis(p.DNS * mult)
	}
	if p.Connect > 0 {
		ret["Connect"] = util.MicrosToMillis(p.Connect * mult)
	}
	if p.Send > 0 {
		ret["Send"] = util.MicrosToMillis(p.Send * mult)
	}
	if p.Wait > 0 {
		ret["Wait"] = util.MicrosToMillis(p.Wait * mult)
	}
	// if p.Receive > 0 {
	// 	ret["Receive"] = util.MicrosToMillis(p.Receive * mult)
	// }
	if p.SSL > 0 {
		ret["SSL"] = util.MicrosToMillis(p.SSL * mult)
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
