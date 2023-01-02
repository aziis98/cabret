package operation

import (
	"log"

	"github.com/aziis98/cabret"
)

func init() {
	registerType("chunk", &Chunk{})
}

// Chunk is a list operation that will group incoming items in groups of the given size
type Chunk struct {
	Size          int
	SkipRemaining bool
}

func (op *Chunk) Configure(options map[string]any) error {
	var err error
	op.Size, err = getKey[int](options, "size")
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
	totalChunks := len(items) / op.Size

	chunks := make([][]cabret.Content, totalChunks, totalChunks+1)

	for i := 0; i < totalChunks; i++ {
		chunks = append(chunks, items[i*op.Size:(i+1)*op.Size])
	}

	if !op.SkipRemaining {
		if len(items)%op.Size != 0 {
			lastExtraChunk := items[totalChunks*op.Size:]
			chunks = append(chunks, lastExtraChunk)

			if len(lastExtraChunk) != len(items)%op.Size {
				panic("ehm last chunk and modulus should be equal")
			}

			log.Printf(`[operation.Chunk] chunked items in %d chunk(s), last with %d items`, totalChunks+1, len(lastExtraChunk))
		} else {
			log.Printf(`[operation.Chunk] chunked items in %d chunk(s)`, totalChunks)
		}
	} else {
		log.Printf(`[operation.Chunk] chunked items in %d chunk(s)`, totalChunks)
	}

	result := make([]cabret.Content, len(chunks))
	for i, chunk := range chunks {
		result[i] = cabret.Content{
			Type: cabret.MetadataOnly,
			Metadata: cabret.Map{
				"Index": i,
				"Total": totalChunks,
				"Items": chunk,
			},
		}
	}

	return result, nil
}
