package language

import (
	"fmt"

	"code-mate-core/core"
	"code-mate-core/language/bash"
	"code-mate-core/language/c"
	"code-mate-core/language/cpp"
	"code-mate-core/language/csharp"
	"code-mate-core/language/css"
	"code-mate-core/language/gorule"
	"code-mate-core/language/html"
	"code-mate-core/language/java"
	"code-mate-core/language/javascript"
	"code-mate-core/language/json"
	"code-mate-core/language/markdown"
	"code-mate-core/language/php"
	"code-mate-core/language/python"
	"code-mate-core/language/rust"
	"code-mate-core/language/sql"
	"code-mate-core/language/typescript"
	"code-mate-core/language/yaml"
)

var builtinLanguages = []core.LanguageAdapter{
	javascript.Language,
	typescript.Language,
	python.Language,
	bash.Language,
	c.Language,
	cpp.Language,
	csharp.Language,
	css.Language,
	gorule.Language,
	html.Language,
	java.Language,
	json.Language,
	markdown.Language,
	php.Language,
	rust.Language,
	sql.Language,
	yaml.Language,
}

func RegisterBuiltins() {
	for _, lang := range builtinLanguages {
		if err := core.RegisterLanguage(lang); err != nil {
			panic(fmt.Sprintf("language: register built-in %q: %v", lang.ID(), err))
		}
	}
}
