package json

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("json", nil, func(code string) core.TokenStream {
	return core.Parse(code, &spec)
})
