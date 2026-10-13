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

func getItemsArr(arr, it []int) (int, int) {
	fmt.Println(it, "endpoint")
	var existingVertex, differVert int
	for i := 0; i < len(it); i++ {
		exist := false
		for j := 0; j < len(arr); j++ {
			if arr[j] == it[i] {
				existingVertex = arr[j]
				exist = true
			}
		}

		if !exist {
			differVert = it[i]
		}
	}

	return existingVertex, differVert
}

func containsAll(subset, chain []int) bool {
	//fmt.Println(chain, "chaintos")
	countExist := 0
	for i := 0; i < len(chain); i++ {
		for j := 0; j < len(subset); j++ {
			if chain[i] == subset[j] {
				countExist++
			}
		}
	}

	//fmt.Println(countExist == len(subset), chain, "countExist")
	return countExist == len(subset)
}

func isStartItem(it int, arr []int) bool {
	return len(arr) > 0 && arr[0] == it
}

func getDist(s, t int) int {
	if s > t {
		return s - t
	} else {
		return t - s
	}
}

func connectByEdge(chain []int, endp []int) []int {
	if len(chain) == 0 {
		//fmt.Println(chain, "chain")
		chain = append(chain, endp...)
		fmt.Println("initial chain ->", chain)
		return chain

	}
	existingVertex, differVert := getItemsArr(chain, endp)

	if isStartItem(existingVertex, chain) {
		fmt.Println(differVert, "<-", chain)
		chain = append([]int{differVert}, chain...)
	} else {
		fmt.Println(chain, "->", differVert)
		chain = append(chain, differVert)
	}

	return chain
}

func connectEndpoint(chain []int) []int {
	leftEndp := chain[0]
	chain = append(chain, leftEndp)
	return chain
}

func ClosestPair(p []int) string {
	var chain []int

	bt := biocoeff.NewBackTrack(p, 2)

	for i := 0; i < len(p)-1; i++ {
		dist := math.MaxInt
		sub := []int{}

		for _, subset := range bt.AcceptedSubsets {
			if containsAll(subset, chain) {
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
		chain = connectByEdge(chain, sub)
	}

	chain = connectEndpoint(chain)

	lenStrChain := len(chain) - 2

	chainStr := make([]string, lenStrChain)
	for i := range lenStrChain {
		chainStr[i] = strconv.Itoa(chain[i])
	}

	lastOne := strconv.Itoa(chain[len(chain)-2])
	lastTwo := strconv.Itoa(chain[len(chain)-1])

	twoEndp := []string{lastOne, lastTwo}

	chainS := slices.Concat(chainStr, twoEndp)

	return strings.Join(chainS, " - ")
}

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
	fmt.Println(ch, "final-chain")

	// lotto.Permutations()
}
