package export

var Helpers = func() map[string][]string {
	ret := map[string][]string{}
	add := func(k string, v string) {
		ret[k] = []string{v}
	}

	add("enum.values", "Possible values for this enumeration")
	add("model.name", "The name, usually corresponding to the name of your database table")
	add("model.package", "package")
	add("model.group", "The model's group, as a path, like [foo/bar]")
	add("model.description", "description")
	add("model.icon", "icon")
	add("model.ordering", "ordering")
	add("model.view", "view")
	add("model.search", "search")
	add("model.tags", "tags")
	add("model.titleOverride", "titleOverride")
	add("model.properOverride", "properOverride")
	add("model.config", "config")
	add("model.columns", "columns")
	add("model.relations", "relations")
	add("model.indexes", "indexes")

	return ret
}()
