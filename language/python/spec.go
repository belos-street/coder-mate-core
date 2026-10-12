package python

import "code-mate-core/core"

var spec = core.TokenizerSpec{
	InitialState:  StateGlobal,
	Rules:         grammarRules,
	FallbackScope: ScopeDefault,
}
