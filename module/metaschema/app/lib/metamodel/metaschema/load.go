package metaschema

import (
	"context"
	"slices"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/filesystem"
	"{{{ .Package }}}/app/lib/jsonschema"
	"{{{ .Package }}}/app/lib/metamodel"
	"{{{ .Package }}}/app/util"
)

func LoadSchemas(
	ctx context.Context, key string, args *metamodel.Args, extraPaths []string, logger util.Logger, filter ...string,
) (*jsonschema.Collection, error) {
	ret := jsonschema.NewCollection()
	for _, x := range args.Enums {
		if len(filter) > 0 && !slices.Contains(filter, x.Name) {
			continue
		}
		_, err := ExportEnum(x, ret, args)
		if err != nil {
			return nil, err
		}
	}
	for _, x := range args.Models {
		if len(filter) > 0 && !slices.Contains(filter, x.Name) {
			continue
		}
		_, err := ExportModel(x, ret, args)
		if err != nil {
			return nil, err
		}
	}
	for _, x := range extraPaths {
		err := parseExtraPath(ctx, x, ret, logger)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse [%s]", x)
		}

	}
	return ret, nil
}

func parseExtraPath(ctx context.Context, pth string, coll *jsonschema.Collection, logger util.Logger) error {
	fs, _ := filesystem.NewFileSystem(".", true, "")
	if fs.IsDir(pth) {
		files := fs.ListJSON(pth, nil, false, logger)
		for _, fn := range files {
			if err := parseExtraPath(ctx, util.StringPath(pth, fn), coll, logger); err != nil {
				return err
			}
		}
	} else {
		b, err := fs.ReadFile(pth)
		if err != nil {
			return err
		}
		sch, err := jsonschema.FromJSON(b)
		if err != nil {
			return err
		}
		coll.AddSchema(sch)
	}
	return nil
}
