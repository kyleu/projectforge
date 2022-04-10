package database

import (
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var serviceRegistry = map[string]*Service{}
var serviceRegistryMu = sync.Mutex{}

func register(s *Service, logger *zap.SugaredLogger) {
	serviceRegistryMu.Lock()
	defer serviceRegistryMu.Unlock()
	if _, ok := serviceRegistry[s.Key]; ok {
		logger.Warnf("double registration for database [%s]", s.Key)
	}
	serviceRegistry[s.Key] = s
}

func unregister(s *Service) {
	delete(serviceRegistry, s.Key)
}

func RegistryGet(key string) (*Service, error) {
	ret, ok := serviceRegistry[key]
	if !ok {
		return nil, errors.Errorf("no database registered with key [%s]", key)
	}
	return ret, nil
}

func RegistryKeys() []string {
	ret := maps.Keys(serviceRegistry)
	slices.Sort(ret)
	return ret
}
