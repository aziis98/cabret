package operation

import (
	"log"

	"github.com/aziis98/cabret"
)

func init() {
	registerType("slice", &Slice{})
}

// Slice is a [cabret.ListOperation] will return a slice of the incoming items.
// Both indices can be negative and will wrap around the items list.
// For example "&Slice{From: -1, To: 0}" will reverse the incoming items list.
type Slice struct {
	From, To int
}

// Configure will configure this operation
//
//	from: <start index> # inclusive, optional and defaults to the start
//	to: <end index>     # exclusive, optional and defaults to the end
func (op *Slice) Configure(options map[string]any) error {
	var err error

	op.From, err = getKey(options, "from", 0)
	if err != nil {
		return err
	}
	op.To, err = getKey(options, "to", -1)
	if err != nil {
		return err
	}

	return nil
}

func (op *Slice) ProcessList(items []cabret.Content) ([]cabret.Content, error) {
	from := op.From
	to := op.To

	if to < 0 {
		to = to + len(items) + 1
	}
	if from < 0 {
		from = from + len(items) + 1
		from, to = to, from
		reverse(items)

		log.Printf(`[operation.Slice] slicing and reversing items from %d to %d`, from, to)
	} else {
		log.Printf(`[operation.Slice] slicing items from %d to %d`, from, to)
	}

	return items[from:to], nil
}

//
// utilities
//

// reverse taken from this answer https://stackoverflow.com/questions/28058278/how-do-i-reverse-a-slice-in-go
func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
