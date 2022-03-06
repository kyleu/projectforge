package inject

import (
	"projectforge.dev/app/export/model"
	"projectforge.dev/app/file"
)

func All(args *model.Args, files file.Files) error {
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
		case "app/controller/routes.go":
			err = Routes(f, args)
		case "app/controller/menu.go":
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
