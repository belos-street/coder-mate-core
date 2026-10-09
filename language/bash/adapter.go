package bash

import "code-mate-core/core"

type adapter struct{}

func (a adapter) ID() string        { return "bash" }
func (a adapter) Aliases() []string { return []string{"sh", "shell"} }
func (a adapter) Parse(code string) core.TokenStream {
	return core.Parse(code, &Spec)
}

var Language core.LanguageAdapter = adapter{}
