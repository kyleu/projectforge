package filter

type Options struct {
	Sort   []string `json:"sort,omitempty"`
	Filter []string `json:"filter,omitempty"`
	Group  []string `json:"group,omitempty"`
	Search string   `json:"search,omitempty"`
	Params *Params  `json:"params,omitempty"`
}

type OptionsMap map[string]*Options

func (m OptionsMap) Get(key string) *Options {
	ret, ok := m[key]
	if ok {
		return ret
	}
	return &Options{Params: &Params{Key: key}}
}
