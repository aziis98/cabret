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

type Operation interface {
	Process(content Content) (*Content, error)
}
