package gorule

import "code-mate-core/core"

type adapter struct{}

func (a adapter) ID() string        { return "go" }
func (a adapter) Aliases() []string { return []string{"golang"} }
func (a adapter) Parse(code string) core.TokenStream {
	return core.Parse(code, &Spec)
}

var Language core.LanguageAdapter = adapter{}
