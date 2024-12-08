package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
    "sync"
    "sync/atomic"
)

type LevelType int

const (
    LevelInvalid LevelType = iota
    LevelIncreasing
    LevelDecreasing
)

type level []int16

func getLevelType(l level) LevelType {
    var increasing, decreasing int
    var incCalcStop, decCalcStop bool
    llen := len(l)
    for i := 1; i < llen; i++ {
        diff := l[i] - l[i-1]
        if !incCalcStop && diff > 0 && diff < 4 {
            increasing++
        } else if !decCalcStop && diff < 0 && diff > -4 {
            decreasing++
        }
        if increasing != i {
            incCalcStop = true
        }
        if decreasing != i {
            decCalcStop = true
        }

        if incCalcStop && decCalcStop {
            return LevelInvalid
        }
    }
    if !incCalcStop {
        return LevelIncreasing
    }
    return LevelDecreasing
}

func main() {

    path := os.Args[1]
    scanner, closef := fileReader(path)
    defer closef()
    var wg sync.WaitGroup
    var result uint32

    for scanner.Scan() {
        wg.Add(1)
        go func(line string) {
            defer wg.Done()
            l := checkLevel(line)
            if l == LevelIncreasing || l == LevelDecreasing {
                atomic.AddUint32(&result, 1)
            }
        }(scanner.Text())
    }
    wg.Wait()

    fmt.Printf("%d reports are safe\n", result)
}

func checkLevel(s string) LevelType {
    stringList := strings.Fields(s)

    var l level
    for _, s := range stringList {
        num, err := strconv.Atoi(s)
        if err != nil {
            panic(fmt.Sprintf("Error converting string to integer: %v", err))
        }
        l = append(l, int16(num))
    }

    return getLevelType(l)
}

func fileReader(path string) (*bufio.Scanner, func()) {
    file, err := os.Open(path)

    if err != nil {
        panic(fmt.Sprintf("error opening file: %v", err))
    }

    f := func() { file.Close() }

    return bufio.NewScanner(file), f
}