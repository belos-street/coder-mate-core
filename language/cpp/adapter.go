package cpp

import "code-mate-core/core"

var Language = core.NewSimpleAdapter("cpp", []string{"c++", "cc", "cxx", "hpp", "hxx", "hh"}, func(code string) core.TokenStream {
	return core.Parse(code, &spec)
})
