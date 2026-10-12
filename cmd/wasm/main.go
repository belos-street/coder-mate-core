//go:build js && wasm

package main

import (
	"sync"
	"syscall/js"

	codematecore "code-mate-core"
	"code-mate-core/language"
	"code-mate-core/theme"
)

var (
	mu           sync.Mutex
	nextID       int
	highlighters = map[int]codematecore.Highlighter{}
	defaultHL    codematecore.Highlighter
)

func main() {
	language.RegisterBuiltins()
	theme.RegisterBuiltInThemes()

	defaultHL = codematecore.NewHighlighter(codematecore.HighlighterOptions{})

	js.Global().Set("coderMateHighlight", js.FuncOf(jsCoderMateHighlight))

	coderMate := map[string]interface{}{
		"createHighlighter": js.FuncOf(jsCreateHighlighter),
		"highlight":         js.FuncOf(jsHighlight),
		"updateTheme":       js.FuncOf(jsUpdateTheme),
		"dispose":           js.FuncOf(jsDispose),
		"getThemes":         js.FuncOf(jsGetThemes),
		"getLanguages":      js.FuncOf(jsGetLanguages),
	}
	js.Global().Set("coderMate", js.ValueOf(coderMate))

	select {}
}

func jsCoderMateHighlight(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return errorResult("expected 2 arguments: code, lang")
	}
	code := args[0].String()
	lang := args[1].String()

	result, err := defaultHL.CodeToHtml(code, lang)
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(result)
}

func jsCreateHighlighter(this js.Value, args []js.Value) interface{} {
	opts := codematecore.HighlighterOptions{}
	if len(args) > 0 && !args[0].IsUndefined() && !args[0].IsNull() {
		obj := args[0]
		if themeVal := obj.Get("theme"); !themeVal.IsUndefined() {
			opts.Theme = themeVal.String()
		}
		if preStyle := obj.Get("preStyle"); !preStyle.IsUndefined() {
			opts.PreStyle = preStyle.String()
		}
		if prefix := obj.Get("lineClassPrefix"); !prefix.IsUndefined() {
			opts.LineClassPrefix = prefix.String()
		}
	}

	hl := codematecore.NewHighlighter(opts)

	mu.Lock()
	nextID++
	id := nextID
	highlighters[id] = hl
	mu.Unlock()

	return id
}

func jsHighlight(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return errorResult("expected 3 arguments: id, code, lang")
	}
	id := args[0].Int()
	code := args[1].String()
	lang := args[2].String()

	mu.Lock()
	hl, ok := highlighters[id]
	mu.Unlock()

	if !ok {
		return errorResult("highlighter not found")
	}

	result, err := hl.CodeToHtml(code, lang)
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(result)
}

func jsUpdateTheme(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return errorResult("expected 2 arguments: id, theme")
	}
	id := args[0].Int()
	themeName := args[1].String()

	mu.Lock()
	hl, ok := highlighters[id]
	mu.Unlock()

	if !ok {
		return errorResult("highlighter not found")
	}

	err := hl.UpdateTheme(themeName)
	if err != nil {
		return errorResult(err.Error())
	}
	return nil
}

func jsDispose(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return nil
	}
	id := args[0].Int()

	mu.Lock()
	delete(highlighters, id)
	mu.Unlock()

	return nil
}

func jsGetThemes(this js.Value, args []js.Value) interface{} {
	themes := theme.ListThemes()
	result := make([]interface{}, len(themes))
	for i, t := range themes {
		result[i] = map[string]interface{}{
			"id":          t.ID,
			"displayName": t.DisplayName,
		}
	}
	return js.ValueOf(result)
}

func jsGetLanguages(this js.Value, args []js.Value) interface{} {
	langs := language.ListLanguages()
	result := make([]interface{}, len(langs))
	for i, l := range langs {
		aliases := l.Aliases()
		jsAliases := make([]interface{}, len(aliases))
		for j, a := range aliases {
			jsAliases[j] = a
		}
		result[i] = map[string]interface{}{
			"id":      l.ID(),
			"aliases": jsAliases,
		}
	}
	return js.ValueOf(result)
}

func successResult(html string) interface{} {
	return js.ValueOf(map[string]interface{}{
		"html": html,
	})
}

func errorResult(msg string) interface{} {
	return js.ValueOf(map[string]interface{}{
		"error": msg,
	})
}
