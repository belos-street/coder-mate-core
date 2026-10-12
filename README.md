# CodeMate Core

A fast, zero-dependency syntax highlighting engine for Go — supporting 17 programming languages with a FSM-based tokenizer, theming system, and WASM browser integration.

```go
hl := codematecore.NewHighlighter(codematecore.HighlighterOptions{})
html, _ := hl.CodeToHtml(`const x: number = 42;`, "typescript")
// → <pre style="..."><code>...</code></pre>
```

## Features

- **17 languages** — JavaScript, TypeScript, Python, Bash, C, C++, C#, CSS, Go, HTML, Java, JSON, Markdown, PHP, Rust, SQL, YAML
- **10 themes** — Dark+ (default), GitHub Light, Dracula, One Dark Pro, Nord, Monokai, Material Ocean, Tokyo Night, Solarized Dark, Solarized Light
- **No external dependencies** — pure Go standard library
- **WASM support** — compile to WebAssembly for browser use
- **Token caching** — avoids re-parsing repeated input
- **Customizable** — per-call pre style, line class prefix, theme switching

## Usage

### Go API

```go
import codematecore "code-mate-core"

hl := codematecore.NewHighlighter(codematecore.HighlighterOptions{})

// Get token stream
tokens, err := hl.CodeToTokens("def hello(): pass", "python")

// Render HTML
html, err := hl.CodeToHtml("def hello(): pass", "python")

// With options
html, err := hl.CodeToHtml("const x = 1;", "javascript",
    codematecore.WithPreStyle("background: #1E1E1E; padding: 12px;"),
    codematecore.WithLineClassPrefix("ln-"),
)

// Switch theme
if err := hl.UpdateTheme("github-light"); err != nil {
    // handle unknown theme
}
```

### Custom Highlighter

```go
hl := codematecore.NewHighlighter(codematecore.HighlighterOptions{
    Theme:           "github-light",
    PreStyle:        "background: #fff; padding: 16px;",
    LineClassPrefix: "line-",
})
```

> Unknown theme names or invalid theme values fall back to the default theme (`dark-plus`).

### WASM (Browser)

```go
// Build
GOOS=js GOARCH=wasm go build -o coder-mate.wasm ./cmd/wasm/
```

```html
<script src="wasm_exec.js"></script>
<script>
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("coder-mate.wasm"), go.importObject)
    .then(result => {
      go.run(result.instance);

      // Quick highlight
      const r = coderMateHighlight("const x = 1;", "javascript");
      document.body.innerHTML = r.html;

      // Multi-instance API
      const id = coderMate.createHighlighter({ theme: "github-light" });
      const r2 = coderMate.highlight(id, "print('hi')", "python");
      coderMate.updateTheme(id, "dark-plus");
      coderMate.dispose(id);
    });
</script>
```

### WASM JS API Reference

| Function | Description |
|----------|-------------|
| `coderMateHighlight(code, lang)` | Highlight with default instance, returns `{ html }` |
| `coderMate.createHighlighter(options)` | Create instance, returns numeric ID |
| `coderMate.highlight(id, code, lang)` | Highlight with specific instance |
| `coderMate.updateTheme(id, theme)` | Switch theme on an instance |
| `coderMate.dispose(id)` | Release instance memory |
| `coderMate.getThemes()` | List available themes |
| `coderMate.getLanguages()` | List available languages |

Options for `createHighlighter`: `theme` (string), `preStyle` (string), `lineClassPrefix` (string).

## Supported Languages

| ID | Aliases |
|----|---------|
| `javascript` | `js` |
| `typescript` | `ts` |
| `python` | `py` |
| `bash` | `sh`, `shell` |
| `c` | `h`, `c89`, `c99`, `c11` |
| `cpp` | `c++`, `cc`, `cxx`, `hpp`, `hxx`, `hh` |
| `csharp` | `cs`, `c#` |
| `css` | |
| `go` | `golang` |
| `html` | `htm` |
| `java` | `jav` |
| `json` | |
| `markdown` | `md` |
| `php` | `phtml`, `php8` |
| `rust` | `rs` |
| `sql` | `postgresql`, `pgsql` |
| `yaml` | `yml` |

## Themes

| ID | Display Name |
|----|-------------|
| `dark-plus` | Dark+ (default) |
| `github-light` | GitHub Light |
| `dracula` | Dracula |
| `one-dark-pro` | One Dark Pro |
| `nord` | Nord |
| `monokai` | Monokai |
| `material-ocean` | Material Ocean |
| `tokyo-night` | Tokyo Night |
| `solarized-dark` | Solarized Dark |
| `solarized-light` | Solarized Light |

## Architecture

```
code  →  Tokenizer (FSM)  →  Token Stream  →  Theme Engine  →  HTML
         ┌─────────────────┐  ┌──────────────┐  ┌──────────────┐
         │ 17 language      │  │ token[].Scope  │  │ scope → color│
         │ specs each with  │  │ token[].Line   │  │ inline styles│
         │ 10-30 grammar    │  │ token[].Col    │  │              │
         │ rules + states   │  │ token[].Style  │  │              │
         └─────────────────┘  └──────────────┘  └──────────────┘
```

The tokenizer uses a **finite state machine** with push/pop state stack. Each language defines grammar rules per state, and rules are matched in priority order. Cross-line tokens (e.g. multi-line strings) are handled by persisting state across line breaks.

## Performance

Benchmarked on **Intel Core Ultra 7 258V** (Windows, Go 1.23):

```
BenchmarkParse_JavaScript-8           2,595    457,296 ns/op   101,666 B/op   1,495 allocs/op
BenchmarkParse_Python-8               3,022    344,932 ns/op    75,905 B/op   1,094 allocs/op
BenchmarkCodeToHtml_SmallJS-8      329,286      3,877 ns/op     4,987 B/op      50 allocs/op
BenchmarkCodeToHtml_Python-8         9,789    126,112 ns/op   138,525 B/op   1,106 allocs/op
BenchmarkCodeToHtml_LargeJS-8        7,552    165,096 ns/op   182,321 B/op   1,419 allocs/op
BenchmarkCodeToHtml_CacheHit-8       6,897    165,755 ns/op   182,291 B/op   1,419 allocs/op
BenchmarkParse_EmptyCode-8      43,432,467         24 ns/op        16 B/op       1 allocs/op
BenchmarkParse_SingleChar-8      1,372,580        778 ns/op       301 B/op       8 allocs/op
BenchmarkParse_ManyShortLines-8     75,010     15,733 ns/op     5,323 B/op      96 allocs/op
BenchmarkUpdateTheme-8          14,727,775         76 ns/op         0 B/op       0 allocs/op
```

- **~0.5 ms** to parse and highlight a 30-line JavaScript file
- **~4 μs** for a small `const x = 1;` (heavily cached)
- **~24 ns** for empty input (near-zero overhead)
- **Cache hit** matches the parse cost (cache avoids re-parsing but still clones token stream)
- **Theme switching** is near-free (76 ns, zero allocations)

### Go WASM vs JS Version

Compared to the [JS version](https://github.com/belos-street/coder-mate-core-js) of the same engine (same CPU, same algorithm):

| Metric | JS (Bun) | Go WASM | Speedup |
|--------|----------|---------|---------|
| JavaScript ~100 lines (full pipeline) | 7.79 ms | ~0.55 ms | **~14x** |
| Python ~120 lines (full pipeline) | 8.36 ms | ~0.84 ms | **~10x** |
| Single small expression | ~0.05 ms | ~0.004 ms | **~12x** |

The Go WASM build runs **10-15x faster** than the equivalent JS bundle for browser-side syntax highlighting, while producing identical token output.

Run benchmarks yourself:

```bash
# Go benchmarks
go test -bench=. -benchmem -count=1 ./benchmark/

# JS benchmarks
git clone https://github.com/belos-street/coder-mate-core-js
cd coder-mate-core-js && bun install
bun run bench:js-language
```

## Related Projects

- [coder-mate-core-js](https://github.com/belos-street/coder-mate-core-js) — The original TypeScript version with the same FSM-based tokenizer, supporting 17 languages and 10 themes. Serves as the reference implementation and is suitable for Node.js / Bun environments.

## Requirements

- Go 1.23+
- No external dependencies

## License

MIT
