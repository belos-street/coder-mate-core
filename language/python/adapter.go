package python

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("python", []string{"py"}, func(code string) core.TokenStream {
	return core.Parse(code, &spec)
})
