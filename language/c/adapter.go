package c

import "code-mate-core/core"

type adapter struct{}

func (a adapter) ID() string        { return "c" }
func (a adapter) Aliases() []string { return []string{"h", "c89", "c99", "c11"} }
func (a adapter) Parse(code string) core.TokenStream {
	return core.Parse(code, &Spec)
}

var Language core.LanguageAdapter = adapter{}
