# Script

Mini scripting language for querying the website "model".

## Syntax

### Values

```go
1
-1
0xFF
3.14
"some text"
'r' 'u' 'n' 'e'
```

### Property Access

```lua
-- Primitives
12
-3.14
"text"

-- Meta
#symbol
#(None)
#(Some 123)

-- Booleans
true
false

-- Lists
[1 2 3]

-- Dicts
{ a = 1, b = 2 }
{ 
  a = 1
  b = 2
}

-- Function call
fn 1 2 3

-- as long as it continues "inline" its fine and still a single expression 
fn [1 2
    3] {
        a = 1
        b = (foo #other {
            bar = 456
        })
    } [
        4 5 6
    ]

-- Arithmetic (no precedence)
-- anything matching the following regex is considered an operator
[+-*/%<>=&|^?#@]+

-- Control Flow (just builtin macros)
if cond trueCaseExpr
ifelse cond trueCaseExpr falseCaseExpr

match value [
    case valuePattern1 case1Expr
    case valuePattern2 case2Expr
    ...
    default defaultExpr
]

for var items body

-- Functions
fn [a] [
    printfln "a = %v" a
]

fn example [a b] [
    printfln "a + b = %v" (a + b)
]

example 42 69

-- Example
build [
  pipeline [
    source "index.html"
    layout "layout/base.html"
    target "dist/index.html"
  ]
  pipeline [
    source "posts/{id}.html"
    markdown
    layout "layout/base.html"
    target "dist/posts/{{ .Id }}/index.html"
  ]
  pipeline [
    source "posts/{id}.md"
    frontmatter
    sort #descending { key = "publish_date" }
    slice { to = 10 }
    template "layouts/list.html"
    layout "layouts/base.html"
    target "dist/posts/index.html"
  ]
  pipeline [
    source "posts/{id}.md"
    frontmatter
    sort #descending { key = "publish_date" }
    slice { from = 10 }
    chunk 10 {
      template "layouts/list.html"
      layout "layouts/base.html"
      target "dist/posts/{{ .Chunk.Index }}/index.html"
    }
  ]
  pipeline [
    source "posts/{id}.md"
    frontmatter
    categorize "tags" {
      template "layouts/tag.html"
      layout "layouts/base.html"
      target "dist/tags/{{ .Category }}/index.html"
    }
  ]
]
```