package gorule

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("go", []string{"golang"}, func(code string) core.TokenStream {
	return core.Parse(code, &spec)
})
