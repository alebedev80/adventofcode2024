package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"sync"
	"sync/atomic"
)

type pair [2]uint32

func main() {

	path := os.Args[1]
	list, _ := getOrderedIntListsFromFile(path)
	fmt.Printf("%v\n", list)
	fmt.Printf("Result: %d\n", getSummaryDistance(&list))
}

func getOrderedIntListsFromFile(path string) ([]pair, error) {

	var l1, l2 []int
	var result []pair
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var v1, v2 int
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d    %d", &v1, &v2)
		l1 = append(l1, v1)
		l2 = append(l2, v2)
		slices.Sort(l1)
		slices.Sort(l2)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	for i := 0; i < len(l1); i++ {
		result = append(result, pair{uint32(l1[i]), uint32(l2[i])})
	}

	return result, nil
}

func getSummaryDistance(list *[]pair) uint32 {
	var sum uint32
	var wg sync.WaitGroup

	for _, val := range *list {
		wg.Add(1)
		go func(p pair) {
			defer wg.Done()
			if p[0] > p[1] {
				atomic.AddUint32(&sum, p[0]-p[1])
			} else {
				atomic.AddUint32(&sum, p[1]-p[0])
			}
		}(val)
	}

	wg.Wait()

	return sum
}
