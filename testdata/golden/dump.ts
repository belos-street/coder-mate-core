/**
 * Golden 基准生成脚本：用 TypeScript 参考实现（coder-mate-core-js）对
 * 固定代码片段生成 token 流，输出 JSON 供 Go 侧 golden 一致性测试比对。
 *
 * 用法：bun run testdata/golden/dump.ts <ts-repo-root> <output-dir>
 *   <ts-repo-root>  coder-mate-core-js 仓库本地克隆路径
 *   <output-dir>    JSON 输出目录（默认与本脚本同目录）
 *
 * 口径：Token 只取 text / scope / line —— Col 在 TS（UTF-16 code unit）
 * 与 Go（UTF-8 byte）间语义不同，属文档化的有意分歧，不参与比对。
 */

import { mkdirSync, writeFileSync } from 'node:fs'
import { pathToFileURL } from 'node:url'
import { resolve, dirname, join } from 'node:path'

const tsRoot = process.argv[2]
const outDir = process.argv[3] ?? dirname(new URL(import.meta.url).pathname)

if (!tsRoot) {
  console.error('usage: bun run dump.ts <ts-repo-root> [output-dir]')
  process.exit(1)
}

const managerPath = pathToFileURL(resolve(tsRoot, 'lib/language/manager.ts')).href
const { tokenize } = await import(managerPath)

interface Case {
  lang: string
  name: string
  code: string
}

const CASES: Case[] = [
  // ---- javascript ----
  { lang: 'javascript', name: 'basics', code: 'const x = 42; // answer\nfunction greet(name) { return `hi ${name}`; }\n' },
  { lang: 'javascript', name: 'strings-comments', code: '/* block\n   comment */\nlet s = "a\\"b" + \'c\\n\';\nif (a === b) x?.y ?? z;\n' },
  { lang: 'javascript', name: 'class-numeric', code: 'class A extends B { constructor() { super(); this.v = new Map(); } }\nexport default 1e+10;\n' },
  // ---- typescript ----
  { lang: 'typescript', name: 'types', code: 'type A = { x: number };\ninterface B extends A { y?: string }\nconst v: Array<Promise<number>> = [];\n' },
  { lang: 'typescript', name: 'class-generics', code: 'class C<T> implements I { public readonly f: Map<string, unknown> = new Map(); }\nfunction g<T>(p: T) { return true as unknown; }\n' },
  { lang: 'typescript', name: 'declare-enum', code: 'declare enum E { A, B }\nlet u: string | null = null;\nconst w = v satisfies Partial<A>;\n' },
  // ---- python ----
  { lang: 'python', name: 'function-annotation', code: 'def greet(name: str) -> str:\n    return f"hi {name}"\n' },
  { lang: 'python', name: 'class-keywords', code: 'class A(Exception):\n    x: int = 0x1F\n    def m(self, *args, **kw) -> None:\n        assert True, \'msg\'\n' },
  { lang: 'python', name: 'control-flow', code: 'import os\nfrom sys import path\ntry:\n    y = [i ** 2 for i in range(10)]\nexcept ValueError as e:\n    raise\n' },
  // ---- bash ----
  { lang: 'bash', name: 'loop', code: '#!/bin/bash\nset -euo pipefail\nfor f in *.txt; do echo "$f"; done\n' },
  { lang: 'bash', name: 'conditional', code: 'if [ -z "$VAR" ]; then\n  export PATH="/usr/bin:$PATH"\nfi\n' },
  { lang: 'bash', name: 'function-case', code: 'func() { local x=1; case $x in a) echo a;; esac }\n' },
  // ---- c ----
  { lang: 'c', name: 'hello', code: '#include <stdio.h>\nint main(void) { printf("hi\\n"); return 0; }\n' },
  { lang: 'c', name: 'struct-comment', code: '/* multi\nline */\nstatic const char *s = "a\\0b";\nstruct P { int x; float y; };\n' },
  { lang: 'c', name: 'loop', code: 'for (int i = 0; i < 10; ++i) { sum += arr[i] * 2; }\n' },
  // ---- cpp ----
  { lang: 'cpp', name: 'vector-rawstring', code: '#include <vector>\nstd::vector<int> v{1, 2, 3};\nauto s = R"(raw\\string)";\n' },
  { lang: 'cpp', name: 'template-class', code: 'template <typename T>\nclass Box { public: T value; };\n' },
  { lang: 'cpp', name: 'namespace', code: 'namespace n { struct S { virtual ~S() = default; }; }\n' },
  // ---- csharp ----
  { lang: 'csharp', name: 'hello', code: 'using System;\nnamespace N { class P { static void Main() { Console.WriteLine("hi"); } } }\n' },
  { lang: 'csharp', name: 'collections', code: 'var list = new List<int> { 1, 2 };\nstring s = @"verbatim \\n";\n' },
  { lang: 'csharp', name: 'record', code: 'public record R(int A, string B) { public int C { get; init; } }\n' },
  // ---- css ----
  { lang: 'css', name: 'basic', code: 'body { color: #fff; margin: 0 auto; }\n' },
  { lang: 'css', name: 'selectors', code: '.cls, #id > a:hover { background: url(img.png) no-repeat; }\n' },
  { lang: 'css', name: 'media-query', code: '@media (max-width: 600px) { .x { --v: calc(1px + 2em); } }\n' },
  // ---- go ----
  { lang: 'go', name: 'hello', code: 'package main\n\nimport "fmt"\n\nfunc main() { fmt.Println("hi") }\n' },
  { lang: 'go', name: 'struct-map', code: 'type S struct { A int `json:"a"` }\nvar m = map[string]int{"x": 1}\n' },
  { lang: 'go', name: 'generics', code: 'func F[T any](v T) *T { s := []byte(`raw`); return nil }\n' },
  // ---- html ----
  { lang: 'html', name: 'document', code: '<!DOCTYPE html>\n<html lang="en">\n<body class="x"><p>hi</p></body>\n</html>\n' },
  { lang: 'html', name: 'comment-entities', code: '<!-- comment -->\n<div id="a" data-x=\'1\'>text &amp; more</div>\n' },
  { lang: 'html', name: 'void-elements', code: '<input type="text" value="v" disabled>\n<img src=a.png>\n' },
  // ---- java ----
  { lang: 'java', name: 'hello', code: 'package com.x;\npublic class Main { public static void main(String[] args) { System.out.println("hi"); } }\n' },
  { lang: 'java', name: 'annotation', code: '/** javadoc */\n@interface A { String v() default "x"; }\n' },
  { lang: 'java', name: 'loop-generics', code: 'for (int i = 0; i < n; i++) { list.add(new Foo(i)); }\n' },
  // ---- json ----
  { lang: 'json', name: 'mixed-values', code: '{"a": 1, "b": [true, false, null]}\n' },
  { lang: 'json', name: 'string-escapes', code: '{"s": "line\\nnext", "n": -1.5e+10}\n' },
  { lang: 'json', name: 'nested', code: '{"nested": {"deep": {"x": []}}}\n' },
  // ---- markdown ----
  { lang: 'markdown', name: 'headings-emphasis', code: '# Title\n\nSome *emph* and **strong** text.\n' },
  { lang: 'markdown', name: 'lists', code: '- item one\n- item two\n  - nested\n' },
  { lang: 'markdown', name: 'code-link', code: '```go\ncode block\n```\n\n[link](http://x.y) and `code`\n' },
  // ---- php ----
  { lang: 'php', name: 'interpolation', code: '<?php\n$name = \'world\';\necho "hi $name";\n?>\n' },
  { lang: 'php', name: 'function', code: 'function f(int $x): ?string { return $x > 0 ? "$x" : null; }\n' },
  { lang: 'php', name: 'array-foreach', code: '$arr = [\'a\' => 1, \'b\' => 2];\nforeach ($arr as $k => $v) { echo $k; }\n' },
  // ---- rust ----
  { lang: 'rust', name: 'hello', code: 'fn main() { let v: Vec<i32> = vec![1, 2]; println!("{}", v[0]); }\n' },
  { lang: 'rust', name: 'lifetime-impl', code: 'struct S<\'a> { name: &\'a str }\nimpl<\'a> S<\'a> { fn new(n: &\'a str) -> Self { S { name: n } } }\n' },
  { lang: 'rust', name: 'rawstring-match', code: '/* block */\nlet r = r"raw\\path";\nmatch x { Some(y) => y, None => 0 }\n' },
  // ---- sql ----
  { lang: 'sql', name: 'select', code: 'SELECT id, name FROM users WHERE age > 21 ORDER BY name;\n' },
  { lang: 'sql', name: 'dml', code: 'INSERT INTO t (a, b) VALUES (1, \'x\');\nUPDATE t SET a = 2 WHERE b IS NULL;\n' },
  { lang: 'sql', name: 'ddl', code: 'CREATE TABLE t (\n  id INT PRIMARY KEY,\n  v VARCHAR(255) NOT NULL\n);\n' },
  // ---- yaml ----
  { lang: 'yaml', name: 'scalars', code: 'name: app\nversion: 1\nenabled: true\n' },
  { lang: 'yaml', name: 'sequence-nested', code: 'items:\n  - a\n  - b\nnested:\n  key: "quoted # not comment"\n' },
  { lang: 'yaml', name: 'anchor-tag', code: 'anchors: &a 1\nref: *a\ntag: !!str 42\n' },
]

mkdirSync(outDir, { recursive: true })

let written = 0
for (const c of CASES) {
  let rows
  try {
    rows = tokenize(c.code, c.lang)
  } catch (err) {
    console.error(`FAIL ${c.lang}-${c.name}: ${err}`)
    process.exit(1)
  }
  const doc = {
    lang: c.lang,
    case: c.name,
    code: c.code,
    // 每行 token 仅保留 text(t) / scope(s) / line(l)；Col 为文档化分歧不比对
    tokens: rows.map((row) => row.map((tok) => ({ t: tok.text, s: tok.scope, l: tok.line }))),
  }
  const file = join(outDir, `${c.lang}-${c.name}.json`)
  writeFileSync(file, JSON.stringify(doc, null, 2) + '\n')
  written++
}

console.log(`wrote ${written} golden files to ${outDir}`)
