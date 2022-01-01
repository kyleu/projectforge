package inject

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
)

func Menu(f *file.File, args *model.Args) error {
	if len(args.Models) == 0 {
		return nil
	}
	out := make([]string, 0, len(args.Models))
	msg := `&menu.Item{Key: %q, Title: %q, Description: %q, Icon: %q, Route: "/%s"},`
	for _, m := range args.Models {
		l := fmt.Sprintf(msg, m.Package, m.ProperPlural(), m.Description, m.Icon, m.Package)
		out = append(out, l)
	}
	content := map[string]string{"codegen": "\n\tret = append(ret,\n\t\t" + strings.Join(out, "\n\t\t") + "\n\t)\n\t// "}
	return file.Inject(f, content)
}
