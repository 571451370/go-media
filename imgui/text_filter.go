package imgui

// Helper: Parse and apply text filters. In format "aaaaa[,bbbb][,ccccc]"
type TextFilter struct {
	Filters []string
}

func (t *TextFilter) IsActive() bool {
	return len(t.Filters) > 0
}