package pdsample

type WDPDFNode struct {
	mark                bool
	Parent, Left, Right *WDPDFNode
	Key                 int
	Weight, SumWeights  float64
}

func (n *WDPDFNode) SetSum() {
	n.SumWeights = n.Weight
	if n.Left != nil {
		n.SumWeights += n.Left.SumWeights
	}
	if n.Right != nil {
		n.SumWeights += n.Right.SumWeights
	}
}

type WDPDF struct {
	root *WDPDFNode
}

func (wd *WDPDF) Insert(item int, weight float64) {
}

func (wd *WDPDF) Update(item int, newWeight float64) {
}

func (wd *WDPDF) Choose(p float64) int {
	if p < 0 || p >= 1 {
		panic("argument(p) outside of valid range")
	}
	if wd.root == nil {
		panic("choose() called on empty tree")
	}

	w := wd.root.SumWeights * p
	n := wd.root

	for {
		if n.Left != nil {
			if w < n.Left.SumWeights {
				n = n.Left
				continue
			} else {
				w -= n.Left.SumWeights
			}
		}

		// shouldn't be necessary, sanity check
		if w < n.Weight || n.Right == nil {
			break
		}
		w -= n.Weight
		n = n.Right
	}

	return 0
}

func (wd *WDPDF) propagateSumsUp(n *WDPDFNode) {
	for n != nil {
		n.SetSum()
		n = n.Parent
	}
}
