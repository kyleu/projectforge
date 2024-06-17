package project

import "projectforge.dev/projectforge/app/util"

func (p *Project) DatabaseEngines() []string {
	ret := util.NewStringSlice(make([]string, 0, 4))
	if p.HasModule(util.DatabaseMySQL) {
		ret.Push(util.DatabaseMySQL)
	}
	if p.HasModule(util.DatabasePostgreSQL) {
		ret.Push(util.DatabasePostgreSQL)
	}
	if p.HasModule(util.DatabaseSQLite) {
		ret.Push(util.DatabaseSQLite)
	}
	if p.HasModule(util.DatabaseSQLServer) {
		ret.Push(util.DatabaseSQLServer)
	}
	return ret.Slice
}

func (p *Project) DatabaseEngineDefault() string {
	if p.Info != nil && p.Info.DatabaseEngine != "" {
		return p.Info.DatabaseEngine
	}
	engines := p.DatabaseEngines()
	if len(engines) == 1 {
		return engines[0]
	}
	if len(engines) > 1 {
		return util.DatabasePostgreSQL
	}
	return ""
}
