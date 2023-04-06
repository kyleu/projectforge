package inject

import (
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/model"
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
		case grpcCall:
			err = GRPC(f, args)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
