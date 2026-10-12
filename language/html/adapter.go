package html

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("html", []string{"htm"}, func(code string) core.TokenStream {
	return core.Parse(code, &spec)
})
