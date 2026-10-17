package main

import (
	"cmp"
	"fmt"
	biocoeff "impl-skiena/bioCoeff"
	"math"
	"slices"
)

/* [5] Implement the two TSP heuristics of Section 1.1 (page 5). Which of them gives
better-quality solutions in practice? Can you devise a heuristic that works better
than both of them?
*/

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

func getItemsArr(arr [][]int, it []int) (int, int) {
	var existingVertex, differVert int
	for i := 0; i < len(it); i++ {
		exist := false
		for j := 0; i < len(arr); i++ {
			for k := 0; k < len(arr[i]); j++ {
				if arr[i][j] == it[i] {
					existingVertex = arr[i][j]
					exist = true
				}
			}

		}
		if !exist {
			differVert = it[i]
		}
	}

	return existingVertex, differVert
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

func getFromPosItem(it int, arr [][]int, pos string) (int, bool) {
	for i := 0; i < len(arr); i++ {
		if pos == "head" {
			if len(arr[i]) > 0 && arr[i][0] == it {
				return i, true
			}
		} else if pos == "tail" {
			if len(arr[i]) > 0 && arr[i][len(arr[i])-1] == it {
				return i, true
			}
		}
	}
	return -1, false
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
		fmt.Println("initial fragments ->", fragments)
		return fragments
	}

	// check if both endp absent add new into fragment

	// for F each fragments
	for i := 0; i < len(fragments); i++ {
		exist := false
		for j := 0; j < len(fragments[i]); j++ {
			for k := 0; k < len(endp); k++ {
				if endp[k] == fragments[i][j] {
					exist = true
				}
			}
		}
		if !exist {
			fragments = append(fragments, endp)
			return fragments
		}
	}

	// check if one of the endp present add with singleton
	existingVertex, differVert := getItemsArr(fragments, endp)

	if index, ok := getFromPosItem(existingVertex, fragments, "head"); ok {
		fmt.Println(differVert, "<-", fragments)
		fragments[index] = append([]int{differVert}, fragments[index]...)
		return fragments
	} else if index, ok := getFromPosItem(existingVertex, fragments, "tail"); ok {
		fmt.Println(fragments, "->", differVert)
		fragments[index] = append(fragments[index], differVert)
		return fragments
	}

	// check if both vertices exist but in differ fragment

	fragHoldEndp := make(map[int]int)
	for k := 0; k < len(endp); k++ {
		if index, ok := getFromPosItem(endp[k], fragments, "head"); ok {
			fragHoldEndp[index] = 0
		} else if index, ok := getFromPosItem(endp[k], fragments, "tail"); ok {
			fragHoldEndp[index] = len(fragments[index]) - 1
		}
	}

	connectMultiVertex()
	// connect multi-fragments

	return fragments
}

func connectMultiVertex(fragHoldEndp map[int]int) {

}

func connectEndpoint(fragments []int) []int {
	leftEndp := fragments[0]
	fragments = append(fragments, leftEndp)
	return fragments
}

func ClosestPair(p []int) string {
	var fragments [][]int

	bt := biocoeff.NewBackTrack(p, 2)

	for i := 0; i < len(p)-1; i++ {
		dist := math.MaxInt
		sub := []int{}

		for _, subset := range bt.AcceptedSubsets {
			if containsAll(subset, fragments) {
				continue
			}

			firstSub := subset[0]
			lastSub := subset[1]
			//fmt.Println(firstSub, lastSub, "firstSub And lastSub")

			if getDist(firstSub, lastSub) <= dist {
				sub = subset
				dist = getDist(firstSub, lastSub)
			}
		}
		//fmt.Println("sub and dist", sub, dist)
		fragments = addEdge(fragments, sub)
	}

	// fragments = connectEndpoint(fragments)

	// lenStrChain := len(fragments) - 2

	// chainStr := make([]string, lenStrChain)
	// for i := range lenStrChain {
	// 	chainStr[i] = strconv.Itoa(fragments[i])
	// }

	// lastOne := strconv.Itoa(fragments[len(fragments)-2])
	// lastTwo := strconv.Itoa(fragments[len(fragments)-1])

	// twoEndp := []string{lastOne, lastTwo}

	// chainS := slices.Concat(chainStr, twoEndp)

	return "2222"
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

	ch := ClosestPair([]int{-10, 1, -5, 2})
	fmt.Println(ch, "final-fragments")

	// lotto.Permutations()
}
