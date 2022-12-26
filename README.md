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
entryPoints:
  ...
  - source: posts/{id}.md
    pipeline:
      - plugin: frontmatter
      - plugin: group
        metadataKey: tag
        key: tags
        pipeline:
          - layout: layouts/tag.html
          - layout: layouts/base.html
          - target: dist/tags/{tag}/index.html # ...{tag}... is the same as "metadataKey" (?)
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