package sql

import "projectforge.dev/projectforge/app/project/export/files/helper"

const endCall = "() %%}"

func sqlFunc(n string) string {
	return helper.TextSQLComment + "{%% func " + n + endCall
}

func sqlCall(n string) string {
	return helper.TextSQLComment + "{%%= " + n + endCall
}

func sqlEnd() string {
	return helper.TextSQLComment + helper.TextEndFunc
}
