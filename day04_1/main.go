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

type substrCalculatorHorizontal struct{}

func (c substrCalculatorHorizontal) count(s InputText, substr SearchWord) uint32 {
	d := uint32(strings.Count(s.text, string(substr)))
	r := uint32(strings.Count(s.text, string(substr.reverse())))
	fmt.Printf("substrCalculatorHorizontal direct = %d, reverse = %d\n", d, r)

	return d + r
}

type substrCalculatorVertical struct{}

func (c substrCalculatorVertical) count(s InputText, substr SearchWord) uint32 {

	var result uint32
	var wg sync.WaitGroup

	wg.Add(s.cols)
	for i := 0; i < s.cols; i++ {
		go func() {
			defer wg.Done()
			result += c.countInColumn(i, s, substr)
		}()
	}
	wg.Wait()
	return result
}

func (c substrCalculatorVertical) countInColumn(x int, s InputText, substr SearchWord) uint32 {
	subLen := substr.len()
	rsubstr := substr.reverse()

	//fmt.Printf("substr: %s subLen %d col: %d, rows: %d cols: %d, llen: %d\n", substr, subLen, x, s.rows, s.cols, s.len)
	var result uint32

	for y := 0; y < s.rows-subLen; y++ {
		//fmt.Printf("i: %d\n", i)
		var j, dcount, rcount int

		for j = 0; j < subLen; j++ {
			chr, err := s.XY(x, y+j)
			if err != nil {
				break
			}
			//fmt.Printf("%d, chr: %c, substr: %c\n", j, chr, rune(substr[j]))
			if chr == rune(substr[j]) {
				dcount++
			}

			if chr == rune(rsubstr[j]) {
				rcount++
			}
		}

		if dcount == subLen {
			//fmt.Printf("dcount: col %d, row %d \n", x, y)
			result++
		}

		if rcount == subLen {
			//fmt.Printf("rcount: col %d, row %d \n", x, y)
			result++
		}
	}
	fmt.Printf("countInColumn x = %d result: %d\n", x, result)
	return result
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
	fmt.Printf("countLeftToRight y = %d result: %d\n", y, result)
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
	fmt.Printf("countRightToLeft y = %d result: %d\n", y, result)
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
	var wg sync.WaitGroup
	wg.Add(4)

	var result uint32

	go func() {
		defer wg.Done()
		h := substrCalculatorHorizontal{}
		atomic.AddUint32(&result, h.count(s, ss))
	}()

	go func() {
		defer wg.Done()
		v := substrCalculatorVertical{}
		atomic.AddUint32(&result, v.count(s, ss))
	}()

	go func() {
		defer wg.Done()
		d := substrCalculatorDiagonalLeftToRight{}
		atomic.AddUint32(&result, d.count(s, ss))
	}()

	go func() {
		defer wg.Done()
		d := substrCalculatorDiagonalRightToLeft{}
		atomic.AddUint32(&result, d.count(s, ss))
	}()

	wg.Wait()

	log.Printf("Result: %d\n", result)

}
