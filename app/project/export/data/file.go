package data

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

type File struct {
	Filename string  `json:"filename"`
	Fields   Fields  `json:"header"`
	Data     [][]any `json:"data"`
}

func (f *File) String() string {
	return fmt.Sprintf("%s [%d fields]: %d rows", f.Filename, len(f.Fields), len(f.Data))
}

func (f *File) ToModel(tableIdx int, groups model.Groups) (*model.Model, error) {
	name := strings.TrimSuffix(f.Filename, ".dat")
	key := util.StringToSnake(name)
	desc := fmt.Sprintf("from file [%s]", f.Filename)
	title := util.StringToTitle(name)
	var group []string
	lo.ForEach(groups, func(g *model.Group, _ int) {
		if key != g.Key && strings.HasPrefix(key, g.Key) {
			group = []string{g.Key}
		}
	})
	var sd [][]any
	if len(f.Data) > 0 {
		sd = f.Data
	}
	ret := &model.Model{Name: key, Package: key, Group: group, Description: desc, Icon: "star", Tags: []string{"json"}, TitleOverride: title, SeedData: sd}
	for colIdx, col := range f.Fields {
		ex := make([]any, 0, 16)
		for rowIdx, row := range f.Data {
			if rowIdx > 15 {
				break
			}
			if len(row) > colIdx {
				ex = append(ex, row[colIdx])
			}
		}
		x, err := exportColumn(col, ex)
		if err != nil {
			return nil, errors.Wrapf(err, "can't export model [%s]", f.Filename)
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
	return lo.Map(f, func(x *File, _ int) map[string]any {
		return map[string]any{"fn": x.Filename, "cols": x.Fields}
	})
}
