// Package filesystem - Content managed by Project Forge, see [projectforge.md] for details.
package filesystem

import (
	"path"

	"projectforge.dev/projectforge/app/util"
)

func (f *FileSystem) listNodes(pth string, ign []string, logger util.Logger, tags ...string) (Nodes, error) {
	files := f.ListFiles(pth, ign, logger)
	nodes := make(Nodes, 0, len(files))
	for _, fi := range files {
		x := path.Join(pth, fi.Name)
		if fi.IsDir {
			kids, err := f.listNodes(x, ign, logger)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, &Node{Name: fi.Name, Dir: true, Children: kids, Tags: tags})
		} else {
			nodes = append(nodes, &Node{Name: fi.Name, Tags: tags, Size: f.Size(x)})
		}
	}
	return nodes.Sort(), nil
}

func (f *FileSystem) ListTree(cfg util.ValueMap, pth string, ign []string, logger util.Logger, tags ...string) (*Tree, error) {
	nodes, err := f.listNodes(pth, ign, logger, tags...)
	if err != nil {
		return nil, err
	}
	return &Tree{Config: cfg, Nodes: nodes}, nil
}
