package model

type Models []*Model

func (m Models) Get(n string) *Model {
	for _, x := range m {
		if x.Name == n {
			return x
		}
	}
	return nil
}

func (m Models) ReverseRelations(t string) Relations {
	var rels Relations
	for _, x := range m {
		for _, rel := range x.Relations {
			if rel.Table == t {
				rels = append(rels, rel.Reverse(x.Name))
			}
		}
	}
	return rels
}
