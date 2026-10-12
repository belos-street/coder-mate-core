package java

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("java", []string{"jav"}, func(code string) core.TokenStream {
	return core.Parse(code, &spec)
})
