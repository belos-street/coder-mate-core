package javascript

import "code-mate-core/core"

var Spec = core.TokenizerSpec{
	InitialState:  StateGlobal,
	Rules:         grammarRules,
	FallbackScope: ScopeDefault,
}
