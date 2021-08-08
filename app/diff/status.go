package diff

import (
	"github.com/kyleu/projectforge/app/util"
)

type Status struct {
	Key   string
	Title string
}

var (
	StatusDifferent = &Status{Key: "different", Title: "The files are different"}
	StatusIdentical = &Status{Key: "identical", Title: "The files are identical"}
	StatusMissing   = &Status{Key: "missing", Title: "File is not present in the source"}
	StatusNew       = &Status{Key: "new", Title: "File is not present in the target"}
	StatusSkipped   = &Status{Key: "skipped", Title: "File was ignored"}
)

var AllStatuses = []*Status{StatusDifferent, StatusIdentical, StatusMissing, StatusNew, StatusSkipped}

func StatusFromString(s string) *Status {
	for _, t := range AllStatuses {
		if t.Key == s {
			return t
		}
	}
	return nil
}

func (t *Status) String() string {
	return t.Key
}

func (t *Status) MarshalJSON() ([]byte, error) {
	return util.ToJSONBytes(t.Key, false), nil
}

func (t *Status) UnmarshalJSON(data []byte) error {
	var s string
	if err := util.FromJSON(data, &s); err != nil {
		return err
	}
	x := StatusFromString(s)
	*t = *x
	return nil
}
