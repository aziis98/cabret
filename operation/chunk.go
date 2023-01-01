package operation

import "github.com/aziis98/cabret"

func init() {
	registerType("chunk", &Chunk{})
}

// Chunk is a list operation that will group incoming items in groups of the given size
type Chunk struct {
	Count         int
	SkipRemaining bool
}

func (op *Chunk) Configure(options map[string]any) error {
	var err error
	op.Count, err = getKey[int](options, "size")
	if err != nil {
		return err
	}
	op.SkipRemaining, err = getKey(options, "skip_remaining", false)
	if err != nil {
		return err
	}

	return nil
}

func (op *Chunk) ProcessList(items []cabret.Content) ([]cabret.Content, error) {
	totalPages := len(items) / op.Count

	chunks := make([][]cabret.Content, totalPages, totalPages+1)

	for i := 0; i < totalPages; i++ {
		chunks = append(chunks, items[i*op.Count:(i+1)*op.Count])
	}

	if !op.SkipRemaining {
		chunks = append(chunks, items[totalPages*op.Count:])
	}

	result := make([]cabret.Content, len(chunks))
	for i, chunk := range chunks {
		result[i] = cabret.Content{
			Type: cabret.MetadataOnly,
			Metadata: cabret.Map{
				"Page":       i + 1,
				"TotalPages": totalPages,
				"Items":      chunk,
			},
		}
	}

	return result, nil
}
