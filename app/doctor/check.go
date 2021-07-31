package doctor

func Check() Prognoses {
	var ret Prognoses
	for _, d := range AllDependencies {
		ret = append(ret, d.Check())
	}
	return ret
}
