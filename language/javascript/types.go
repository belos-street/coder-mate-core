package javascript

// Grammar states
const (
	StateGlobal                = "global"
	StateMultilineComment      = "multiline-comment"
	StateStringDouble          = "string-double"
	StateStringSingle          = "string-single"
	StateStringBacktick        = "string-backtick"
	StateTemplateInterpolation = "template-interpolation"
	StateImportDynamic         = "import-dynamic"
)

// Token scopes
const (
	ScopeCommentLine        = "comment.line.double-slash.js"
	ScopeCommentBlock       = "comment.block.js"
	ScopeKeywordControl     = "keyword.control.js"
	ScopeKeywordAsync       = "keyword.control.async.js"
	ScopeKeywordClass       = "keyword.control.class.js"
	ScopeKeywordModule      = "keyword.control.module.js"
	ScopeKeywordImport      = "keyword.control.import.js"
	ScopeKeywordDeclaration = "keyword.declaration.js"
	ScopeBoolean            = "constant.language.boolean.js"
	ScopeNull               = "constant.language.null.js"
	ScopeLanguageConstant   = "constant.language.js"
	ScopeNumeric            = "constant.numeric.js"
	ScopeGlobalThis             = "variable.language.global-this.js"
	ScopeVariableThis           = "variable.language.this.js"
	ScopeIdentifier             = "variable.identifier.js"
	ScopeStringDouble       = "string.quoted.double.js"
	ScopeStringSingle       = "string.quoted.single.js"
	ScopeStringBacktick     = "string.quoted.backtick.js"
	ScopeOperator           = "operator.js"
	ScopeOptionalChaining   = "operator.optional-chaining.js"
	ScopeNullishCoalescing  = "operator.nullish-coalescing.js"
	ScopeArrowFunction      = "operator.arrow-function.js"
	ScopeTemplateExpression = "punctuation.definition.template-expression.js"
	ScopeSupportPromise     = "support.function.promise.js"
	ScopeDefault            = "default"
)
