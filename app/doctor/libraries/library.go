package libraries

var AllLibraries = Libraries{
	{Key: "string", Name: "Rich String", Icon: "cog"},
	{Key: "filesystem", Name: "Filesystem", Icon: "file", Files: map[string]string{
		"module/filesystem/app/lib/filesystem": "./",
	}},
}

type Library struct {
	Key   string            `json:"key"`
	Name  string            `json:"name"`
	Icon  string            `json:"icon"`
	Files map[string]string `json:"files"`
}

func (l *Library) String() string {
	return l.Name
}

type Libraries []*Library

func (l Libraries) Get(key string) *Library {
	for _, x := range l {
		if x.Key == key {
			return x
		}
	}
	return nil
}

func (l Libraries) Index(key string) int {
	for i, x := range l {
		if x.Key == key {
			return i
		}
	}
	return -1
}
