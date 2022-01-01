package model

func ColToString(c *Column, prefix string) string {
	return TypeToString(c.Type, prefix+c.Proper())
}

func TypeToString(t *Type, prop string) string {
	switch t.Key {
	case TypeUUID.Key:
		return prop + ".String()"
	default:
		return prop
	}
}

func ColToViewString(c *Column, prefix string) string {
	return TypeToViewString(c.Type, prefix+c.Proper(), c.Nullable)
}

func TypeToViewString(t *Type, prop string, nullable bool) string {
	ret := t.ToGoString(prop)
	switch t.Key {
	case TypeTimestamp.Key:
		if nullable {
			return "{%%= components.DisplayTimestamp(" + ret + ") %%}"
		}
		return "{%%= components.DisplayTimestamp(&" + ret + ") %%}"
	default:
		return "{%%s " + ret + " %%}"
	}
}
