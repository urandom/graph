package base

import (
	"sync"

	"github.com/urandom/graph"
)

type Walker struct {
	start graph.Linker
	roots []graph.Linker
	wgm   map[graph.Id]*sync.WaitGroup
	count int
}

type edge struct {
	from graph.Id
	to   graph.Id
}

func NewWalker(start graph.Linker) Walker {
	v := NewVisitor()
	roots, count, wgm := findRoots(start, v)

	w := Walker{start: start, roots: roots,
		count: count, wgm: wgm}

	return w
}

func (w Walker) Walk() <-chan graph.WalkData {
	nodes := make(chan graph.WalkData)
	counter := make(chan struct{})
	v := NewVisitor()

	for _, r := range w.roots {
		go linkWalker(r, []graph.Connector{}, nodes, counter, w.wgm, v)
	}

	go closeNodes(nodes, counter, w.Total())

	return nodes
}

func (w Walker) Total() int {
	return w.count
}

func linkWalker(
	l graph.Linker,
	connectors []graph.Connector,
	nodes chan<- graph.WalkData,
	counter chan<- struct{},
	wgm map[graph.Id]*sync.WaitGroup,
	v Visitor,
) {
	if !v.Add(l.Node()) {
		return
	}

	if wg, ok := wgm[l.Node().Id()]; ok {
		wg.Wait()
	}

	done := make(chan struct{})

	nodes <- graph.NewWalkData(l.Node(), connectors, done)
	counter <- struct{}{}

	for _, out := range l.Connectors(graph.OutputType) {
		if t, _ := out.Target(); t != nil {
			go linkWalker(t, t.Connectors(), nodes, counter, wgm, v)
		}
	}

	go func() {
		select {
		case <-done:
			for _, out := range l.Connectors(graph.OutputType) {
				if tl, _ := out.Target(); tl != nil {
					if wg, ok := wgm[tl.Node().Id()]; ok {
						wg.Done()
					}
				}
			}
		}
	}()
}

func closeNodes(nodes chan graph.WalkData, counter <-chan struct{}, total int) {
	for {
		select {
		case <-counter:
			total--
			if total == 0 {
				close(nodes)
			}
		}
	}
}

func findRoots(l graph.Linker, v Visitor) (roots []graph.Linker, count int, wgm map[graph.Id]*sync.WaitGroup) {
	wgm = make(map[graph.Id]*sync.WaitGroup)
	roots = append(roots, l)
	count++

	v.Add(l.Node())

	wgv := map[edge]bool{}
	rr, rc := findBacktrackable(l, v, wgm, wgv)
	roots = append(roots, rr...)
	count += rc

	return
}

func findBacktrackable(
	l graph.Linker,
	v Visitor,
	wgm map[graph.Id]*sync.WaitGroup,
	wgv map[edge]bool,
) (roots []graph.Linker, count int) {
	for _, c := range l.Connectors(graph.OutputType) {
		if t, _ := c.Target(); t != nil {
			if !v.Add(t.Node()) {
				continue
			}
			count++

			if len(t.Connectors()) > 1 {
				rr, rc := findRootsBacktrack(t, v, wgm, wgv)

				roots = append(roots, rr...)
				count += rc
			} else {
				var wg sync.WaitGroup
				hasParents := false
				for _, in := range t.Connectors() {
					if tc, _ := in.Target(); tc != nil {
						edge := edge{from: tc.Node().Id(), to: t.Node().Id()}
						if !wgv[edge] {
							wg.Add(1)
							hasParents = true
							wgv[edge] = true
						}
					}
				}

				if hasParents {
					wgm[t.Node().Id()] = &wg
				}

			}

			rr, rc := findBacktrackable(t, v, wgm, wgv)
			roots = append(roots, rr...)
			count += rc
		}
	}

	return
}

func findRootsBacktrack(
	l graph.Linker,
	v Visitor,
	wgm map[graph.Id]*sync.WaitGroup,
	wgv map[edge]bool,
) (roots []graph.Linker, count int) {
	var wg sync.WaitGroup
	var hasParents bool

	for _, in := range l.Connectors() {
		if t, _ := in.Target(); t != nil {
			edge := edge{from: t.Node().Id(), to: l.Node().Id()}
			if !wgv[edge] {
				wg.Add(1)
				hasParents = true
				wgv[edge] = true
			}

			if !v.Add(t.Node()) {
				continue
			}

			count++
			rr, rc := findRootsBacktrack(t, v, wgm, wgv)
			roots = append(roots, rr...)
			count += rc
		}
	}

	if count == 0 {
		roots = append(roots, l)
	}
	if hasParents {
		wgm[l.Node().Id()] = &wg
	}

	return
}
