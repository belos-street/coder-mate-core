package css

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("css", nil, func(code string) core.TokenStream {
	return core.Parse(code, &spec)
})
