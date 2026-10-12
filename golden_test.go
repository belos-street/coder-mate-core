package codematecore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"code-mate-core/core"
	"code-mate-core/language"
)

// golden 基准一致性测试：testdata/golden/*.json 由 testdata/golden/dump.ts
// 以 TypeScript 参考实现（coder-mate-core-js）生成，记录固定代码片段的
// token 流（text/scope/line）。本测试将 Go 实现的输出与基准逐 token 比对，
// 兑现 spec 08 §2 承诺的跨实现一致性验证。
//
// Col 不参与比对：TS 按 UTF-16 code unit 计列，Go 按 UTF-8 byte 计列，
// 属文档化的有意分歧（见 core/types.go）。

type goldenToken struct {
	Text  string `json:"t"`
	Scope string `json:"s"`
	Line  int    `json:"l"`
}

type goldenCase struct {
	Lang   string          `json:"lang"`
	Case   string          `json:"case"`
	Code   string          `json:"code"`
	Tokens [][]goldenToken `json:"tokens"`
}

// knownDivergences：Go 相对 TS 参考实现的「已知分歧清单」（首跑 2026-10-09 全量归因）。
// 清单内的 case 只断言分词无损性与行号一致性；清单外的任何逐 token 差异都是回归：
// 要么修 Go，要么给出新的分歧理由并更新本清单。
//
// 结构性根因说明：TS 规则大量使用 lookahead（如函数名 /^(...)\s*(?=\()/），
// Go regexp（RE2）不支持 lookahead，等价改写必须吞掉分隔符，
// 由此产生 token 边界类分歧——这不是缺陷，是引擎能力的文档化取舍。
var knownDivergences = map[string]string{
	// —— 有意修正：TS 原版存在同样的规则缺口（2.4）——
	"javascript-class-numeric":  "Go 补齐 JS 关键字 typeof/new/this/super 与数字指数符号 1e+10（TS 原版缺口）",
	"typescript-class-generics": "Go 补齐 JS 关键字 new（TS 原版缺口，经 typescript 继承）",
	"python-class-keywords":     "Go 补齐 assert/del/global/nonlocal 关键字与 0x 数字字面量（TS 原版缺口）",
	// —— 结构性分歧：TS lookahead 拆分函数名与 '('，Go 无法等价改写 ——
	"bash-function-case":   "RE2 无 lookahead：Go 函数名 token 吞 '('（TS 拆分为 name + operator）",
	"c-hello":              "RE2 无 lookahead：Go 函数名 token 吞 '('（TS 拆分为 name + operator）",
	"cpp-namespace":        "RE2 无 lookahead：Go 函数名 token 吞 '('（TS 拆分为 name + operator）",
	"csharp-hello":         "RE2 无 lookahead：Go 函数名 token 吞 '('（TS 拆分为 name + operator）",
	"csharp-record":        "RE2 无 lookahead：Go 函数名 token 吞 '('（TS 拆分为 name + operator）",
	"go-hello":             "Go go 规则缺 entity.name.function scope 且补齐后仍无法拆括号（TS 用 lookahead）",
	"go-generics":          "Go go 类型参数 scope 吞 '['（TS 用 lookahead 拆分）",
	"java-annotation":      "RE2 无 lookahead：Go 函数名 token 吞 '('（TS 拆分为 name + operator）",
	"java-hello":           "RE2 无 lookahead：Go 函数名 token 吞 '('（TS 拆分为 name + operator）",
	"java-loop-generics":   "RE2 无 lookahead：Go 函数名 token 吞 '('（TS 拆分为 name + operator）",
	"php-function":         "RE2 无 lookahead：Go 函数名 token 吞 '('（TS 拆分为 name + operator）",
	"rust-hello":           "RE2 无 lookahead：Go 函数名 token 吞 '('（TS 拆分为 name + operator）",
	"rust-lifetime-impl":   "RE2 无 lookahead：Go 函数名 token 吞 '('（TS 拆分为 name + operator）",
	"rust-rawstring-match": "RE2 无 lookahead：Go 函数名 token 吞 '('（TS 拆分为 name + operator）",
	// —— 结构性分歧：TS 属性/分隔符/命名空间用 lookahead 拆分，Go 合并 ——
	"bash-conditional":      "RE2 无 lookahead：Go bash 变量赋值 token 含 '='（TS 拆分为 var + operator）",
	"cpp-vector-rawstring":  "RE2 无 lookahead：Go cpp 命名空间 token 吞 '::'（TS 拆分为 std + ::，2026-10-12 归因）",
	"css-basic":             "RE2 无 lookahead：Go css 属性名 token 含 ':'（TS 拆分为 property + separator）",
	"css-media-query":       "RE2 无 lookahead：Go css 属性名 token 含 ':'（TS 拆分为 property + separator）",
	"css-selectors":         "RE2 无 lookahead：Go css 属性名 token 含 ':'（TS 拆分为 property + separator）",
	"yaml-anchor-tag":       "RE2 无 lookahead：Go yaml 属性名 token 含 ':'（TS 拆分为 property + separator）",
	"yaml-scalars":          "RE2 无 lookahead：Go yaml 属性名 token 含 ':'（TS 拆分为 property + separator）",
	"yaml-sequence-nested":  "RE2 无 lookahead：Go yaml 属性名 token 含 ':'（TS 拆分为 property + separator）",
	// —— 粒度差异 ——
	"markdown-headings-emphasis": "markdown 强调/斜体的 token 粒度不同（TS 整词，Go 逐字符，渲染等价）",
}

// scopeOnlyDivergences：token 数量与 text 序列与 TS 完全一致、仅个别 scope
// 有文档化分歧的 case（TS quirk）。做硬化比对：text 严格相等，scope 按
// 映射表（golden scope → Go scope）转换后严格相等——scope 级回归在此显式失败。
var scopeOnlyDivergences = map[string]map[string]string{
	"c-struct-comment":   {"entity.name.macro.c": "entity.name.type.c"},
	"cpp-template-class": {"entity.name.macro.cpp": "entity.name.type.cpp"},
}

func flattenGolden(rows [][]goldenToken) []goldenToken {
	var out []goldenToken
	for _, row := range rows {
		out = append(out, row...)
	}
	return out
}

func flattenStream(rows core.TokenStream) []goldenToken {
	var out []goldenToken
	for _, row := range rows {
		for _, tok := range row {
			out = append(out, goldenToken{Text: tok.Text, Scope: tok.Scope, Line: tok.Line})
		}
	}
	return out
}

func TestGoldenConsistencyWithTSReference(t *testing.T) {
	dir := "testdata/golden"
	entries, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil || len(entries) == 0 {
		t.Fatalf("no golden files found in %s (generate with: bun run testdata/golden/dump.ts <ts-repo> testdata/golden)", dir)
	}

	for _, file := range entries {
		file := file
		t.Run(strings.TrimSuffix(filepath.Base(file), ".json"), func(t *testing.T) {
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("read golden: %v", err)
			}

			var gc goldenCase
			if err := json.Unmarshal(data, &gc); err != nil {
				t.Fatalf("unmarshal golden: %v", err)
			}

			rows, err := language.Tokenize(gc.Code, gc.Lang)
			if err != nil {
				t.Fatalf("tokenize: %v", err)
			}

			got := flattenStream(rows)
			want := flattenGolden(gc.Tokens)

			// 无损性：所有实现都必须满足
			var b strings.Builder
			for _, tok := range got {
				b.WriteString(tok.Text)
			}
			if b.String() != gc.Code {
				t.Errorf("tokenization is lossy: got %q, want %q", b.String(), gc.Code)
			}

			// 硬化分支：仅 scope 有文档化分歧的 case——text 严格相等，
			// scope 经映射表转换后严格相等（scope 级回归在此显式失败）
			key := gc.Lang + "-" + gc.Case
			if scopeMap, ok := scopeOnlyDivergences[key]; ok {
				if len(got) != len(want) {
					t.Fatalf("scope-divergent case changed token count: got %d, golden %d（如为有意变更，请更新 scopeOnlyDivergences/knownDivergences）", len(got), len(want))
				}
				for i := range want {
					if got[i].Text != want[i].Text {
						t.Errorf("token %d text: got %q, want %q", i, got[i].Text, want[i].Text)
						continue
					}
					expectedScope := want[i].Scope
					if mapped, ok := scopeMap[want[i].Scope]; ok {
						expectedScope = mapped
					}
					if got[i].Scope != expectedScope {
						t.Errorf("token %d scope: got %q, want %q（golden %q 经映射表 %v）", i, got[i].Scope, expectedScope, want[i].Scope, scopeMap)
					}
				}
				return
			}

			reason, divergent := knownDivergences[key]
			if divergent {
				// 已知分歧 case：仅要求无损性（上方已断言）与行号单调递增，
				// 不做 token 数量/逐 token 比对——lookahead 类分歧本身就改变 token 数
				prev := 0
				for i, tok := range got {
					if tok.Line < prev {
						t.Errorf("known-divergent case (%s): line regression at token %d: %d < %d", reason, i, tok.Line, prev)
					}
					prev = tok.Line
				}
				return
			}

			// 严格比对：text + scope + line
			if len(got) != len(want) {
				t.Fatalf("token count mismatch: got %d, want %d\nfirst diff:\n%s", len(got), len(want), diffSnippet(got, want))
			}
			for i := range got {
				if got[i] != want[i] {
					t.Errorf("token %d mismatch: got {text:%q scope:%q line:%d}, want {text:%q scope:%q line:%d}\n%s",
						i, got[i].Text, got[i].Scope, got[i].Line, want[i].Text, want[i].Scope, want[i].Line, diffSnippet(got, want))
					return
				}
			}
		})
	}
}

// diffSnippet 输出首个差异点附近的上下文，便于定位
func diffSnippet(got, want []goldenToken) string {
	n := len(got)
	if len(want) < n {
		n = len(want)
	}
	idx := n
	for i := 0; i < n; i++ {
		if got[i] != want[i] {
			idx = i
			break
		}
	}

	lo := idx - 2
	if lo < 0 {
		lo = 0
	}
	hi := idx + 4
	var b strings.Builder
	b.WriteString("  got:  ")
	for i := lo; i < hi && i < len(got); i++ {
		b.WriteString(got[i].Text + "/" + got[i].Scope + " ")
	}
	b.WriteString("\n  want: ")
	for i := lo; i < hi && i < len(want); i++ {
		b.WriteString(want[i].Text + "/" + want[i].Scope + " ")
	}
	return b.String()
}
