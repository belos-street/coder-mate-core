package sql

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("sql", []string{"postgresql", "pgsql"}, func(code string) core.TokenStream {
	return core.Parse(code, &spec)
})
