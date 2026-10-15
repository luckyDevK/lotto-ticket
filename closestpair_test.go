package main

// Tests for ClosestPair.
//
// All points here lie on a line, so the oracle is simple: an optimal tour over
// collinear points costs exactly 2*(max-min) -- you must walk the span and come
// back. Two properties are checked:
//
//   P1 (validity)    the returned tour is a permutation of the input, closed
//                    (first vertex repeated at the end), length n+1.
//   P2 (optimality)  its cost equals 2*(max-min).
//
// A heuristic is not obliged to be optimal in general, but on collinear points
// ClosestPair is: every merge joins two nearest fragment endpoints, which on a
// line never crosses. So a P2 failure here is a real defect, not bad luck.

import (
	"fmt"
	"math/rand"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func cpParseTour(t *testing.T, s string) []int {
	t.Helper()
	out := []int{}
	for _, f := range strings.Split(s, " - ") {
		n, err := strconv.Atoi(strings.TrimSpace(f))
		if err != nil {
			t.Fatalf("unparsable tour %q: %v", s, err)
		}
		out = append(out, n)
	}
	return out
}

func cpTourCost(tour []int) int {
	total := 0
	for i := 0; i+1 < len(tour); i++ {
		total += getDist(tour[i], tour[i+1])
	}
	return total
}

func cpOptimalCost(p []int) int {
	q := slices.Clone(p)
	sort.Ints(q)
	return 2 * (q[len(q)-1] - q[0])
}

// cpRun calls ClosestPair and turns a panic into a test failure.
func cpRun(t *testing.T, p []int) (tour []int, ok bool) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("input %v: panic: %v", p, r)
			ok = false
		}
	}()
	return cpParseTour(t, ClosestPair(slices.Clone(p))), true
}

// cpCheck verifies P1 and P2. Returns false on the first property violated.
func cpCheck(t *testing.T, p []int) bool {
	t.Helper()

	tour, ok := cpRun(t, p)
	if !ok {
		return false
	}

	if len(tour) != len(p)+1 {
		t.Errorf("input %v -> %v: want %d vertices (n+1), got %d", p, tour, len(p)+1, len(tour))
		return false
	}
	if tour[0] != tour[len(tour)-1] {
		t.Errorf("input %v -> %v: tour not closed (%d != %d)", p, tour, tour[0], tour[len(tour)-1])
		return false
	}

	visited := slices.Clone(tour[:len(tour)-1])
	want := slices.Clone(p)
	sort.Ints(visited)
	sort.Ints(want)
	if !slices.Equal(visited, want) {
		t.Errorf("input %v -> %v: not a permutation of the input (visited %v)", p, tour, visited)
		return false
	}

	if got, opt := cpTourCost(tour), cpOptimalCost(p); got != opt {
		t.Errorf("input %v -> %v: cost %d, optimal %d", p, tour, got, opt)
		return false
	}
	return true
}

func TestClosestPairTable(t *testing.T) {
	cases := [][]int{
		// trivial sizes
		{7, 3},
		{0, 1, 2},

		// already sorted / evenly spaced
		{1, 2, 3, 4},
		{-3, 0, 3, 6, 9},

		// the input you were tracing by hand
		{-10, 1, -5, 2},

		// two tight clusters far apart: the second cluster is a pair whose two
		// vertices are both absent from the chain when it gets picked
		{0, 1, 100, 101},
		{0, 1, 2, 50, 51, 52},

		// three clusters, all on the same side, equal internal spacing.
		// which cluster is picked second depends on subset iteration order.
		{0, 1, 10, 11, 20, 21},
		{0, 1, -10, -11, -20, -21},

		// nearest neighbour is "inward", so a greedy outward sweep goes wrong
		{0, 10, 20, 9},
		{-13, 5, 29, 14, -16},

		// the closest pair is in the middle of the span
		{-30, -1, 0, 1, 30},
		{-30, 0, 1, 30},

		// --- cases that need more than one fragment --------------------------
		// Everything below is chosen so that several fragments are alive at
		// once. A single-chain implementation cannot represent these at all.

		// the second-shortest pair has BOTH vertices absent -- the branch a
		// single chain has nowhere to put
		{0, 1, 10, 12},

		// the 5-point trace: born, born, absorb singleton, same-fragment
		// reject, then merge two multi-vertex fragments
		{0, 1, 5, 10, 11},

		// four fragments alive at once, all gaps distinct so there are no ties
		// to argue about: births at d=1,2,3,4 then merges at d=9,18,27
		{0, 1, 10, 12, 30, 33, 60, 64},

		// same shape, but internal gaps DECREASE left to right, so fragments
		// are born right to left and the collection is never value-ordered
		{0, 8, 20, 25, 40, 42, 60, 61},

		// the lower-valued fragment ends up LAST in path order: [0,4] joins
		// [-50,-49] at its tail. Fails anything that sorts a merge by value.
		{-50, -49, 0, 4, 100, 102},
		{0, 4, 50, 51, 100, 102},

		// three 3-vertex fragments, then two multi-vertex merges back to back.
		// Catches a stale fragment lookup for the absorbed fragment's vertices.
		{0, 1, 2, 20, 21, 22, 40, 41, 42},

		// five clusters, every internal gap tied at 1 and every bridge tied at
		// 9, so the whole run is decided by tie-break order
		{0, 1, 10, 11, 20, 21, 30, 31, 40, 41},

		// the bridge pair is in the middle and is picked first; the outer
		// singletons then attach to opposite ends
		{0, 30, 31, 61},

		// one tight cluster plus scattered singletons: every step after the
		// cluster forms is the one-present / absorb-singleton case
		{0, 1, 2, 3, 100, 200, 201},

		// the absorbed fragment is larger than the one it merges into
		{0, 1, 2, 3, 4, 50, 51},

		// fragments straddling zero, so signs change across a merge
		{-7, -6, 6, 7},
		{-20, -19, -1, 0, 1, 19, 20},

		// equidistant neighbours: the middle vertex can be taken either way
		{0, 5, 10},
		{-10, -5, 0, 5, 10},

		// large span, tiny clusters: every intra-cluster pair is consumed
		// before any bridge is even considered
		{0, 1, 1000, 1001, 2000, 2001, 3000, 3001},

		// big enough that bookkeeping errors compound instead of cancelling
		{0, 1, 7, 8, 9, 23, 24, 40, 41, 42, 55, 70, 71, 90},
	}

	for _, c := range cases {
		t.Run(fmt.Sprint(c), func(t *testing.T) {
			cpCheck(t, slices.Clone(c))
		})
	}
}

// TestClosestPairFuzz samples random point sets with distinct coordinates and
// reports the failure rate plus the smallest counterexample, so you can shrink
// a failure down to something traceable by hand.
func TestClosestPairFuzz(t *testing.T) {
	const iterations = 2000

	r := rand.New(rand.NewSource(1))
	failures := 0
	var smallest []int

	for i := 0; i < iterations; i++ {
		n := 3 + r.Intn(5)
		seen := map[int]bool{}
		p := []int{}
		for len(p) < n {
			v := r.Intn(61) - 30
			if !seen[v] {
				seen[v] = true
				p = append(p, v)
			}
		}

		if !cpCheckQuiet(p) {
			failures++
			if smallest == nil || len(p) < len(smallest) {
				smallest = slices.Clone(p)
			}
		}
	}

	if failures > 0 {
		t.Errorf("%d/%d random inputs failed; smallest counterexample: %v -> %q (optimal cost %d)",
			failures, iterations, smallest, ClosestPair(slices.Clone(smallest)), cpOptimalCost(smallest))
	}
}

// cpCheckQuiet is cpCheck without the reporting, for use inside the fuzz loop.
func cpCheckQuiet(p []int) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()

	fields := strings.Split(ClosestPair(slices.Clone(p)), " - ")
	tour := make([]int, 0, len(fields))
	for _, f := range fields {
		n, err := strconv.Atoi(strings.TrimSpace(f))
		if err != nil {
			return false
		}
		tour = append(tour, n)
	}

	if len(tour) != len(p)+1 || tour[0] != tour[len(tour)-1] {
		return false
	}

	visited := slices.Clone(tour[:len(tour)-1])
	want := slices.Clone(p)
	sort.Ints(visited)
	sort.Ints(want)
	if !slices.Equal(visited, want) {
		return false
	}

	return cpTourCost(tour) == cpOptimalCost(p)
}

// TestClosestPairDuplicates is a separate test because supporting repeated
// coordinates is a design decision, not a given. Delete this test if you decide
// duplicates are out of scope; keep it if you want value-indexed lookups
// replaced by position-indexed ones.
func TestClosestPairDuplicates(t *testing.T) {
	for _, c := range [][]int{
		{4, 4},
		{5, 5, 1, 9},
		{0, 0, 0},
	} {
		t.Run(fmt.Sprint(c), func(t *testing.T) {
			cpCheck(t, slices.Clone(c))
		})
	}
}
