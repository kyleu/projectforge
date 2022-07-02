package data

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/project/export/model"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

type File struct {
	Filename string  `json:"filename"`
	Fields   Fields  `json:"header"`
	Data     [][]any `json:"data"`
}

func (d *File) String() string {
	return fmt.Sprintf("%s [%d fields]: %d rows", d.Filename, len(d.Fields), len(d.Data))
}

func (d *File) ToModel(tableIdx int, groups model.Groups) (*model.Model, error) {
	name := strings.TrimSuffix(d.Filename, ".dat")
	key := util.StringToSnake(name)
	desc := fmt.Sprintf("from file [%s]", d.Filename)
	title := util.StringToTitle(name)
	var group []string
	for _, g := range groups {
		if key != g.Key && strings.HasPrefix(key, g.Key) {
			group = []string{g.Key}
		}
	}
	var sd [][]any
	if len(d.Data) > 0 {
		sd = d.Data
	}
	ret := &model.Model{Name: key, Package: key, Group: group, Description: desc, Icon: "star", Tags: []string{"json"}, TitleOverride: title, SeedData: sd}
	for _, col := range d.Fields {
		x, err := exportColumn(col)
		if err != nil {
			return nil, errors.Wrapf(err, "can't export model [%s]", d.Filename)
		}
		ret.Columns = append(ret.Columns, x)
	}
	if len(ret.PKs()) == 0 {
		idCol := &model.Column{Name: "id", Type: types.NewUUID(), PK: true, Search: true, HelpString: "Synthetic identifier"}
		ret.Columns = append(model.Columns{idCol}, ret.Columns...)
		appended := make([][]any, 0, len(ret.SeedData))
		for rowIdx, x := range ret.SeedData {
			id := util.UUIDFromString(fmt.Sprintf("%08d-0000-0000-0000-%012d", tableIdx, rowIdx))
			if id == nil {
				return nil, errors.Errorf("invalid UUID")
			}
			appended = append(appended, append([]any{*id}, x...))
		}
		ret.SeedData = appended
	}

	return ret, nil
}

type Files []*File

func (f Files) Headers() []map[string]any {
	ret := make([]map[string]any, 0, len(f))
	for _, x := range f {
		ret = append(ret, map[string]any{"fn": x.Filename, "cols": x.Fields})
	}
	return ret
}
