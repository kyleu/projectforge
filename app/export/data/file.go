package data

import (
	"fmt"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/export/model"
)

type DataFile struct {
	Filename string  `json:"filename"`
	Fields   Fields  `json:"header"`
	Data     [][]any `json:"data"`
}

func (d *DataFile) String() string {
	return fmt.Sprintf("%s [%d fields]: %d rows", d.Filename, len(d.Fields), len(d.Data))
}

func (d *DataFile) ToModel() (*model.Model, error) {
	ret := &model.Model{}
	for _, col := range d.Fields {
		x, err := exportCol(col)
		if err != nil {
			return nil, errors.Wrapf(err, "can't export model [%s]", d.Filename)
		}
		ret.Columns = append(ret.Columns, x)
	}
	return ret, nil
}

type DataFiles []*DataFile

func (d DataFiles) Headers() []map[string]any {
	ret := make([]map[string]any, 0, len(d))
	for _, x := range d {
		ret = append(ret, map[string]any{"fn": x.Filename, "cols": x.Fields})
	}
	return ret
}
