package css

import "code-mate-core/core"

type adapter struct{}

func (a adapter) ID() string        { return "css" }
func (a adapter) Aliases() []string { return nil }
func (a adapter) Parse(code string) core.TokenStream {
	return core.Parse(code, &Spec)
}

var Language core.LanguageAdapter = adapter{}
