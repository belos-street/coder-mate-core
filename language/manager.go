package language

import "code-mate-core/core"

// EnsureBuiltinsRegistered 幂等注册内置语言：任一内置语言已存在时跳过；
// 注册表被 ClearRegistry 清空（如测试场景）后再次调用会自动重建。
// 存在性检查与注册均由 core 注册表锁保护，重复注册幂等，并发调用安全。
func EnsureBuiltinsRegistered() {
	if _, ok := core.GetLanguage(builtinLanguages[0].ID()); !ok {
		RegisterBuiltins()
	}
}

func GetLanguage(id string) (core.LanguageAdapter, bool) {
	EnsureBuiltinsRegistered()
	return core.GetLanguage(id)
}

func ListLanguages() []core.LanguageAdapter {
	EnsureBuiltinsRegistered()
	return core.ListLanguages()
}

func Tokenize(code, langID string) (core.TokenStream, error) {
	EnsureBuiltinsRegistered()
	return core.Tokenize(code, langID)
}

// ResetBuiltins 清空语言注册表，供测试重建场景使用。
func ResetBuiltins() {
	core.ClearRegistry()
}
