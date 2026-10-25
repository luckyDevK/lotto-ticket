package main

import (
	"fmt"
	"strings"
	"time"
)

// CPTrace is the JSON-friendly form of a ClosestPair run: everything needed
// to animate it in a browser, with actual point values instead of indices.
type CPTrace struct {
	Name   string       `json:"name"`
	Points []int        `json:"points"`
	Steps  []CPJSONStep `json:"steps"`
	Tour   string       `json:"tour"`
}

type CPJSONStep struct {
	Added     []int   `json:"added"`
	Dist      int     `json:"dist"`
	Fragments [][]int `json:"fragments"`
}

// BuildCPTrace runs ClosestPair on p and converts the recorded steps (which
// reference p by index) into a trace that records actual values, suitable
// for json.Marshal and use outside this process (e.g. a browser viz).
func BuildCPTrace(name string, p []int) CPTrace {
	var steps []CPStep
	tour := ClosestPair(append([]int{}, p...), func(s CPStep) {
		steps = append(steps, s)
	})

	jsteps := make([]CPJSONStep, len(steps))
	for i, s := range steps {
		added := []int{p[s.Added[0]], p[s.Added[1]]}
		frags := make([][]int, len(s.Fragments))
		for j, f := range s.Fragments {
			vals := make([]int, len(f))
			for k, idx := range f {
				vals[k] = p[idx]
			}
			frags[j] = vals
		}
		jsteps[i] = CPJSONStep{Added: added, Dist: s.Dist, Fragments: frags}
	}

	return CPTrace{Name: name, Points: p, Steps: jsteps, Tour: tour}
}

// ansi fragment colors, cycled if there are more fragments than colors.
var cpColors = []string{
	"\033[31m", // red
	"\033[32m", // green
	"\033[33m", // yellow
	"\033[34m", // blue
	"\033[35m", // magenta
	"\033[36m", // cyan
}

const cpReset = "\033[0m"
const cpBold = "\033[1m"

// cpClearScreen moves the cursor home and clears the terminal.
const cpClearScreen = "\033[H\033[2J"

// cpRenderer holds everything needed to draw a single frame for a fixed
// input p at a fixed terminal width.
type cpRenderer struct {
	p        []int
	width    int
	min, max int
	axisLine string
}

func newCPRenderer(p []int, width int) *cpRenderer {
	min, max := p[0], p[0]
	for _, v := range p {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	axis := strings.Repeat("-", width)
	minLabel, maxLabel := fmt.Sprint(min), fmt.Sprint(max)
	axisLine := minLabel + axis[len(minLabel):width-len(maxLabel)] + maxLabel

	return &cpRenderer{p: p, width: width, min: min, max: max, axisLine: axisLine}
}

func (r *cpRenderer) col(v int) int {
	if r.max == r.min {
		return 0
	}
	return (v - r.min) * (r.width - 1) / (r.max - r.min)
}

// frame renders one line of colored points for the given fragments, with the
// just-added edge (if any) bolded and shown as "O" instead of "o".
func (r *cpRenderer) frame(fragments [][]int, added []int) string {
	cells := make([]string, r.width)
	for i := range cells {
		cells[i] = "."
	}
	for fi, frag := range fragments {
		color := cpColors[fi%len(cpColors)]
		for _, idx := range frag {
			c := r.col(r.p[idx])
			mark := "o"
			style := color
			if len(added) == 2 && (idx == added[0] || idx == added[1]) {
				style = color + cpBold
				mark = "O"
			}
			cells[c] = style + mark + cpReset
		}
	}
	return strings.Join(cells, "")
}

func (r *cpRenderer) fragmentLegend(fragments [][]int) string {
	fragStr := make([]string, len(fragments))
	for j, f := range fragments {
		vals := make([]int, len(f))
		for k, idx := range f {
			vals[k] = r.p[idx]
		}
		fragStr[j] = cpColors[j%len(cpColors)] + fmt.Sprint(vals) + cpReset
	}
	return strings.Join(fragStr, " ")
}

// VisualizeClosestPair runs ClosestPair on p and prints each merge step as a
// colored ASCII number line, one step per block (scrolls, like a log).
func VisualizeClosestPair(p []int, width int) {
	var steps []CPStep
	tour := ClosestPair(append([]int{}, p...), func(s CPStep) {
		steps = append(steps, s)
	})

	r := newCPRenderer(p, width)

	fmt.Printf("input: %v  (n=%d)\n\n", p, len(p))
	fmt.Println(r.axisLine)
	fmt.Println(strings.Repeat(".", width))
	fmt.Println()

	for i, s := range steps {
		fmt.Printf("step %d: added edge (%d,%d) d=%d\n", i+1, p[s.Added[0]], p[s.Added[1]], s.Dist)
		fmt.Println(r.frame(s.Fragments, s.Added))
		fmt.Println("fragments:", r.fragmentLegend(s.Fragments))
		fmt.Println()
	}

	fmt.Println("final tour:", tour)
}

// AnimateClosestPair runs ClosestPair on p and plays each merge step as an
// animation: every step clears the terminal and redraws the current frame in
// place, pausing delay between steps, like a flipbook.
func AnimateClosestPair(p []int, width int, delay time.Duration) {
	var steps []CPStep
	tour := ClosestPair(append([]int{}, p...), func(s CPStep) {
		steps = append(steps, s)
	})

	r := newCPRenderer(p, width)

	draw := func(stepNum int, s CPStep) {
		fmt.Print(cpClearScreen)
		fmt.Printf("input: %v  (n=%d)\n\n", p, len(p))
		fmt.Println(r.axisLine)
		if s.Fragments == nil {
			fmt.Println(strings.Repeat(".", width))
			fmt.Println()
			fmt.Println("start (no fragments yet)")
		} else {
			fmt.Println(r.frame(s.Fragments, s.Added))
			fmt.Println()
			fmt.Printf("step %d/%d: added edge (%d,%d) d=%d\n", stepNum, len(steps), p[s.Added[0]], p[s.Added[1]], s.Dist)
			fmt.Println("fragments:", r.fragmentLegend(s.Fragments))
		}
	}

	draw(0, CPStep{})
	time.Sleep(delay)
	for i, s := range steps {
		draw(i+1, s)
		time.Sleep(delay)
	}

	fmt.Println("\nfinal tour:", tour)
}
