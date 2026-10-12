package core

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

type LanguageAdapter interface {
	ID() string
	Aliases() []string
	Parse(code string) TokenStream
}

var (
	registryMu       sync.RWMutex
	languageRegistry = make(map[string]LanguageAdapter)
)

func normalizeID(id string) string {
	return strings.TrimSpace(strings.ToLower(id))
}

func ensureValidID(id string) (string, error) {
	normalized := normalizeID(id)
	if normalized == "" {
		return "", fmt.Errorf("language id cannot be empty")
	}
	return normalized, nil
}

func RegisterLanguage(lang LanguageAdapter) error {
	canonicalID, err := ensureValidID(lang.ID())
	if err != nil {
		return err
	}
	keys := []string{canonicalID}
	for _, alias := range lang.Aliases() {
		aliasID, err := ensureValidID(alias)
		if err != nil {
			return err
		}
		keys = append(keys, aliasID)
	}

	registryMu.Lock()
	defer registryMu.Unlock()

	for _, key := range keys {
		existing, ok := languageRegistry[key]
		if ok && !reflect.DeepEqual(existing, lang) {
			return fmt.Errorf("language id or alias %q is already registered", key)
		}
	}

	for _, key := range keys {
		languageRegistry[key] = lang
	}

	return nil
}

func GetLanguage(id string) (LanguageAdapter, bool) {
	normalized, err := ensureValidID(id)
	if err != nil {
		return nil, false
	}

	registryMu.RLock()
	defer registryMu.RUnlock()
	lang, ok := languageRegistry[normalized]
	return lang, ok
}

func ListLanguages() []LanguageAdapter {
	registryMu.RLock()
	defer registryMu.RUnlock()

	seen := make(map[string]bool)
	var languages []LanguageAdapter

	for _, lang := range languageRegistry {
		key := normalizeID(lang.ID())
		if seen[key] {
			continue
		}
		seen[key] = true
		languages = append(languages, lang)
	}

	// 按 ID 排序，保证返回顺序确定（Go map 遍历顺序随机）
	sort.Slice(languages, func(i, j int) bool {
		return languages[i].ID() < languages[j].ID()
	})

	return languages
}

func Tokenize(code, langID string) (TokenStream, error) {
	if langID == "" {
		return nil, fmt.Errorf("language id cannot be empty")
	}
	lang, ok := GetLanguage(langID)
	if !ok {
		return nil, fmt.Errorf("language %q is not registered", normalizeID(langID))
	}
	return lang.Parse(code), nil
}

func ClearRegistry() {
	registryMu.Lock()
	defer registryMu.Unlock()
	languageRegistry = make(map[string]LanguageAdapter)
}
