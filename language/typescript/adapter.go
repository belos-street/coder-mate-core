package typescript

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("typescript", []string{"ts"}, func(code string) core.TokenStream {
	return core.Parse(code, spec)
})
