package rust

import "code-mate-core/core"

var spec = core.TokenizerSpec{
	InitialState:  GrammarStateGlobal,
	Rules:         grammarRules,
	FallbackScope: ScopeDefault,
}
