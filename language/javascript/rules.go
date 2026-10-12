package javascript

import (
	"regexp"

	"code-mate-core/core"
)

func GetRules() map[string][]core.GrammarRule {
	// 深拷贝，防止调用方篡改内部规则表
	cloned := make(map[string][]core.GrammarRule, len(grammarRules))
	for k, v := range grammarRules {
		rules := make([]core.GrammarRule, len(v))
		copy(rules, v)
		cloned[k] = rules
	}
	return cloned
}

var grammarRules = map[string][]core.GrammarRule{
	// ==================== global 状态 ====================
	StateGlobal: {
		// === 1. 注释规则 ===
		{
			Regex: regexp.MustCompile(`^\/\/.*`),
			Scope: ScopeCommentLine,
		},
		{
			Regex:     regexp.MustCompile(`^\/\*`),
			Scope:     ScopeCommentBlock,
			PushState: StateMultilineComment,
		},

		// === 2. 模块化语法（ES2020）===
		{
			Regex:     regexp.MustCompile(`^(import)\s*\(`),
			Scope:     ScopeKeywordImport,
			PushState: StateImportDynamic,
		},
		{
			Regex: regexp.MustCompile(`^(export)\s*\*\s*as\s*(\w+)\s*from`),
			Scope: ScopeKeywordModule,
		},
		{
			Regex: regexp.MustCompile(`^(import|export)\b`),
			Scope: ScopeKeywordModule,
		},

		// === 3. ES2020 关键字/全局变量 ===
		{
			Regex: regexp.MustCompile(`^(globalThis)\b`),
			Scope: ScopeGlobalThis,
		},
		{
			Regex: regexp.MustCompile(`^(Promise\.allSettled)\b`),
			Scope: ScopeSupportPromise,
		},

		// === 4. 控制流/声明关键字 ===
		{
			Regex: regexp.MustCompile(`^(if|else|for|while|do|switch|case|break|continue|return|throw|try|catch|finally)\b`),
			Scope: ScopeKeywordControl,
		},
		{
			Regex: regexp.MustCompile(`^(async|await)\b`),
			Scope: ScopeKeywordAsync,
		},
		{
			Regex: regexp.MustCompile(`^(typeof|new|delete|void|yield|instanceof)\b`),
			Scope: ScopeKeywordControl,
		},
		{
			Regex: regexp.MustCompile(`^(this|super)\b`),
			Scope: ScopeVariableThis,
		},
		{
			Regex: regexp.MustCompile(`^(class|extends|static|constructor)\b`),
			Scope: ScopeKeywordClass,
		},
		{
			Regex: regexp.MustCompile(`^(function|var|let|const)\b`),
			Scope: ScopeKeywordDeclaration,
		},
		{
			Regex: regexp.MustCompile(`^(true|false)\b`),
			Scope: ScopeBoolean,
		},
		{
			Regex: regexp.MustCompile(`^(null|undefined)\b`),
			Scope: ScopeNull,
		},

		// === 5. 运算符（ES2020新增）===
		{
			Regex: regexp.MustCompile(`^(\?\.)`),
			Scope: ScopeOptionalChaining,
		},
		{
			Regex: regexp.MustCompile(`^(\?\?)`),
			Scope: ScopeNullishCoalescing,
		},
		{
			Regex: regexp.MustCompile(`^=>`),
			Scope: ScopeArrowFunction,
		},
		{
			Regex: regexp.MustCompile(`^(===|!==|==|!=|>=|<=|&&|\|\||\+\+|--|\*\*|\+|-|\*|\/|%|=|>|<|!|\.|,|:|;|\(|\)|\{|\}|\[|\])`),
			Scope: ScopeOperator,
		},

		// === 6. 值类型 ===
		{
			Regex: regexp.MustCompile(`^(0b[01]+n?|0o[0-7]+n?|0x[0-9a-fA-F]+n?|\d+\.?\d*(?:[eE][+-]?\d+)?n?|\.\d+(?:[eE][+-]?\d+)?n?)`),
			Scope: ScopeNumeric,
		},
		{
			Regex:     regexp.MustCompile("^`"),
			Scope:     ScopeStringBacktick,
			PushState: StateStringBacktick,
		},
		{
			Regex:     regexp.MustCompile(`^"`),
			Scope:     ScopeStringDouble,
			PushState: StateStringDouble,
		},
		{
			Regex:     regexp.MustCompile(`^'`),
			Scope:     ScopeStringSingle,
			PushState: StateStringSingle,
		},

		// === 7. 标识符（普通变量/函数名）===
		{
			Regex: regexp.MustCompile(`^[a-zA-Z_$][a-zA-Z0-9_$]*`),
			Scope: ScopeIdentifier,
		},
	},

	// ==================== multiline-comment 状态 ====================
	StateMultilineComment: {
		{
			Regex:    regexp.MustCompile(`^[\s\S]*?\*\/`),
			Scope:    ScopeCommentBlock,
			PopState: true,
		},
		{
			Regex: regexp.MustCompile(`^[\s\S]*`),
			Scope: ScopeCommentBlock,
		},
	},

	// ==================== string-double 状态 ====================
	StateStringDouble: {
		{
			Regex:    regexp.MustCompile(`^[^\\"]*(?:\\.[^\\"]*)*"`),
			Scope:    ScopeStringDouble,
			PopState: true,
		},
		{
			Regex: regexp.MustCompile(`^.*`),
			Scope: ScopeStringDouble,
		},
	},

	// ==================== string-single 状态 ====================
	StateStringSingle: {
		{
			Regex:    regexp.MustCompile(`^[^\\']*(?:\\.[^\\']*)*'`),
			Scope:    ScopeStringSingle,
			PopState: true,
		},
		{
			Regex: regexp.MustCompile(`^.*`),
			Scope: ScopeStringSingle,
		},
	},

	// ==================== string-backtick 状态 ====================
	StateStringBacktick: {
		{
			Regex:     regexp.MustCompile(`^\$\{`),
			Scope:     ScopeTemplateExpression,
			PushState: StateTemplateInterpolation,
		},
		{
			Regex: regexp.MustCompile("^[^\\\\\x60\\$]+"),
			Scope: ScopeStringBacktick,
		},
		{
			Regex: regexp.MustCompile(`^\\.`),
			Scope: ScopeStringBacktick,
		},
		{
			Regex: regexp.MustCompile(`^\$`),
			Scope: ScopeStringBacktick,
		},
		{
			Regex:    regexp.MustCompile("^\x60"),
			Scope:    ScopeStringBacktick,
			PopState: true,
		},
	},

	// ==================== template-interpolation 状态 ====================
	StateTemplateInterpolation: {
		{
			Regex:    regexp.MustCompile(`^}`),
			Scope:    ScopeTemplateExpression,
			PopState: true,
		},
		{
			Regex: regexp.MustCompile(`^(if|else|const|let|var|async|await)\b`),
			Scope: ScopeKeywordControl,
		},
		{
			Regex: regexp.MustCompile(`^(true|false|null|undefined)\b`),
			Scope: ScopeLanguageConstant,
		},
		{
			Regex: regexp.MustCompile(`^(\?\.|\?\?)`),
			Scope: ScopeOperator,
		},
		{
			Regex:     regexp.MustCompile(`^"`),
			Scope:     ScopeStringDouble,
			PushState: StateStringDouble,
		},
		{
			Regex:     regexp.MustCompile(`^'`),
			Scope:     ScopeStringSingle,
			PushState: StateStringSingle,
		},
		{
			Regex: regexp.MustCompile(`^\d+n?`),
			Scope: ScopeNumeric,
		},
		{
			Regex: regexp.MustCompile(`^[a-zA-Z_$][a-zA-Z0-9_$]*`),
			Scope: ScopeIdentifier,
		},
		{
			Regex: regexp.MustCompile(`^[+\-*/=<>!&|.,;:()[\]]`),
			Scope: ScopeOperator,
		},
		{
			Regex: regexp.MustCompile(`^\{`),
			Scope: ScopeOperator,
		},
		{
			Regex: regexp.MustCompile(`^.`),
			Scope: ScopeIdentifier,
		},
	},

	// ==================== import-dynamic 状态 ====================
	StateImportDynamic: {
		{
			Regex:    regexp.MustCompile(`^\)`),
			Scope:    ScopeOperator,
			PopState: true,
		},
		{
			Regex:     regexp.MustCompile(`^"`),
			Scope:     ScopeStringDouble,
			PushState: StateStringDouble,
		},
		{
			Regex:     regexp.MustCompile(`^'`),
			Scope:     ScopeStringSingle,
			PushState: StateStringSingle,
		},
		{
			Regex: regexp.MustCompile(`^[a-zA-Z_$][a-zA-Z0-9_$]*`),
			Scope: ScopeIdentifier,
		},
		{
			Regex: regexp.MustCompile(`^[+\-*/=<>!&|.,;:()[\]{}]`),
			Scope: ScopeOperator,
		},
		{
			Regex: regexp.MustCompile(`^.`),
			Scope: ScopeIdentifier,
		},
	},
}
