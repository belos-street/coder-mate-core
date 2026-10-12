package rust

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("rust", []string{"rs"}, func(code string) core.TokenStream {
	return core.Parse(code, &spec)
})
