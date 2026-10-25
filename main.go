package main

import (
	"cmp"
	"fmt"
	biocoeff "impl-skiena/bioCoeff"
	"math"
	"slices"
	"strconv"
	"strings"
)

/* [5] Implement the two TSP heuristics of Section 1.1 (page 5). Which of them gives
better-quality solutions in practice? Can you devise a heuristic that works better
than both of them?
*/

const (
	absent = iota
	tail
	head
	interior
)

// create object with these entity -> name : string , isVisited : bool, distance: int

type RobotTour struct {
	Name      string
	isVisited bool
	point     int
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func getIndexFromPoint(robotTours []RobotTour, targetPoint int) int {
	for i, v := range robotTours {
		if targetPoint == v.point {
			return i
		}
	}

	return -1
}

func countUnvisitedRobot(robotTours *[]RobotTour) int {
	unvisited := 0

	for i := range *robotTours {
		if !(*robotTours)[i].isVisited {
			unvisited++
		}
	}
	return unvisited
}

func isNegativeNum(point int) bool {
	return point < 0
}

func Visit(robotTours []RobotTour) {
	// get smallest absolute value of robotTours
	smallest := slices.MinFunc(robotTours, func(a, b RobotTour) int {
		return cmp.Compare(absInt(a.point), absInt(b.point))
	})
	// get smallest index
	index := getIndexFromPoint(robotTours, smallest.point)
	// pick initial p
	currVisit := &robotTours[index]
	fmt.Printf("[START] %s (%d)\n", robotTours[index].Name, robotTours[index].point)
	currVisit.isVisited = true
	no := 1
	// create for loop to visit unvissited robot
	for u := countUnvisitedRobot(&robotTours); u > 0; u = countUnvisitedRobot(&robotTours) {
		closestDist := math.MaxInt
		closestIndex := -1
		for i := range robotTours {
			if robotTours[i].isVisited {
				continue
			}

			currentDist := absInt(robotTours[i].point - currVisit.point)
			if currentDist <= closestDist {
				closestDist = currentDist
				closestIndex = i
			}
		}

		if closestIndex > -1 {
			fmt.Printf("%d. %v(%d) -> %v(%d)\n", no, currVisit.Name, currVisit.point, robotTours[closestIndex].Name, robotTours[closestIndex].point)
			no++
			currVisit = &robotTours[closestIndex]
		}

		if !currVisit.isVisited {
			//fmt.Printf("curr visited %s %d\n", currVisit.Name, currVisit.point)
			currVisit.isVisited = true
		}

		if u == 1 {
			fmt.Printf("revisit initial robot at the end: %s", currVisit.Name)
		}

	}
}

func containsAll(subset []int, fragments [][]int) bool {
	//fmt.Println(fragments, "chaintos")
	for i := 0; i < len(fragments); i++ {
		countExist := 0
		for j := 0; j < len(fragments[i]); j++ {
			for k := 0; k < len(subset); k++ {
				if subset[k] == fragments[i][j] {
					countExist++
				}

				if countExist == len(subset) {
					return true
				}
			}
		}
	}

	//fmt.Println(countExist == len(subset), fragments, "countExist")
	return false
}

type MarkerFrag struct {
	index, marker int
}

func getDist(s, t int) int {
	if s > t {
		return s - t
	} else {
		return t - s
	}
}

func addEdge(fragments [][]int, endp []int) [][]int {
	if len(fragments) == 0 {
		//fmt.Println(fragments, "fragments")
		fragments = append(fragments, endp)
		return fragments
	}

	// check if both vertices exist but in differ fragment
	fragHoldEndp := make([]MarkerFrag, len(endp))
	for i := 0; i < len(endp); i++ {
		mf := locate(endp[i], fragments)
		fragHoldEndp[i] = mf
	}

	a, b := fragHoldEndp[0], fragHoldEndp[1]

	// if both endp absent add into fragments
	if a.marker == absent && b.marker == absent {
		fragments = append(fragments, endp)
		return fragments
	}

	// if one of the endp exist in tail or head then add that new frag
	if (b.marker == absent) != (a.marker == absent) {
		present, newVertex := b, endp[0]
		if present.marker == absent {
			present, newVertex = a, endp[1]
		}

		switch present.marker {
		case head:
			fragments[present.index] = append([]int{newVertex}, fragments[present.index]...)
			return fragments
		case tail:
			fragments[present.index] = append(fragments[present.index], newVertex)
			return fragments
		}
	}

	return connectMultiVertex(fragHoldEndp, fragments)
}

func locate(v int, fragments [][]int) MarkerFrag {
	for i := 0; i < len(fragments); i++ {
		for j := 0; j < len(fragments[i]); j++ {
			if fragments[i][j] == v {
				switch {
				case j == 0:
					return MarkerFrag{i, head}
				case j == len(fragments[i])-1:
					return MarkerFrag{i, tail}
				default:
					return MarkerFrag{i, interior}
				}
			}
		}
	}
	return MarkerFrag{index: -1, marker: absent}
}

func connectMultiVertex(fragHoldEndp []MarkerFrag, fragments [][]int) [][]int {
	a, b := fragHoldEndp[0], fragHoldEndp[1]

	isAValidFrag := a.marker != absent && a.marker != interior
	isBValidFrag := b.marker != absent && b.marker != interior

	if !isAValidFrag && !isBValidFrag {
		return fragments
	}

	switch {
	case a.marker == head && b.marker == head:
		swp := reverse(fragments[a.index])
		fragments[b.index] = slices.Concat(swp, fragments[b.index])
		fragments = append(fragments[:a.index], fragments[a.index+1:]...)
	case a.marker == tail && b.marker == tail:
		swp := reverse(fragments[b.index])
		fragments[a.index] = slices.Concat(fragments[a.index], swp)
		fragments = append(fragments[:b.index], fragments[b.index+1:]...)
	default:
		if a.marker == tail && b.marker == head {
			fragments[a.index] = slices.Concat(fragments[a.index], fragments[b.index])
			fragments = append(fragments[:b.index], fragments[b.index+1:]...)
		} else if b.marker == tail && a.marker == head {
			fragments[a.index] = slices.Concat(fragments[b.index], fragments[a.index])
			fragments = append(fragments[:b.index], fragments[b.index+1:]...)
		}
	}

	return fragments
}

func reverse(arr []int) []int {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}

	return arr
}

func connectEndpoint(fragment []int) []int {
	leftEndp := fragment[0]
	fragment = append(fragment, leftEndp)
	return fragment
}

func isInterior(subset []int, fragments [][]int) bool {
	for i := 0; i < len(fragments); i++ {
		for j := 0; j < len(fragments[i]); j++ {
			for k := 0; k < len(subset); k++ {
				if len(fragments[i]) > 2 && subset[k] == fragments[i][j] && j > 0 && j < len(fragments[i])-1 {
					return true
				}
			}
		}
	}
	return false
}

func ClosestPair(p []int) string {
	var fragments [][]int

	if len(p) == 1 {
		return strconv.Itoa(p[0])
	}

	indicies := make([]int, len(p))
	for i := 0; i < len(p); i++ {
		indicies[i] = i
	}

	// use  inline loop rather than backtracking
	bt := biocoeff.NewBackTrack(indicies, 2)
	// fmt.Println(bt.AcceptedSubsets, "acceptedSubset")
	// fmt.Println(len(p), "pragos")

	// bAcc := make([][]int, len(bt.AcceptedSubsets))
	// // iterate over ba
	// for i, ba := range bt.AcceptedSubsets {
	// 	bs := make([]int, len(ba))
	// 	for j, b := range ba {
	// 		val := p[b]
	// 		bs[j] = val
	// 		// fmt.Println("abs", math.Abs(float64(p[i]-p[j])))
	// 	}
	// 	slices.Sort(bs)
	// 	bAcc[i] = bs

	// }

	slices.SortFunc(bt.AcceptedSubsets, func(a, b []int) int {
		aAbs := math.Abs(float64(p[a[0]] - p[a[1]]))
		bAbs := math.Abs(float64(p[b[0]] - p[b[1]]))

		return cmp.Compare(aAbs, bAbs)
	})

	// fmt.Println(bAcc, "bAcc")

	for i := 0; i < len(p)-1; i++ {

		dist := math.MaxInt
		sub := []int{}

		for _, subset := range bt.AcceptedSubsets {
			if containsAll(subset, fragments) {
				continue
			}

			if isInterior(subset, fragments) {
				continue
			}

			firstSub := subset[0]
			lastSub := subset[1]
			d := getDist(p[firstSub], p[lastSub])

			if d <= dist {
				sub = subset
				dist = d
			}
		}
		fragments = addEdge(fragments, sub)
	}

	fragments[0] = connectEndpoint(fragments[0])

	chainStr := make([]string, len(fragments[0]))
	for i := range chainStr {
		chainStr[i] = strconv.Itoa(p[fragments[0][i]])
	}

	// lastOne := strconv.Itoa(fragments[len(fragments)-2])
	// lastTwo := strconv.Itoa(fragments[len(fragments)-1])

	// twoEndp := []string{lastOne, lastTwo}

	// chainS := slices.Concat(chainStr, twoEndp)

	return strings.Join(chainStr, " - ")
}

// problem =
// closestPairDuplicate
//

func main() {

	// robots := []RobotTour{
	// 	{Name: "u1", point: -21},
	// 	{Name: "u2", point: -5},
	// 	{Name: "u3", point: -1},
	// 	{Name: "u4", point: 0},
	// 	{Name: "u5", point: 1},
	// 	{Name: "u6", point: 3},
	// 	{Name: "u7", point: 11},
	// }
	// Visit(robots)

	// bt := biocoeff.NewBackTrack([]int{5, 5, 1, 9}, 2)

	// fmt.Println(bt.AcceptedSubsets, "acceptedSubsets")
	ch := ClosestPair([]int{-13, -25, -16, 15, -18})
	fmt.Println(ch, "final-fragments")

	//lotto.Permutations()
}
