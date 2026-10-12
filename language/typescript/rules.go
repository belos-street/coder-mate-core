package typescript

import (
	"regexp"
	"sync"

	"code-mate-core/core"
	js "code-mate-core/language/javascript"
)

var (
	grammarRules     map[string][]core.GrammarRule
	grammarRulesOnce sync.Once
)

var tsBuiltinTypes = regexp.MustCompile(`^(string|number|boolean|any|unknown|never|void|object|symbol|bigint|Array|Promise|Record|Partial|Pick|Omit|Readonly|Required)\b`)
var operatorJS = regexp.MustCompile(`^(\.\.\.|===|!==|==|!=|>=|<=|&&|\|\||\+\+|--|\+|-|\*|\/|%|=|>|<|!|&|\||\?|\.|,|:|;|\(|\)|\{|\}|\[|\])`)
var punctuation = regexp.MustCompile(`^(\.\.\.|[+\-*/=<>!&|.,;:?()[\]])`)

func getGrammarRules() map[string][]core.GrammarRule {
	grammarRulesOnce.Do(func() {
		grammarRules = buildGrammarRules()
	})
	return grammarRules
}

func buildGrammarRules() map[string][]core.GrammarRule {
	base := js.GetRules()
	global := base["global"]

	tsRules := []core.GrammarRule{
		{Regex: regexp.MustCompile(`^(type|interface|enum|namespace)\b`), Scope: ScopeKeywordDeclTypeTS, PushState: StateExpectTypeName},
		{Regex: regexp.MustCompile(`^(implements)\b`), Scope: ScopeKeywordDeclTypeTS, PushState: StateExpectTypeName},
		{Regex: regexp.MustCompile(`^(declare)\b`), Scope: ScopeKeywordDeclTypeTS},
		{Regex: regexp.MustCompile(`^(public|private|protected|readonly|abstract|override)\b`), Scope: ScopeKeywordModifierAccessTS},
		{Regex: regexp.MustCompile(`^(as)\b`), Scope: ScopeKeywordOperatorAssertionTS},
		{Regex: regexp.MustCompile(`^(satisfies|keyof|typeof|infer|is|asserts|in|extends)\b`), Scope: ScopeKeywordOperatorTypeTS},
		{Regex: tsBuiltinTypes, Scope: ScopeSupportBuiltinTS},
		{Regex: regexp.MustCompile(`^[A-Z][a-zA-Z0-9_$]*`), Scope: ScopeEntityTypeTS},
	}

	var filtered []core.GrammarRule
	for _, rule := range global {
		s := rule.Scope
		if s == js.ScopeKeywordClass {
			filtered = append(filtered, core.GrammarRule{
				Regex:     regexp.MustCompile(`^(class)\b`),
				Scope:     js.ScopeKeywordClass,
				PushState: StateClassAfterName,
			})

		} else if s == js.ScopeIdentifier {
			continue
		} else {
			filtered = append(filtered, rule)
		}
	}

	// TS 类型关键字须排在 JS 关键字规则之前，但注释与动态 import 规则必须
	// 保持更靠前——锚定「动态 import」规则定位插入点，替代对 JS 规则顺序的
	// 位置硬编码（insertPos := 3）
	insertPos := len(filtered)
	for i, rule := range filtered {
		if rule.PushState == js.StateImportDynamic {
			insertPos = i + 1
			break
		}
	}

	var merged []core.GrammarRule
	for i, rule := range filtered {
		if i == insertPos {
			merged = append(merged, tsRules...)
		}
		merged = append(merged, rule)
	}
	if insertPos >= len(filtered) {
		merged = append(merged, tsRules...)
	}

	merged = append(merged, core.GrammarRule{
		Regex: regexp.MustCompile(`^[a-zA-Z_$][a-zA-Z0-9_$]*`),
		Scope: js.ScopeIdentifier,
	})
	merged = append(merged, core.GrammarRule{
		Regex: punctuation,
		Scope: js.ScopeOperator,
	})
	merged = append(merged, core.GrammarRule{
		Regex: regexp.MustCompile(`^.`),
		// TS 参考实现的 global 状态无 catch-all（未匹配走 fallbackScope=default），
		// 这里保持一致：兜底给 default 而非 identifier
		Scope: js.ScopeDefault,
	})

	base["global"] = merged

	base[StateExpectTypeName] = []core.GrammarRule{
		{Regex: regexp.MustCompile(`^\s+`), Scope: js.ScopeDefault},
		{Regex: regexp.MustCompile(`^[a-zA-Z_$][a-zA-Z0-9_$]*`), Scope: ScopeEntityTypeTS, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: js.ScopeDefault, PopState: true},
	}
	base[StateClassAfterName] = []core.GrammarRule{
		{Regex: regexp.MustCompile(`^\s+`), Scope: js.ScopeDefault},
		{Regex: regexp.MustCompile(`^(implements)\b`), Scope: ScopeKeywordDeclTypeTS, PushState: StateExpectTypeName},
		{Regex: regexp.MustCompile(`^(extends)\b`), Scope: js.ScopeKeywordClass, PushState: StateExpectTypeName},
		{Regex: regexp.MustCompile(`^(satisfies|keyof|typeof|infer|is|asserts|in)\b`), Scope: ScopeKeywordOperatorTypeTS},
		{Regex: regexp.MustCompile(`^<`), Scope: js.ScopeOperator, PushState: StateClassTypeParam},
		{Regex: regexp.MustCompile(`^\{`), Scope: js.ScopeOperator, PopState: true},
		{Regex: regexp.MustCompile(`^[A-Z][a-zA-Z0-9_$]*`), Scope: ScopeEntityTypeTS},
		{Regex: regexp.MustCompile(`^[a-zA-Z_$][a-zA-Z0-9_$]*`), Scope: js.ScopeIdentifier},
		{Regex: regexp.MustCompile(`^(\.\.\.|[+\-*/=<>!&|.,;:?()[\]])`), Scope: js.ScopeOperator},
		{Regex: regexp.MustCompile(`^.`), Scope: js.ScopeDefault},
	}
	base[StateClassTypeParam] = []core.GrammarRule{
		{Regex: regexp.MustCompile(`^\s+`), Scope: js.ScopeDefault},
		{Regex: regexp.MustCompile(`^<`), Scope: js.ScopeOperator, PushState: StateClassTypeParam},
		{Regex: regexp.MustCompile(`^>`), Scope: js.ScopeOperator, PopState: true},
		{Regex: regexp.MustCompile(`^(satisfies|keyof|typeof|infer|is|asserts|in|extends)\b`), Scope: ScopeKeywordOperatorTypeTS},
		{Regex: tsBuiltinTypes, Scope: ScopeSupportBuiltinTS},
		{Regex: regexp.MustCompile("^`"), Scope: js.ScopeStringBacktick, PushState: js.StateStringBacktick},
		{Regex: regexp.MustCompile(`^"`), Scope: js.ScopeStringDouble, PushState: js.StateStringDouble},
		{Regex: regexp.MustCompile(`^'`), Scope: js.ScopeStringSingle, PushState: js.StateStringSingle},
		{Regex: regexp.MustCompile(`^[A-Z][a-zA-Z0-9_$]*`), Scope: ScopeEntityTypeTS},
		{Regex: regexp.MustCompile(`^[a-zA-Z_$][a-zA-Z0-9_$]*`), Scope: js.ScopeIdentifier},
		{Regex: operatorJS, Scope: js.ScopeOperator},
		{Regex: regexp.MustCompile(`^.`), Scope: js.ScopeDefault},
	}

	return base
}
