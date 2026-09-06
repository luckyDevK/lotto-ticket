package main

import (
	"cmp"
	"fmt"
	"impl-skiena/lotto"
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

	lotto.Permutations()
}
