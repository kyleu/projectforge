package inject

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/util"
)

func Services(f *file.File, args *model.Args) error {
	if len(args.Models) == 0 {
		return nil
	}

	svcSize := 0
	for _, m := range args.Models {
		if len(m.Proper()) > svcSize {
			svcSize = len(m.Proper())
		}
	}

	svcs := make([]string, 0, len(args.Models))
	refs := make([]string, 0, len(args.Models))
	for _, m := range args.Models {
		svcs = append(svcs, fmt.Sprintf("%s *%s.Service", util.StringPad(m.Proper(), svcSize), m.Package))
		refs = append(refs, fmt.Sprintf("%s %s.NewService(st.DB, st.Logger),", util.StringPad(m.Proper()+":", svcSize+1), m.Package))
	}
	svcTxt := fmt.Sprintf("\n\t%s\n\t// ", strings.Join(svcs, "\n\t"))
	refTxt := fmt.Sprintf("\n\t\t%s\n\t\t// ", strings.Join(refs, "\n\t\t"))
	content := map[string]string{"services": svcTxt, "refs": refTxt}
	return file.Inject(f, content)
}
