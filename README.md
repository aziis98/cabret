# Cabret

A yaml based static site generator, ideally with the same features as Hugo but with a simpler model. 

```yaml
entryPoints:
  - source: index.html
    pipeline:
      - layout: layouts/base.html
      - target: dist/index.html
  - source: posts/{id}.md
    pipeline:
      - plugin: markdown
      - layout: layouts/base.html
      - target: dist/posts/{id}/index.html
```

## ToDo

### Tags

A case of fan-in (get all posts and group by tags) and fan-out (generate all tag pages with back-links to posts)

```yaml
build:
  # Render homepage
  - pipeline:
      - source: index.html
      - use: layout
        path: layout/base.html
      - target: dist/index.html
  # Render each post
  - pipeline:
      - source: posts/{id}.md
      - use: markdown
      - use: layout
        path: layouts/base.html
      - target: dist/posts/{id}/index.html
  # Render "posts" page
  - pipeline:
      - source: posts/{id}.md
      - use: frontmatter
      - use: sort
        key: publish_date
        direction: descending
      - use: slice
        from: 0
        to: 10
      - use: template
        path: layouts/list.html
      - use: layout
        path: layouts/base.html
      - target: dist/posts/index.html
  # Render next pages for "posts" page
  - pipeline:
      - source: posts/{id}.md
      - use: frontmatter
      - use: sort
        key: publish_date
        direction: descending
      - use: slice
        from: 10
        to: end
      - use: chunk # paginate items
        size: 10
        pipeline: # this pipeline gets called with every [size * n, size * n + size] range of items
          - use: template # aggregate this items chunk in a single item
            path: layouts/list.html
          - use: layout
            path: layouts/base.html
          - target: dist/posts/{.Chunk.Index}/index.html
  # Render "/tags/{tag}/" pages
  - pipeline:
      - source: posts/{id}.md
      - use: frontmatter
      - use: categorize
        key: tags # each post contains a  metadata field called "tags"  
        pipeline: # this pipeline gets called with all posts relative to one category
          - use: template
            path: layouts/tag.html
          - use: layout
            path: layouts/base.html
          - target: dist/tags/{category}/index.html
```

### Pagination

A case of fan-out with (various data leakages)

```yaml
entryPoints:
  ...
  - pipeline:
      - plugin: paginate
        items:
          pipeline:
            - source: posts/{id}.md
            - plugin: frontmatter
        pageSize: 10
        metadataKey: page
        pipeline:
          - layout: layouts/list.html

```

### Custom DSL

```
12
"text"
#symbol
#(None)
#(Some 123)
true
false
[1 2 3]
{ a = 1, b = 2 }
{ 
  a = 1
  b = 2
}
fn 1 2 3

# Example
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