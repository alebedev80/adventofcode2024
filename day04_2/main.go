package main

import (
    "fmt"
    "log"
    "os"
    "strings"
    "sync"
    "sync/atomic"
    //"sync"
    //"sync/atomic"
)

type SearchWord string

func (s SearchWord) reverse() SearchWord {

    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }

    return SearchWord(runes)
}

func (s SearchWord) len() int {
    return len(s)
}

type InputText struct {
    text string
    rows int
    cols int
    len  int
}

func NewInputText(text string) InputText {
    return InputText{
        text: strings.TrimSpace(text),
        rows: strings.Count(text, "\n") + 1,
        cols: strings.Index(text, "\n"),
        len:  len(text),
    }
}

func (s InputText) XY(x int, y int) (rune, error) {
    pos := y*s.cols + y + x
    if y*s.cols+y+x >= s.len {
        return 0, fmt.Errorf("invalid coordinates: %d, %d", x, y)
    }
    return rune(s.text[pos]), nil
}

type substrCalculator interface {
    count(s InputText, substr SearchWord) uint32
}

type substrXCalculator struct{}

func (c substrXCalculator) count(s InputText, substr SearchWord) uint32 {
    var result uint32
    var wg sync.WaitGroup

    leftToRight := substrCalculatorDiagonalLeftToRight{}
    rightToLeft := substrCalculatorDiagonalRightToLeft{}

    for x := 0; x <= s.cols-substr.len(); x++ {
        for y := 0; y <= s.rows-substr.len(); y++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                frame, err := c.getFrame(s, x, y, substr.len())
                if err != nil {
                    return
                }
                l2r := leftToRight.count(frame, substr)
                r2l := rightToLeft.count(frame, substr)
                if l2r > 0 && r2l > 0 {
                    fmt.Printf("Found X, cords %d,%d\n", x, y)
                    atomic.AddUint32(&result, 1)
                }
            }()

        }
    }

    wg.Wait()
    return result
}

func (c substrXCalculator) getFrame(s InputText, x int, y int, size int) (InputText, error) {
    text := ""
    for i := 0; i < size; i++ {
        for j := 0; j < size; j++ {
            chr, err := s.XY(x+j, y+i)
            if err != nil {
                return InputText{}, err
            }
            text += string(chr)
        }
        text += "\n"
    }
    return NewInputText(text), nil

}

type substrCalculatorDiagonalLeftToRight struct{}

func (c substrCalculatorDiagonalLeftToRight) count(s InputText, substr SearchWord) uint32 {
    var result uint32
    var wg sync.WaitGroup

    wg.Add(s.rows)
    for i := 0; i < s.rows; i++ {
        go func() {
            defer wg.Done()
            result += c.countDiagonal(i, s, substr)
        }()
    }
    wg.Wait()
    return result
}

func (c substrCalculatorDiagonalLeftToRight) countDiagonal(y int, s InputText, substr SearchWord) uint32 {
    subLen := substr.len()
    rsubstr := substr.reverse()

    var result uint32

    for x := 0; x <= s.cols-subLen; x++ {
        var j, dcount, rcount int

        for j = 0; j < subLen; j++ {

            chr, err := s.XY(x+j, y+j)
            if err != nil {
                break
            }

            if chr == rune(substr[j]) {
                dcount++
            }

            if chr == rune(rsubstr[j]) {
                rcount++
            }
        }

        if dcount == subLen {
            result++
        }

        if rcount == subLen {
            result++
        }
    }
    //fmt.Printf("countLeftToRight y = %d result: %d\n", y, result)
    return result

}

type substrCalculatorDiagonalRightToLeft struct{}

func (c substrCalculatorDiagonalRightToLeft) count(s InputText, substr SearchWord) uint32 {
    var result uint32
    var wg sync.WaitGroup

    wg.Add(s.rows)
    for i := 0; i < s.rows; i++ {
        go func() {
            defer wg.Done()
            result += c.countDiagonal(i, s, substr)
        }()
    }
    wg.Wait()
    return result
}
func (c substrCalculatorDiagonalRightToLeft) countDiagonal(y int, s InputText, substr SearchWord) uint32 {
    subLen := substr.len()
    rsubstr := substr.reverse()

    var result uint32

    for x := s.cols - 1; x >= subLen-1; x-- {
        var j, dcount, rcount int

        for j = 0; j < subLen; j++ {
            chr, err := s.XY(x-j, y+j)
            if err != nil {
                break
            }

            if chr == rune(substr[j]) {
                dcount++
            }

            if chr == rune(rsubstr[j]) {
                rcount++
            }
        }

        if dcount == subLen {
            result++
        }

        if rcount == subLen {
            result++
        }
    }
    //fmt.Printf("countRightToLeft y = %d result: %d\n", y, result)
    return result
}

func main() {

    ss := SearchWord(os.Args[1])
    path := os.Args[2]
    data, err := os.ReadFile(path)
    if err != nil {
        log.Fatalf("Failed to read file: %v", err)
    }

    s := NewInputText(string(data))
    c := substrXCalculator{}
    result := c.count(s, ss)

    log.Printf("Result: %d\n", result)

}
