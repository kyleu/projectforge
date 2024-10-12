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

func (s *Status) String() string {
	return s.Key
}

func (s *Status) Matches(x *Status) bool {
	return s.Key == x.Key
}

func (s *Status) StringFor(act string) string {
	if act == "audit" {
		switch s.Key {
		case StatusDifferent.Key:
			return "Invalid Header"
		case StatusMissing.Key:
			return "Empty Folder"
		default:
			return util.StringToTitle(s.Key)
		}
	}
	return util.StringToTitle(s.Key)
}

func (s *Status) MarshalJSON() ([]byte, error) {
	return util.ToJSONBytes(s.Key, false), nil
}

func (s *Status) UnmarshalJSON(data []byte) error {
	str, err := util.FromJSONString(data)
	if err != nil {
		return err
	}
	x := StatusFromString(str)
	*s = *x
	return nil
}
