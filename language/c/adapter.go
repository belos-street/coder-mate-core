package c

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("c", []string{"h", "c89", "c99", "c11"}, func(code string) core.TokenStream {
	return core.Parse(code, &spec)
})
