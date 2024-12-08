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

func (l level) getLevelType() LevelType {
    var increasing, decreasing int
    var incCalcStop, decCalcStop bool
    for i := 1; i < l.len(); i++ {
        diff := l[i] - l[i-1]
        if !incCalcStop && diff > 0 && diff < 4 {
            increasing++
        } else if !decCalcStop && diff < 0 && diff > -4 {
            decreasing++
        }

        incCalcStop = increasing != i
        decCalcStop = decreasing != i

        if incCalcStop && decCalcStop {
            return LevelInvalid
        }
    }
    if !incCalcStop {
        return LevelIncreasing
    }
    return LevelDecreasing
}

func (l level) len() int {
    return len(l)
}

func (l level) levelWithOutElement(i int) level {
    var newLevel level
    for j, v := range l {
        if j != i {
            newLevel = append(newLevel, v)
        }
    }
    return newLevel
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
            if checkLevel(line) != LevelInvalid {
                atomic.AddUint32(&result, 1)
            }
        }(scanner.Text())
    }
    wg.Wait()

    fmt.Printf("%d reports are safe\n", result)
}

func checkLevel(s string) LevelType {
    stringList := strings.Fields(s)

    l := level{}
    for _, s := range stringList {
        num, err := strconv.Atoi(s)
        if err != nil {
            panic(fmt.Sprintf("Error converting string to integer: %v", err))
        }
        l = append(l, int16(num))
    }

    result := LevelInvalid
    for i := 0; result == LevelInvalid && i < l.len(); i++ {
        result = l.levelWithOutElement(i).getLevelType()
    }
    return result
}

func fileReader(path string) (*bufio.Scanner, func()) {
    file, err := os.Open(path)

    if err != nil {
        panic(fmt.Sprintf("error opening file: %v", err))
    }

    f := func() { file.Close() }

    return bufio.NewScanner(file), f
}
