package php

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("php", []string{"phtml", "php8"}, func(code string) core.TokenStream {
	return core.Parse(code, &spec)
})
