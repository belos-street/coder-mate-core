package csharp

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("csharp", []string{"cs", "c#"}, func(code string) core.TokenStream {
	return core.Parse(code, &spec)
})
