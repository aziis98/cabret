package layout

type Template interface {
	Render(ctx map[string]any) ([]byte, error)
}
