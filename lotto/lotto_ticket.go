package lotto

import (
	"fmt"
	"slices"
)

// P = {1,2,3,4,5,6,7,8,9,10,11...}

// []

// root (P) = []
// reject (P,c) = return false always
// accept (P,c) = return true if the candidate length same as input else false
// first (P,c) = generate the first extension of canidate of c for example the candidate [1,2,3] = [1],[2],[3]
// next (P,s) = generate the next candidate after s extension

// procedure backtrack(P, c) is
//     if reject(P, c) then return
//     if accept(P, c) then output(P, c)
//     s ← first(P, c)
//     while s ≠ NULL do
//         backtrack(P, s)
//         s ← next(P, s)

// create reject with return always false
// accept when len is same as input
// firtst -> create the first extension of c
// next -> generate the next extension of a can after s (first)

type Backtracking interface {
	root(p []int) []int
	reject(P []int, c []int) bool
	accept(P []int, c []int, k int)
	first(P []int, c []int) []int
	next(P []int, s []int) []int
	validateOrder(it []int) bool
	generateTickets()
}

type BackTrack struct {
	p, rot  []int
	k       int
	path    string
	content map[string]any

	acceptedSubsets [][]int
	acceptedCount   *int
}

type Node struct {
	Subset   []int  `json:"subset"`
	Children []Node `json:"children"`
}

type LottoSearch struct {
	totalSubsets int
}

type SATicket struct {
	noImprovementStreak, iterationPerTemp int
}

func (ls *LottoSearch) LottTicketSet(n []int, k, l int) [][]int {
	// initialize the (n l)-element bit-vector V to all false
	V := make([]bool, ls.totalSubsets)

	// while there are exist a false entry in V
	for ls.hasUncovered(V) {
		// Select k-subset T of of n as the next ticket to buy
		// For each of the l-subsets Ti of T, V[rank(Ti)] = true
		// report the set of tickets bought
	}

	return [][]int{}
}

func (s *SATicket) GenerateTicket(kTicket []int) {
	// currentS = kTicket
	// initialize temp = 1
	// while noImprovementStreak <= 10 do
	// for i until iterationPerTemp do
	//
}

func (s *SATicket) Transition() {

}

func (s *SATicket) C(ticket []int) {
	//
}

func (b *LottoSearch) hasUncovered(V []bool) bool {
	for _, covered := range V {
		if !covered {
			return true
		}
	}

	return false
}

func (b *BackTrack) root(p []int) []int {
	if !b.validateOrder(p) {
		return nil
	}

	return []int{}
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
		b.acceptedSubsets = append(b.acceptedSubsets, append([]int(nil), c...))
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

func (b *BackTrack) getAcceptedSubsets(node Node) {
	//
}

func (b *BackTrack) validateOrder(it []int) bool {
	for i := 0; i < len(it)-1; i++ {
		if it[i] >= it[i+1] {
			return false
		}
	}

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

func Permutations() {
	count := 0

	n := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	l := 3
	// k := 6

	acceptedSubsets := make([][]int, 0)

	backtrackk := BackTrack{p: n, rot: []int{}, k: l, acceptedCount: &count, acceptedSubsets: acceptedSubsets}

	backtrackk.rot = backtrackk.root(backtrackk.p)

	_, _ = backtrack(&backtrackk, []int{})

	//fmt.Println("ss", backtrackk.acceptedSubsets)

	for _, subset := range backtrackk.acceptedSubsets {
		fmt.Printf("subset: %v\n", subset)
	}

	fmt.Println("accepted total:", *backtrackk.acceptedCount)
	// fmt.Println("ven:", root)

	// backtrackk.content = map[string]any{
	// 	"root": root,
	// }

	// backtrackk.LottTicketSet(n, l, k)

	// // for k, v := range backtrackk.acceptedSubset.BVNode {
	// // 	fmt.Printf("subset k:%v val:%v\n", k, v)
	// // }

	// file, err := os.Create("subsets.json")
	// if err != nil {
	// 	log.Fatalf("Failed to create file: %s", err)
	// }

	// defer file.Close()

	// encoder := json.NewEncoder(file)
	// encoder.SetIndent("", " ")
	// err = encoder.Encode(backtrackk.content)

	// if err != nil {
	// 	log.Fatalf("Failed to Encode file: %s", err)
	// }

}

//  {
//     "root": {
//       "subset": [],
//       "children": [
//         {
//           "subset": [1],
//           "children": [
//             {
//               "subset": [1,2],
//               "children": [
//                 { "subset": [1,2,3], "children": [] },
//                 { "subset": [1,2,4], "children": [] }
//               ]
//             },
//             {
//               "subset": [1,3],
//               "children": [
//                 { "subset": [1,3,4], "children": [] }
//               ]
//             },
//             {
//               "subset": [1,4],
//               "children": []
//             }
//           ]
//         }
//       ]
//     }
//   }

// subset problem
// accept when c is the len is equal to k

// procedure backtrack(P, c) is
//     if reject(P, c) then return
//     if accept(P, c) then output(P, c)
//     s ← first(P, c)
//     while s ≠ NULL do
//         backtrack(P, s)
//         s ← next(P, s)

// create reject with return always false
// accept when len is same as input
// firtst -> create the first extension of c
// next -> generate the next extension of a can after s (first)
