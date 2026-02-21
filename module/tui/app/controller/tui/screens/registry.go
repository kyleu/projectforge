package screens

import (
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/menu"
)

type Registry struct {
	screens map[string]Screen
	menu    menu.Items
}

func NewRegistry() *Registry {
	return &Registry{screens: map[string]Screen{}}
}

func (r *Registry) AddScreen(scr Screen) {
	r.screens[scr.Key()] = scr
}

func (r *Registry) Register(item *menu.Item, scr Screen) {
	if item.Route == "" {
		item.Route = scr.Key()
	}
	r.menu = append(r.menu, item)
	r.AddScreen(scr)
}

func (r *Registry) Screen(key string) (Screen, error) {
	ret, ok := r.screens[key]
	if !ok {
		return nil, errors.Errorf("screen [%s] is not registered", key)
	}
	return ret, nil
}

func (r *Registry) MustScreen(key string) Screen {
	s, err := r.Screen(key)
	if err != nil {
		panic(err)
	}
	return s
}

func (r *Registry) Menu() menu.Items {
	return r.menu
}
