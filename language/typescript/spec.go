package typescript

import "code-mate-core/core"

var spec = buildSpec()

// buildSpec 合并 JS 基础规则与 TS 扩展规则，构造 TS 词法规范
func buildSpec() *core.TokenizerSpec {
	return &core.TokenizerSpec{
		InitialState:  "global",
		Rules:         buildGrammarRules(),
		FallbackScope: "default",
	}
}
