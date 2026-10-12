package benchmark

import (
	"testing"

	codematecore "code-mate-core"
	"code-mate-core/language"
	js "code-mate-core/language/javascript"
	py "code-mate-core/language/python"
	"code-mate-core/theme"
)

func init() {
	language.RegisterBuiltins()
	theme.RegisterBuiltInThemes()
}

var jsBenchCode = `import { useState, useEffect } from 'react';

interface Props {
  name: string;
  age?: number;
}

function Greeting({ name, age }: Props): JSX.Element {
  const [count, setCount] = useState(0);

  useEffect(() => {
    document.title = "Hello, " + name;
  }, [name]);

  return (
    <div className="greeting">
      <h1>Hello, {name}!</h1>
      {age && <p>Age: {age}</p>}
      <p>Count: {count}</p>
      <button onClick={() => setCount(count + 1)}>+</button>
    </div>
  );
}

export default Greeting;
`

var pyBenchCode = `from typing import Optional, List
from dataclasses import dataclass

@dataclass
class User:
    name: str
    email: str
    age: Optional[int] = None

def get_users() -> List[User]:
    users = []
    for i in range(10):
        user = User(
            name=f"User_{i}",
            email=f"user{i}@example.com",
            age=20 + i,
        )
        users.append(user)
    return users

async def process_users():
    result = []
    users = get_users()
    for user in users:
        if user.age and user.age > 25:
            result.append(user)
    return result
`

// --- Tokenizer benchmarks (core.Parse) ---

func BenchmarkParse_JavaScript(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows := js.Language.Parse(jsBenchCode)
		_ = len(rows)
	}
}

func BenchmarkParse_Python(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows := py.Language.Parse(pyBenchCode)
		_ = len(rows)
	}
}

// --- Full pipeline benchmarks (CodeToHtml) ---

func BenchmarkCodeToHtml_SmallJS(b *testing.B) {
	hl := codematecore.NewHighlighter(codematecore.HighlighterOptions{Theme: "dark-plus"})
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		html, err := hl.CodeToHtml("const x = 1;", "javascript")
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		_ = html
	}
}

func BenchmarkCodeToHtml_Python(b *testing.B) {
	hl := codematecore.NewHighlighter(codematecore.HighlighterOptions{Theme: "dark-plus"})
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		html, err := hl.CodeToHtml(pyBenchCode, "python")
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		_ = html
	}
}

func BenchmarkCodeToHtml_LargeJS(b *testing.B) {
	hl := codematecore.NewHighlighter(codematecore.HighlighterOptions{Theme: "dark-plus"})
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		html, err := hl.CodeToHtml(jsBenchCode, "javascript")
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		_ = html
	}
}

// --- Cache hit benchmark ---

func BenchmarkCodeToHtml_CacheHit(b *testing.B) {
	hl := codematecore.NewHighlighter(codematecore.HighlighterOptions{Theme: "dark-plus"})
	// Warm up cache
	_, err := hl.CodeToHtml(jsBenchCode, "javascript")
	if err != nil {
		b.Fatalf("warmup error: %v", err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		html, err := hl.CodeToHtml(jsBenchCode, "javascript")
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		_ = html
	}
}

// --- Edge case benchmarks ---

func BenchmarkParse_EmptyCode(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows := js.Language.Parse("")
		_ = len(rows)
	}
}

func BenchmarkParse_SingleChar(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows := js.Language.Parse("x")
		_ = len(rows)
	}
}

func BenchmarkParse_ManyShortLines(b *testing.B) {
	code := "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\n"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows := js.Language.Parse(code)
		_ = len(rows)
	}
}

// --- Theme switching benchmark ---

func BenchmarkUpdateTheme(b *testing.B) {
	hl := codematecore.NewHighlighter(codematecore.HighlighterOptions{})
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hl.UpdateTheme("github-light")
		hl.UpdateTheme("dark-plus")
	}
}
