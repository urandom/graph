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

func NewWalker(start graph.Linker) Walker {
	v := NewVisitor()
	roots, count := findRoots(start, v)

	w := Walker{start: start, roots: roots,
		count: count, wgm: make(map[graph.Id]*sync.WaitGroup)}

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

func linkWalker(l graph.Linker,
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
			if wg, ok := wgm[t.Node().Id()]; ok {
				wg.Add(1)
			} else {
				wg := new(sync.WaitGroup)
				wg.Add(1)
				wgm[t.Node().Id()] = wg
			}

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
			if !v.Add(t.Node()) {
				continue
			}
			count++

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
			if !v.Add(t.Node()) {
				continue
			}

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
