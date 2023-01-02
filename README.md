# Cabret

Cabret is a yaml based static site generator, ideally with the same features as Hugo but with a simpler model. Here is a simple example of a _Cabretfile.yaml_

## Build

```bash shell
$ go build -v -o ./bin/cabret ./cmd/cabret
```

## Introduction

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

Cabret is based on the idea of a data pipeline. Along a pipeline flows **a list** of `Content` items that are just items with a data, mime type and metadata.

```go
type Content struct {
    Type      string
    Metadata  map[string]any
    Data      []byte          
}
```

The `Type` is the _mime type_ of the `Data`. The `Metadata` can be whatever the user wants and operations can populate it with special fields as they please. In some cases its useful to have `Data` be `nil`, when that happens the `Type` should be `cabret.MetadataOnly` (just the string `"metadata-only"`)

## Operations

Each pipeline is a list of operations, the first field in an operation should be one of `source`, `use` or `target`.

- `source: <path pattern>`

    A **path pattern** is a "glob pattern" with captures, some examples are

    ```
    posts/{id}.md 
    {{folders}}/{year}-{month}-{day}.json
    ``` 

    the only difference is that `{{...}}` captures can match slashes so can be used to match more than one "path part", while `{...}` can only match non slash characters.

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

- `target: <path template>`

    A **path template** can contains `{...}` interpolations, the key inside the braces will be replaced using the `.Metadata.MatchResult` map. Double braces can also be used and are evaluated using Go `text/template` with context the whole item metadata.

    For each incoming item this will render the given path template with the item metadata and write the content to disk.

- `use: <operation>`

    This will apply the provided operations to the incoming items, for now available operations are

    - **Template** is a list operation that will concatenate all incoming items and render them using the specified template engine

        ```yaml
        use: template
        engine: <name>
        ```

        for now `engine` can be one of

        - `html` &mdash; golang html template engine

        - `text` &mdash; golang text template engine

        - (_TODO_) `handlebars` &mdash; the famous Handlebars template engine 

    - **Layout** is an item operation that works like the previous one as it will use a template engine to render the content of the output item. The main difference is that this will pass the previous item content in the template context in the `.Content` field.

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

        This operation will automatically choose between the go html or text template engines based on the first file name extension, _TODO: make this configurable like the `template` operation_)

        The template context is the incoming item metadata as well its data passed in the `.Content` variable.

        This is useful for wrapping partial html pages into full documents.

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

    - **Chunk** can be used to paginate some items

        ```yaml
        use: chunk
        size: <number>
        skip_remaining: <true or false> # optional, defaults to "false"
        ```

        This operation will group the incoming items in chunk with the provided size (except for the last one that can end up holding less items). 
        
        The option `skip_remaining` will exclude the last chunk if it doesn't have exactly "size" items.

        The output items hold no data but have the following structure

        ```go
        Metadata: {
            Index: /* 0-indexed current chunk number */,
            Total: /* total number of chunks */,
            Items: /* list of items of this chunk */,
        }
        ```

    - **Slice** extract a sub range of items

        ```yaml
        use: slice
        from: <start> # inclusive, defaults to 0
        to: <end>     # exclusive, defaults to len(Items)
        ```

        Does what you expect to the incoming list of items.

    - **Program** is an item operation that will run a specified program passing the current item as stdin. The program output is then processed and becomes the new item data.

        ```yaml
        use: program
        io: <stdin/stdout format> # optional, defaults to "raw"
        command: <shell command>  # required
        ```

        The io format specifies how the current item is passed to the given program

        - `raw` &mdash; will just pass the item `.Data` raw as stdin to the specified command.
        
        - `json` &mdash; will encode the whole item (also including type and metadata) as json ad pass that string as stdin to the specified command.

        _TODO: Make the shell command be a go text/template string, for passing metadata directly thorough the command string._

    - (_TODO_) **Sort** sorts the incoming items using the provided key and direction (by default is ascending)

        ```yaml
        use: sort
        direction: <ascending or descending>
        key: <key to use for sorting>
        ```

## Architecture

- `cabret.go` &mdash; all abstract types and functions 

- `path/` &mdash; handles path patterns and templates

- `config/` &mdash; contains config structure definitions and handles loading from YAML

- `cmd/cabret/` &mdash; module that puts it all together in a CLI application

- `operation/` &mdash; module containing all operations

    - `operation/template/` &mdash; a small module that abstract the concept of a template, used by the layout operation.

- `runner/` &mdash; a module that depends on `config`, `parse` and `operations` and evaluates operations.

- `parse/` &mdash; handles the conversion of pipelines and lists of operations to core types.

