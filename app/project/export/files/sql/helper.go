package sql

func sqlFunc(n string) string {
	return "-- {%% func " + n + "() %%}"
}

func sqlCall(n string) string {
	return "-- {%%= " + n + "() %%}"
}

func sqlEnd() string {
	return "-- {%% endfunc %%}"
}
