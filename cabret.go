package cabret

const MatchResult = "MatchResult"

type Map map[string]any

type File struct {
	Path string
	Content
}

type Content struct {
	// Type for known content formats is just the mime-type
	Type string

	// Data is the content of the file
	Data []byte

	// Metadata is any extra data of the file (e.g. yaml frontmatter) or injected by plugins
	Metadata Map
}

type ListOperation interface {
	MapAll(contents []Content) ([]Content, error)
}

type ItemOperation interface {
	FlatMap(content Content) (*Content, error)
}

type FlatMapToMapAll struct{ FlatMapOperation }

func (op FlatMapToMapAll) MapAll(contents []Content) ([]Content, error) {
	mapped := []Content{}

	for _, item := range contents {
		result, err := op.FlatMap(item)
		if err != nil {
			return nil, err
		}

		// skip terminal operations
		if result == nil {
			continue
		}

		mapped = append(mapped, *result)
	}

	return mapped, nil
}
