package main

import (
    "bufio"
    "fmt"
    "os"
    "sync"
    "sync/atomic"
)

type aggregatedList map[int]int

func main() {

    path := os.Args[1]
    list1, list2, _ := getAggregatedListsFromFile(path)
    result := getSimilarityScore(list1, list2)
    fmt.Printf("Similarity Score = %d\n", result)
}

func getAggregatedListsFromFile(path string) (aggregatedList, aggregatedList, error) {

    scanner, closef := fileReader(path)
    defer closef()

    var l1, l2 = make(aggregatedList), make(aggregatedList)
    var v1, v2 int
    for scanner.Scan() {
        fmt.Sscanf(scanner.Text(), "%d    %d", &v1, &v2)
        addToAggList(l1, v1)
        addToAggList(l2, v2)
    }

    return l1, l2, nil
}

func fileReader(path string) (*bufio.Scanner, func()) {
    file, err := os.Open(path)

    if err != nil {
        panic(fmt.Sprintf("error opening file: %v", err))
    }

    f := func() { file.Close() }

    return bufio.NewScanner(file), f
}

func addToAggList(l aggregatedList, v int) {
    if _, ok := l[v]; !ok {
        l[v] = 1
    } else {
        l[v]++
    }
}

func getSimilarityScore(l1 aggregatedList, l2 aggregatedList) int {
    var score int
    for k, v := range l1 {
        if val, ok := l2[k]; ok {
            score += k * v * val
        }
    }
    return score

}

func getSummaryDistance(l1 aggregatedList, l2 aggregatedList) uint32 {
    var sum uint32
    var wg sync.WaitGroup

    for k, v := range l1 {
        go func() {
            if val, ok := l2[k]; ok {
                atomic.AddUint32(&sum, uint32(k*v*val))
            }
        }()

    }

    wg.Wait()

    return sum
}
