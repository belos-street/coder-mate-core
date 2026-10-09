package json

import "code-mate-core/core"

type adapter struct{}

func (a adapter) ID() string        { return "json" }
func (a adapter) Aliases() []string { return nil }
func (a adapter) Parse(code string) core.TokenStream {
	return core.Parse(code, &Spec)
}

var Language core.LanguageAdapter = adapter{}
