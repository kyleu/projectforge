package inject

import (
	"context"

	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/util"
)

func All(ctx context.Context, args *model.Args, files file.Files, logger util.Logger) error {
	if args == nil {
		return nil
	}
	grpcCall := ""
	if args.HasModule("grpc") {
		grpcPackage := args.Config.GetStringOpt("grpcPackage")
		if grpcPackage == "" {
			grpcPackage = "grpc"
		}
		grpcCall = "app/" + grpcPackage + "/handle.go"
	}
	for _, f := range files {
		var err error
		switch f.FullPath() {
		case "app/services.go":
			err = Services(f, args)
		case "app/controller/cmenu/menu.go":
			err = Menu(f, args)
		case "app/lib/search/search.go":
			err = Search(f, args)
		case grpcCall:
			err = GRPC(f, args)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
