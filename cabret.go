package cabret

import (
	"fmt"
	"log"
)

func init() {
	log.SetFlags(0)
}

const MatchResultKey = "MatchResult"

type MatchResult struct {
	File     string
	Captures map[string]string
}

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
	Load(config map[string]any) error
}

type ListOperation interface {
	Operation
	ProcessList(contents []Content) ([]Content, error)
}

type ItemOperation interface {
	Operation
	ProcessItem(content Content) (*Content, error)
}

func ProcessOperation(op Operation, inputs []Content) ([]Content, error) {
	switch op := op.(type) {
	case ListOperation:
		return op.ProcessList(inputs)

	case ItemOperation:
		outputs := []Content{}
		for _, item := range inputs {
			result, err := op.ProcessItem(item)
			if err != nil {
				return nil, err
			}

			// skip terminal operations
			if result == nil {
				continue
			}
			outputs = append(outputs, *result)
		}
		return outputs, nil
	default:
		return nil, fmt.Errorf(`invalid operation type %T`, op)
	}
}
