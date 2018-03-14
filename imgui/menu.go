package imgui

func (m *MenuColumns) Update(count int, spacing float64, clear bool) {
	m.Count = count
	m.Width = 0
	m.NextWidth = 0
	m.Spacing = spacing
	if clear {
		for i := range m.NextWidths {
			m.NextWidths[i] = 0
		}
	}
	for i := 0; i < m.Count; i++ {
		if i > 0 && m.NextWidths[i] > 0 {
			m.Width += m.Spacing
		}
		m.Pos[i] = float64(int(m.Width))
		m.Width += m.NextWidths[i]
		m.NextWidths[i] = 0
	}
}