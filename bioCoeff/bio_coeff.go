package biocoeff

import (
	"slices"
)

type BackTrack struct {
	p, rot  []int
	k       int
	content map[string]any

	AcceptedSubsets [][]int
	acceptedCount   *int
}

func (b *BackTrack) root(p []int) []int {
	if !b.validateOrder(p) {
		return nil
	}

	return []int{}
}

type Node struct {
	Subset   []int  `json:"subset"`
	Children []Node `json:"children"`
}

func (b *BackTrack) reject(P []int, c []int) bool {

	if len(c) > b.k {
		return true
	}

	if len(c) == 0 {
		return false
	}

	available := 0
	for _, item := range P {
		if item > c[len(c)-1] {
			available++
		}
	}

	// return by compare availabily_len(c) and b.k
	return available+len(c) < b.k
}

func (b *BackTrack) accept(P []int, c []int, k int) bool {
	if len(P) == 0 {
		return false
	}

	if len(c) == k {
		// return/store to output
		b.AcceptedSubsets = append(b.AcceptedSubsets, append([]int(nil), c...))
		// fmt.Printf("subset: %v\n", c)
		return true

	}

	return false
}

func (b *BackTrack) first(P []int, c []int) []int {
	if P == nil {
		return nil
	}

	if !b.validateOrder(P) {
		return nil
	}

	// if c doesnt have items in P then add it from smallest of it
	cMap := make(map[int]struct{})
	for _, item := range c {
		cMap[item] = struct{}{}
	}

	if len(c) == 0 {
		for _, item := range P {
			if _, exist := cMap[item]; !exist {
				// then add smallest item from it once
				first := make([]int, len(c), len(c)+1)
				copy(first, c)
				first = append(first, item)

				return first
			}
		}
	}

	for _, item := range P {
		if _, exist := cMap[item]; !exist && item > c[len(c)-1] {
			// then add smallest item from it once
			first := make([]int, len(c), len(c)+1)
			copy(first, c)
			first = append(first, item)

			return first
		}
	}

	return nil
}

func (b *BackTrack) next(P []int, s []int) []int {
	if P == nil {
		return nil
	}

	if !b.validateOrder(P) {
		return nil
	}

	prefix := s[:len(s)-1]

	for _, item := range P {
		// if item greater than s[]-1
		if item > s[len(s)-1] && !slices.Contains(prefix, item) {
			// then sibling[len(sibling)-1] = item
			sibling := make([]int, len(s))
			copy(sibling, s)
			sibling[len(sibling)-1] = item
			return sibling
		}
	}

	return nil
}

func (b *BackTrack) Subsets() [][]int { return b.AcceptedSubsets }

func (b *BackTrack) getAcceptedSubsets(node Node) {
	//
}

func (b *BackTrack) validateOrder(it []int) bool {
	// for i := 0; i < len(it)-1; i++ {
	// 	if it[i] >= it[i+1] {
	// 		return false
	// 	}
	// }

	return true
}

func backtrack(b *BackTrack, c []int) (Node, bool) {
	if b.reject(b.p, c) {
		return Node{}, false
	}

	node := Node{
		Subset:   c,
		Children: []Node{},
	}

	if b.accept(b.p, c, b.k) {
		(*b.acceptedCount)++
		return node, true
	}
	s := b.first(b.p, c)

	for s != nil {
		if child, ok := backtrack(b, s); ok {
			node.Children = append(node.Children, child)
		}
		s = b.next(b.p, s)
	}

	return node, true
}

func NewBackTrack(n []int, l int) *BackTrack {
	count := 0

	AcceptedSubsets := make([][]int, 0)

	p := slices.Clone(n)
	slices.Sort(p)

	backtrackk := BackTrack{p: p, rot: []int{}, k: l, acceptedCount: &count, AcceptedSubsets: AcceptedSubsets}

	backtrackk.rot = backtrackk.root(backtrackk.p)

	_, _ = backtrack(&backtrackk, []int{})

	return &backtrackk
}
