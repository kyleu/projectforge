package menu

var Separator = &Item{}

type Item struct {
	Key         string `json:"key"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Route       string `json:"route,omitempty"`
	Children    Items  `json:"children,omitempty"`
}

func (i *Item) AddChild(child *Item) {
	i.Children = append(i.Children, child)
}

func (i *Item) Desc() string {
	if i.Description != "" {
		return i.Title + ": " + i.Description
	}
	return i.Title
}

type Items []*Item

func (i Items) Get(key string) *Item {
	for _, item := range i {
		if item.Key == key {
			return item
		}
	}
	return nil
}

func (i Items) GetByPath(path []string) *Item {
	if len(path) == 0 {
		return nil
	}
	ret := i.Get(path[0])
	if ret == nil {
		return nil
	}
	if len(path) > 1 {
		return ret.Children.GetByPath(path[1:])
	}
	return ret
}

func (i Items) Keys() []string {
	ret := make([]string, 0, len(i))
	for _, x := range i {
		ret = append(ret, x.Key)
	}
	return ret
}
