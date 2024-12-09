package main

import (
    "fmt"
    "log"
    "os"
    "regexp"
    "strconv"
    "strings"
)

func main() {

    path := os.Args[1]
    data, err := os.ReadFile(path)
    if err != nil {
        log.Fatalf("Failed to read file: %v", err)
    }

    //sss := "xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))"

    s := string(data)
    var result int
    for _, s1 := range strings.Split(s, "do()") {
        l2 := strings.Split(s1, "don't()")
        fmt.Printf("%s\n", l2[0])
        result += parseString(l2[0])
    }

    fmt.Printf("Result: %d\n", result)
}

func parseString(s string) int {
    var result int
    re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)

    matches := re.FindAllStringSubmatch(s, -1)
    fmt.Printf("\nArguments: ")
    for _, match := range matches {
        if len(match) == 3 {
            fmt.Printf("(%s, %s) ", match[1], match[2])
            a1, err := strconv.Atoi(match[1])
            a2, err2 := strconv.Atoi(match[2])

            if err != nil || err2 != nil {
                continue
            }
            result += a1 * a2
        }
    }
    fmt.Println()
    fmt.Println()
    return result
}
