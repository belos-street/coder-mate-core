package javascript

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("javascript", []string{"js"}, func(code string) core.TokenStream {
	return core.Parse(code, &spec)
})
