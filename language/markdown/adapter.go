package markdown

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("markdown", []string{"md"}, func(code string) core.TokenStream {
	return core.Parse(code, &spec)
})
