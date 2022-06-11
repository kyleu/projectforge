package data

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/export/model"
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

func (d *File) ToModel() (*model.Model, error) {
	name := strings.TrimSuffix(d.Filename, ".dat")
	key := util.StringToSnake(name)
	ret := &model.Model{
		Name:          key,
		Package:       key,
		Description:   fmt.Sprintf("from file [%s]", d.Filename),
		Icon:          "star",
		TitleOverride: util.StringToTitle(name),
	}
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
	}

	return ret, nil
}

type DataFiles []*File

func (d DataFiles) Headers() []map[string]any {
	ret := make([]map[string]any, 0, len(d))
	for _, x := range d {
		ret = append(ret, map[string]any{"fn": x.Filename, "cols": x.Fields})
	}
	return ret
}