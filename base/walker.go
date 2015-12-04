package base

import "github.com/urandom/graph"

type Walker struct {
	start graph.Linker
	roots []graph.Linker
	count int
}

func NewWalker(start graph.Linker) Walker {
	v := NewVisitor()
	roots, count := findRoots(start, v)
	w := Walker{start: start, roots: roots, count: count}

	return w
}

func (w Walker) Walk() <-chan graph.WalkData {
	nodes := make(chan graph.WalkData)

	return nodes
}

func (w Walker) Total() int {
	return w.count
}

func findRoots(l graph.Linker, v Visitor) (roots []graph.Linker, count int) {
	roots = append(roots, l)
	count++

	v.Add(l.Node())

	rr, rc := findBacktrackable(l, v)
	roots = append(roots, rr...)
	count += rc

	return
}

func findBacktrackable(l graph.Linker, v Visitor) (roots []graph.Linker, count int) {
	for _, c := range l.Connectors(graph.OutputType) {
		t, o := c.Target()

		if t != nil {
			if v.Visited(t.Node()) {
				continue
			}

			count++
			v.Add(t.Node())

			if len(t.Connectors()) > 1 {
				rr, rc := findRootsBacktrack(t, v, o)

				roots = append(roots, rr...)
				count += rc
			}

			rr, rc := findBacktrackable(t, v)
			roots = append(roots, rr...)
			count += rc
		}
	}

	return
}

func findRootsBacktrack(l graph.Linker, v Visitor, ignore ...graph.Connector) (roots []graph.Linker, count int) {
	for _, c := range l.Connectors() {
		if len(ignore) > 0 && c == ignore[0] {
			continue
		}

		t, _ := c.Target()
		if t != nil {
			if v.Visited(t.Node()) {
				continue
			}

			v.Add(t.Node())
			count++

			rr, rc := findRootsBacktrack(t, v)
			roots = append(roots, rr...)
			count += rc
		}
	}

	if count == 0 {
		roots = append(roots, l)
	}

	return
}
