package main

// Adversarial / property tests for ClosestPair, on top of closestpair_test.go.
//
// Same oracle as there: all points are collinear, so an optimal closed tour
// costs exactly 2*(max-min), and the nearest-endpoint merge heuristic is
// provably optimal on a line (no two merges can cross). So any cost above the
// oracle is an implementation defect, never heuristic slack.
//
// What this file attacks that a hand-written table cannot:
//
//   A1 input order      fragments are built in INDEX order, not value order
//                       (addEdge appends the subset as-is), so a descending or
//                       scrambled input is what forces head/head and tail/tail
//                       merges in connectMultiVertex. Sorted inputs never
//                       reach those two branches at all.
//   A2 exhaustiveness   every point set drawn from a small universe, rather
//                       than the handful someone thought of.
//   A3 shrinking        a random failure is reduced to a minimal one you can
//                       trace on paper.
//   A4 algebra          cost is invariant under translation and reflection and
//                       scales linearly, so a sign or abs bug shows up as a
//                       broken identity rather than a single odd number.
//   A5 purity           ClosestPair must not mutate the caller's slice and must
//                       be deterministic. Both are live risks: fragments hold
//                       the subset slices owned by bt.AcceptedSubsets, and
//                       reverse() rewrites them in place.

import (
	"fmt"
	"math/rand"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// cpTourQuiet runs ClosestPair and parses the tour. ok is false on a panic or
// on unparsable output; no reporting, so it is safe inside search loops.
func cpTourQuiet(p []int) (tour []int, ok bool) {
	ok = true
	defer func() {
		if recover() != nil {
			tour, ok = nil, false
		}
	}()
	for _, f := range strings.Split(ClosestPair(slices.Clone(p)), " - ") {
		n, err := strconv.Atoi(strings.TrimSpace(f))
		if err != nil {
			return nil, false
		}
		tour = append(tour, n)
	}
	return tour, ok
}

// cpVerdict returns "" when p is handled correctly, or the reason it is not.
func cpVerdict(p []int) string {
	tour, ok := cpTourQuiet(p)
	if !ok {
		return "panicked or produced unparsable output"
	}
	if len(tour) != len(p)+1 {
		return fmt.Sprintf("want %d vertices (n+1), got %d: %v", len(p)+1, len(tour), tour)
	}
	if tour[0] != tour[len(tour)-1] {
		return fmt.Sprintf("tour not closed (%d != %d): %v", tour[0], tour[len(tour)-1], tour)
	}
	visited, want := slices.Clone(tour[:len(tour)-1]), slices.Clone(p)
	sort.Ints(visited)
	sort.Ints(want)
	if !slices.Equal(visited, want) {
		return fmt.Sprintf("not a permutation of the input (visited %v): %v", visited, tour)
	}
	if got, opt := cpTourCost(tour), cpOptimalCost(p); got != opt {
		return fmt.Sprintf("cost %d, optimal %d: %v", got, opt, tour)
	}
	return ""
}

// cpShrink reduces a failing input to a minimal still-failing one by dropping
// one point at a time, so a 20-point random failure lands as something you can
// hand-trace. Returns p unchanged if nothing smaller fails.
func cpShrink(p []int) []int {
	cur := slices.Clone(p)
	for shrunk := true; shrunk && len(cur) > 2; {
		shrunk = false
		for i := range cur {
			cand := slices.Concat(slices.Clone(cur[:i]), slices.Clone(cur[i+1:]))
			if cpVerdict(cand) != "" {
				cur, shrunk = cand, true
				break
			}
		}
	}
	return cur
}

// -- A1: input order -------------------------------------------------------

// TestClosestPairInputOrder feeds the SAME point set in every possible order.
// Fragment construction follows index order, so these 720-permutation sweeps
// are the only thing that exercises the head/head and tail/tail merge branches
// and any stale-index bookkeeping after a fragment is deleted from the slice.
func TestClosestPairInputOrder(t *testing.T) {
	bases := [][]int{
		// four fragments, all gaps distinct -- no ties to hide behind
		{0, 1, 10, 12, 30, 33},
		// gaps decrease left to right: fragments are born right to left
		{0, 8, 20, 25, 40, 42},
		// two 3-vertex fragments merging multi-to-multi
		{0, 1, 2, 20, 21, 22},
		// every internal gap tied at 1, every bridge tied at 9: pure tie-break
		{0, 1, 10, 11, 20, 21},
		// merge partners sit on opposite sides of zero
		{-50, -49, 0, 4, 100, 102},
		// one long fragment absorbing a far singleton last
		{0, 1, 2, 3, 4, 50},
		// perfectly uniform: every pair at distance 5 is a candidate
		{0, 5, 10, 15, 20, 25},

		// the exact mirror of the first base. Same fragment structure, every
		// coordinate negated: if one of this pair passes and the other fails,
		// the defect is in sign handling, not in the merge bookkeeping.
		{-33, -30, -12, -10, -1, 0},
		// strictly negative, so no coordinate is ever >= 0 at any point
		{-33, -30, -12, -10, -2, -1},
		// fragments born on both sides of zero, then bridged ACROSS it, so a
		// merge has to join a negative endpoint to a positive one
		{-12, -10, -1, 1, 10, 12},
	}

	for _, base := range bases {
		t.Run(fmt.Sprint(base), func(t *testing.T) {
			n, failed := 0, 0
			var firstBad []int
			cpPermute(slices.Clone(base), 0, func(order []int) {
				n++
				if cpVerdict(order) != "" {
					failed++
					if firstBad == nil {
						firstBad = slices.Clone(order)
					}
				}
			})
			if failed > 0 {
				t.Errorf("%d/%d input orderings of %v failed; e.g. %v: %s",
					failed, n, base, firstBad, cpVerdict(firstBad))
			}
		})
	}
}

// TestClosestPairDescending gives every table-worthy shape in reverse and in a
// fixed scramble, which is where index-ordered fragments bite.
func TestClosestPairDescending(t *testing.T) {
	shapes := [][]int{
		{0, 1, 10, 12, 30, 33, 60, 64},
		{0, 1, 5, 10, 11},
		{0, 1, 2, 20, 21, 22, 40, 41, 42},
		{0, 1, 1000, 1001, 2000, 2001, 3000, 3001},
		{-20, -19, -1, 0, 1, 19, 20},
		{0, 1, 7, 8, 9, 23, 24, 40, 41, 42, 55, 70, 71, 90},
		// mirrors of the above, so reversing/scrambling is exercised on
		// all-negative and zero-straddling coordinates too
		{-64, -60, -33, -30, -12, -10, -1, 0},
		{-90, -71, -70, -55, -42, -41, -40, -24, -23, -9, -8, -7, -1},
		{-31, -30, -1, 1, 30, 31},
	}

	for _, s := range shapes {
		desc := slices.Clone(s)
		slices.Reverse(desc)
		t.Run("desc "+fmt.Sprint(s), func(t *testing.T) {
			if why := cpVerdict(desc); why != "" {
				t.Errorf("input %v: %s", desc, why)
			}
		})

		// deterministic scramble: interleave from both ends inward, so no
		// fragment is ever born from adjacent indices
		scram := make([]int, 0, len(s))
		for i, j := 0, len(s)-1; i <= j; i, j = i+1, j-1 {
			scram = append(scram, s[i])
			if i != j {
				scram = append(scram, s[j])
			}
		}
		t.Run("scramble "+fmt.Sprint(s), func(t *testing.T) {
			if why := cpVerdict(scram); why != "" {
				t.Errorf("input %v: %s", scram, why)
			}
		})
	}
}

// -- A2: exhaustive small universes ---------------------------------------

// TestClosestPairExhaustive enumerates EVERY point set of size 3..6 drawn from
// two small coordinate universes. One universe is uniform, so ties dominate and
// the run is decided entirely by tie-break order; the other has strictly
// growing gaps, so every merge distance is distinct and the fragment order is
// forced. Together that is a few hundred inputs no hand-written table covers.
func TestClosestPairExhaustive(t *testing.T) {
	universes := map[string][]int{
		"uniform":      {0, 1, 2, 3, 4, 5, 6, 7},
		"growing gaps": {0, 1, 3, 6, 10, 15, 21, 28},
		// straddles zero, so a subset may land entirely left of it, entirely
		// right of it, or across it -- and when it lands across, every merge
		// distance is a difference of opposite-signed coordinates
		"straddling zero": {-12, -7, -3, -1, 0, 1, 4, 9},
		// strictly negative: no subset drawn from this universe contains a
		// single non-negative coordinate, which is the case the other three
		// universes cannot produce at all
		"all negative": {-29, -22, -16, -11, -7, -4, -2, -1},
		// the exact mirror of "uniform". Any subset here is the negation of a
		// subset there with the same optimal cost, so the two must agree.
		"uniform mirrored": {-7, -6, -5, -4, -3, -2, -1, 0},
	}

	for name, universe := range universes {
		t.Run(name, func(t *testing.T) {
			checked, failed := 0, 0
			var firstBad []int
			for k := 3; k <= 6; k++ {
				cpCombinations(universe, k, func(set []int) {
					checked++
					if cpVerdict(set) != "" {
						failed++
						if firstBad == nil {
							firstBad = slices.Clone(set)
						}
					}
				})
			}
			if failed > 0 {
				t.Errorf("%d/%d subsets of %v failed; e.g. %v: %s",
					failed, checked, universe, firstBad, cpVerdict(firstBad))
			}
			if checked == 0 {
				t.Fatal("enumerated nothing -- the combination generator is broken")
			}
		})
	}
}

// -- A3: deep fuzz with shrinking ------------------------------------------

// TestClosestPairFuzzDeep is the wide net: several structural generators, sizes
// from the degenerate end up to 20 points, and a shrinker so a failure arrives
// as the smallest input that still reproduces it.
func TestClosestPairFuzzDeep(t *testing.T) {
	gens := map[string]func(r *rand.Rand) []int{
		// scattered distinct coordinates over a wide window
		"scattered": func(r *rand.Rand) []int {
			n := 2 + r.Intn(19)
			seen := map[int]bool{}
			p := []int{}
			for len(p) < n {
				v := r.Intn(2001) - 1000
				if !seen[v] {
					seen[v] = true
					p = append(p, v)
				}
			}
			return p
		},
		// clusters separated by gaps much larger than the internal spacing:
		// every intra-cluster pair is consumed before any bridge, so the run is
		// a long stretch of births followed by a long stretch of merges
		"clustered": func(r *rand.Rand) []int {
			clusters := 2 + r.Intn(4)
			p := []int{}
			// start left of zero by a random amount, so a run either stays
			// entirely negative or crosses zero partway through the clusters
			base := -r.Intn(900)
			for c := 0; c < clusters; c++ {
				size := 1 + r.Intn(4)
				step := 1 + r.Intn(3)
				for i := 0; i < size; i++ {
					p = append(p, base+i*step)
				}
				base += 100 + r.Intn(400)
			}
			r.Shuffle(len(p), func(i, j int) { p[i], p[j] = p[j], p[i] })
			return p
		},
		// gaps that double: every merge distance exceeds the sum of all smaller
		// ones, so there is exactly one legal merge order and any mis-step shows
		"exponential": func(r *rand.Rand) []int {
			n := 3 + r.Intn(10)
			// shift the whole staircase so it can sit entirely left of zero,
			// entirely right of it, or across it, without changing the gap
			// structure that makes the merge order unique
			p := []int{r.Intn(1201) - 600}
			gap := 1
			for i := 1; i < n; i++ {
				p = append(p, p[i-1]+gap)
				gap *= 2
			}
			if r.Intn(2) == 0 {
				slices.Reverse(p)
			} else {
				r.Shuffle(len(p), func(i, j int) { p[i], p[j] = p[j], p[i] })
			}
			return p
		},
		// equal spacing at a random stride: maximum ties at every step
		"uniform": func(r *rand.Rand) []int {
			n := 3 + r.Intn(16)
			step := 1 + r.Intn(9)
			// wide enough that a whole run frequently lands below zero, not just
			// straddling it
			off := r.Intn(241) - 180
			p := []int{}
			for i := 0; i < n; i++ {
				p = append(p, off+i*step)
			}
			r.Shuffle(len(p), func(i, j int) { p[i], p[j] = p[j], p[i] })
			return p
		},
		// mirrored around zero, so merges repeatedly cross the sign boundary
		"mirrored": func(r *rand.Rand) []int {
			half := 1 + r.Intn(5)
			seen := map[int]bool{}
			p := []int{}
			for len(p) < 2*half {
				v := 1 + r.Intn(60)
				if !seen[v] {
					seen[v] = true
					p = append(p, v, -v)
				}
			}
			r.Shuffle(len(p), func(i, j int) { p[i], p[j] = p[j], p[i] })
			return p
		},
	}

	const perGen = 1500

	for name, gen := range gens {
		t.Run(name, func(t *testing.T) {
			r := rand.New(rand.NewSource(0xC0FFEE))
			failures := 0
			var smallest []int
			for i := 0; i < perGen; i++ {
				p := gen(r)
				if len(p) < 2 {
					continue
				}
				if cpVerdict(p) == "" {
					continue
				}
				failures++
				min := cpShrink(p)
				if smallest == nil || len(min) < len(smallest) {
					smallest = min
				}
			}
			if failures > 0 {
				t.Errorf("%d/%d %s inputs failed; smallest counterexample %v: %s",
					failures, perGen, name, smallest, cpVerdict(smallest))
			}
		})
	}
}

// -- A4: algebraic identities ----------------------------------------------

// TestClosestPairInvariance checks that the cost obeys the identities the
// geometry demands: translating the whole set leaves it alone, reflecting it
// leaves it alone, and scaling by k multiplies it by k. A single wrong abs,
// sign comparison or unsigned subtraction breaks one of these while still
// looking plausible on any one input.
func TestClosestPairInvariance(t *testing.T) {
	bases := [][]int{
		{0, 1, 10, 12, 30, 33},
		{0, 1, 5, 10, 11},
		{3, 1, 4, 1 + 6, 9},
		{0, 8, 20, 25, 40, 42, 60, 61},
	}

	cost := func(t *testing.T, p []int) int {
		t.Helper()
		tour, ok := cpTourQuiet(p)
		if !ok {
			t.Fatalf("input %v: panicked or produced unparsable output", p)
		}
		return cpTourCost(tour)
	}

	for _, base := range bases {
		t.Run(fmt.Sprint(base), func(t *testing.T) {
			want := cost(t, base)
			if opt := cpOptimalCost(base); want != opt {
				t.Fatalf("input %v: cost %d, optimal %d", base, want, opt)
			}

			for _, shift := range []int{-1000, -7, 7, 1000, 1 << 30} {
				moved := make([]int, len(base))
				for i, v := range base {
					moved[i] = v + shift
				}
				if got := cost(t, moved); got != want {
					t.Errorf("translation by %d changed cost: %d != %d (%v)", shift, got, want, moved)
				}
			}

			mirror := make([]int, len(base))
			for i, v := range base {
				mirror[i] = -v
			}
			if got := cost(t, mirror); got != want {
				t.Errorf("reflection changed cost: %d != %d (%v)", got, want, mirror)
			}

			for _, k := range []int{2, 13, 1000} {
				scaled := make([]int, len(base))
				for i, v := range base {
					scaled[i] = v * k
				}
				if got := cost(t, scaled); got != want*k {
					t.Errorf("scaling by %d: cost %d, want %d (%v)", k, got, want*k, scaled)
				}
			}
		})
	}
}

// -- A5: purity ------------------------------------------------------------

// TestClosestPairDoesNotMutateInput pins down that the caller's slice survives
// the call. fragments store the very slices owned by bt.AcceptedSubsets and
// reverse() rewrites them in place, so aliasing here is not hypothetical.
func TestClosestPairDoesNotMutateInput(t *testing.T) {
	for _, p := range [][]int{
		{0, 1, 10, 12, 30, 33},
		{33, 30, 12, 10, 1, 0},
		{0, 1, 2, 20, 21, 22, 40, 41, 42},
	} {
		t.Run(fmt.Sprint(p), func(t *testing.T) {
			in := slices.Clone(p)
			before := slices.Clone(in)
			_ = ClosestPair(in)
			if !slices.Equal(in, before) {
				t.Errorf("input mutated: %v -> %v", before, in)
			}
		})
	}
}

// TestClosestPairDeterministic calls the same input repeatedly. Map iteration
// or an unstable sort leaking into the tie-break shows up here as flapping
// output even though every single run looks optimal.
func TestClosestPairDeterministic(t *testing.T) {
	for _, p := range [][]int{
		{0, 1, 10, 11, 20, 21, 30, 31},
		{0, 5, 10, 15, 20, 25, 30},
		{21, 0, 31, 10, 1, 30, 11, 20},
	} {
		t.Run(fmt.Sprint(p), func(t *testing.T) {
			first := ClosestPair(slices.Clone(p))
			for i := 0; i < 20; i++ {
				if got := ClosestPair(slices.Clone(p)); got != first {
					t.Fatalf("run %d differs: %q != %q", i+1, got, first)
				}
			}
		})
	}
}

// -- stress ----------------------------------------------------------------

// TestClosestPairStress runs sizes where bookkeeping errors compound instead of
// cancelling. Uniform spacing is deliberate: it maximises ties, so the greedy
// has the largest possible set of equally short candidates at every step.
func TestClosestPairStress(t *testing.T) {
	for _, n := range []int{16, 20, 24, 28} {
		t.Run(fmt.Sprintf("uniform n=%d", n), func(t *testing.T) {
			p := make([]int, n)
			for i := range p {
				// centred on zero, so half the ladder is negative at every n
				p[i] = (i - n/2) * 3
			}
			if why := cpVerdict(p); why != "" {
				t.Errorf("n=%d ascending: %s", n, why)
			}

			rev := slices.Clone(p)
			slices.Reverse(rev)
			if why := cpVerdict(rev); why != "" {
				t.Errorf("n=%d descending: %s", n, why)
			}

			r := rand.New(rand.NewSource(int64(n)))
			for trial := 0; trial < 5; trial++ {
				sh := slices.Clone(p)
				r.Shuffle(len(sh), func(i, j int) { sh[i], sh[j] = sh[j], sh[i] })
				if why := cpVerdict(sh); why != "" {
					t.Errorf("n=%d shuffled %v: %s", n, sh, why)
				}
			}
		})
	}

	// clusters at huge magnitude: checks that nothing truncates or overflows on
	// the way through, while staying far from MaxInt so the oracle is exact
	t.Run("large magnitudes", func(t *testing.T) {
		for _, p := range [][]int{
			{-1000000000, -999999999, 999999999, 1000000000},
			{1 << 40, 1<<40 + 1, 1 << 41, 1<<41 + 1},
			{-(1 << 40), 0, 1 << 40},
			{-(1 << 41), -(1 << 41) + 1, -(1 << 40), -(1 << 40) + 1},
			{-1000000000, -999999998, -999999996},
		} {
			if why := cpVerdict(p); why != "" {
				t.Errorf("input %v: %s", p, why)
			}
		}
	})
}

// -- degenerate sizes ------------------------------------------------------

// TestClosestPairDegenerateSize covers the two sizes where the output shape is
// a contract decision rather than a consequence of the algorithm. Both are
// counterexamples this file turned up; every other property above held across
// ~10k adversarial inputs.
//
//	n == 1  the tour is "5 - 5": one vertex, repeated so the output parses like
//	        every other size (count n+1, first == last). Returning a bare "5"
//	        would break a caller that parses uniformly on exactly this input.
//	n == 0  there is NO tour, so it cannot be checked against P1/P2 at all --
//	        n+1 = 1 vertex is not a thing an empty point set has. ClosestPair
//	        signals it with the sentinel "-", so this subtest asserts that
//	        string directly instead of routing through cpVerdict, which parses
//	        its argument as a tour and can only report "unparsable".
//
// If the n == 0 sentinel ever changes -- to "", to a second return value, or to
// a documented panic -- this subtest is the one place that has to change with
// it, and cpVerdict stays out of it either way.
func TestClosestPairDegenerateSize(t *testing.T) {
	t.Run("n=1", func(t *testing.T) {
		if why := cpVerdict([]int{5}); why != "" {
			t.Errorf("input [5]: %s", why)
		}
	})

	t.Run("n=0", func(t *testing.T) {
		const wantSentinel = "-"

		got, panicked := "", false
		func() {
			defer func() {
				if r := recover(); r != nil {
					panicked = true
					got = fmt.Sprint(r)
				}
			}()
			got = ClosestPair([]int{})
		}()

		if panicked {
			t.Fatalf("input []: panicked (%s), want sentinel %q", got, wantSentinel)
		}
		if got != wantSentinel {
			t.Errorf("input []: got %q, want sentinel %q", got, wantSentinel)
		}
	})

	t.Run("n=2", func(t *testing.T) {
		for _, p := range [][]int{{0, 1}, {1, 0}, {-5, 5}, {7, 7}} {
			if why := cpVerdict(p); why != "" {
				t.Errorf("input %v: %s", p, why)
			}
		}
	})
}

// -- generators ------------------------------------------------------------

// cpPermute calls f once per permutation of p, reusing the same backing slice,
// so f must copy anything it keeps.
func cpPermute(p []int, i int, f func([]int)) {
	if i == len(p) {
		f(p)
		return
	}
	for j := i; j < len(p); j++ {
		p[i], p[j] = p[j], p[i]
		cpPermute(p, i+1, f)
		p[i], p[j] = p[j], p[i]
	}
}

// cpCombinations calls f once per k-subset of universe, in ascending order.
// f must copy anything it keeps.
func cpCombinations(universe []int, k int, f func([]int)) {
	cur := make([]int, 0, k)
	var rec func(start int)
	rec = func(start int) {
		if len(cur) == k {
			f(cur)
			return
		}
		for i := start; i < len(universe); i++ {
			cur = append(cur, universe[i])
			rec(i + 1)
			cur = cur[:len(cur)-1]
		}
	}
	rec(0)
}
