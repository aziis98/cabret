# Cabret

A yaml based static site generator, ideally with the same features as Hugo but with a simpler model. Here is a basic example of a _Cabretfile.yaml_

```yaml
build:
  - pipeline:
      - source: index.html
      - use: layout
        path: layouts/base.html
      - target: dist/index.html
  - pipeline:
      - source: posts/{{ .Id }}.md
      - plugin: markdown
      - use: layout
        path: layouts/base.html
      - target: dist/posts/{{ .Id }}/index.html
```

Cabret is based on the idea of a data pipeline. Along a pipeline flows **a list** of `Content` objects that are just

```go
type Content struct {
    Type      string
    Metadata  map[string]any
    Data      []byte          
}
```

The `Type` is the _mime type_ of the `Data`. The `Metadata` can be whatever the user wants and operations and populate it with special fields as they please. In some cases its useful to have `Data` be `nil`, when that happens the `Type` should be `cabret.MetadataOnly` (just the string `metadata-only`)

## Operations

Each pipeline is a list of operations, the first field in an operation must be one of `source`, `use` or `target`.

- `source: <pattern>`

    Load files matching the given pattern and populates the field `.Metadata.MatchResult` with the captures.

    Other supported forms are

    ```yaml
    use: source
    path: <path pattern>
    ```

    ```yaml
    use: source
    paths: 
      - <path pattern1>
      ...
      - <path patternN>
    ```

- `target: <template>`

    For each incoming item this will render the given path template with the item metadata and write the content to disk.

- `use: <operation>`

    This will apply the provided operations to the incoming items, for now available operations are

    - **Layout** is available in the following forms

        ```yaml
        use: layout
        path: <glob pattern>
        ```

        ```yaml
        use: layout
        paths: 
          - <glob pattern 1>
          ...
          - <glob pattern N>
        ```

        This operation will load the given pattern with `.ParseFiles(...)` (this automatically chooses between the Go html or text template engines based on the first file name extension)

        The template context is the incoming item metadata as well its data inside the `.Content` variable.

    - **Markdown** doesn't need any other options (TODO: add some options to configure _goldmark_)

        ```yaml
        use: markdown
        ```

        This will render each incoming item using goldmark with some reasonable default options.

        By default this will also read YAML frontmatter.

    - **Frontmatter** as the previous one but just reads frontmatter.

        ```yaml
        use: frontmatter
        ```

    - **Categorize** can be used in the following forms

        ```yaml
        use: categorize
        key: <key>
        bind: <variable> # optional
        ```

        This will categorize all items based on the provided key, for example if key is `tags` then the incoming list of posts will be converted in a list of empty contents with a `Category` field telling the item's category and an `Items` field containing a list of posts with that tag.

        By default the output item category is placed into `Category` but this can be changed with the `bind` option.

    - (_TODO_) **Chunk** can be used to paginate some items

        ```yaml
        use: chunk
        size: <number>
        ```

        This operation will group the incoming items in chunk with the provided size (except for the last one that can end up holding less items). The output items hold no data but have the following structure

        ```go
        Metadata: {
          Page: <Page Number>
          TotalPages: <Total Page Count>
          Items: <Chunk>
        }
        ```

    - (_TODO_) **Slice** extract a sub range of items

        ```yaml
        use: slice
        from: <start> # inclusive, defaults to 0
        to: <end>     # exclusive, defaults to len(Items)
        ```

        Does what you expect to the incoming list of items.

    - (_TODO_) **Sort** sorts the incoming items using the provided key and direction (by default is ascending)

        ```yaml
        use: sort
        direction: <ascending or descending>
        key: <key to use for sorting>
        ```

