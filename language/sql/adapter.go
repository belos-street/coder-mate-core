package sql

import "code-mate-core/core"

type adapter struct{}

func (a adapter) ID() string        { return "sql" }
func (a adapter) Aliases() []string { return []string{"postgresql", "pgsql"} }
func (a adapter) Parse(code string) core.TokenStream {
	return core.Parse(code, &Spec)
}

var Language core.LanguageAdapter = adapter{}
