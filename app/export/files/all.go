package files

import (
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
)

func All(args *model.Args) (file.Files, error) {
	ret := make(file.Files, 0, len(args.Models)*10)
	for _, m := range args.Models {
		calls := file.Files{
			Model(m, args), DTO(m, args), Service(m, args), Controller(m, args),
			ViewList(m, args), ViewTable(m, args), ViewDetail(m, args), ViewEdit(m, args),
		}
		ret = append(ret, calls...)
		if args.HasModule("grpc") {
			f, err := GRPC(m, args)
			if err != nil {
				return nil, err
			}
			ret = append(ret, f)
		}
	}
	return ret, nil
}
