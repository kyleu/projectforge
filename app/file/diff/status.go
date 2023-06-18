package diff

import (
	"github.com/samber/lo"
	"projectforge.dev/projectforge/app/util"
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
	return lo.FindOrElse(AllStatuses, nil, func(t *Status) bool {
		return t.Key == s
	})
}

func (t *Status) String() string {
	return t.Key
}

func (t *Status) StringFor(act string) string {
	if act == "audit" {
		switch t.Key {
		case StatusDifferent.Key:
			return "Invalid Header"
		case StatusMissing.Key:
			return "Empty Folder"
		default:
			return util.StringToTitle(t.Key)
		}
	}
	return util.StringToTitle(t.Key)
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
