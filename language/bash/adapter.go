package bash

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("bash", []string{"sh", "shell"}, func(code string) core.TokenStream {
	return core.Parse(code, &spec)
})
