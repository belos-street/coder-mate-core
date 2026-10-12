package yaml

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("yaml", []string{"yml"}, func(code string) core.TokenStream {
	return core.Parse(code, &spec)
})
