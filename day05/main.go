package main

import (
    "log"
    "os"
    "strconv"
    "strings"
    "sync"
    "sync/atomic"
)

type Rule [2]int

func NewRuleFromStr(str string) Rule {
    pair := strings.Split(str, "|")
    r1, _ := strconv.Atoi(pair[0])
    r2, _ := strconv.Atoi(pair[1])
    return Rule{r1, r2}
}

type Rules []Rule

func NewRulesFromStr(str string) Rules {
    rules := make(Rules, 0)
    for _, s := range strings.Split(str, "\n") {
        rules = append(rules, NewRuleFromStr(s))
    }
    return rules
}

func (r Rule) first() int {
    return r[0]
}

func (r Rule) second() int {
    return r[1]
}

type Update []int

func NewUpdateFromStr(str string) Update {
    update := make(Update, 0)
    for _, s := range strings.Split(str, ",") {
        v, _ := strconv.Atoi(s)
        update = append(update, v)
    }
    return update
}

func (u Update) ApplyRules(rules Rules) int {
    var swapCounter int
    for _, rule := range rules {
        el1, ok1 := u.getPosition(rule.first())
        el2, ok2 := u.getPosition(rule.second())

        if ok1 && ok2 && el1 > el2 {
            u[el1], u[el2] = u[el2], u[el1]
            swapCounter++
        }

    }
    return swapCounter
}

func (u Update) Middle() uint32 {
    return uint32(u[len(u)/2])
}

func (u Update) getPosition(el int) (int, bool) {
    for i, e := range u {
        if e == el {
            return i, true
        }
    }
    return -1, false
}

type Updates []Update

func NewUpdatesFromStr(str string) Updates {
    lines := strings.Split(str, "\n")
    updates := make(Updates, len(lines))
    for i, s := range lines {
        updates[i] = NewUpdateFromStr(s)
    }
    return updates
}

func getInput(path string) (Rules, Updates) {
    data, err := os.ReadFile(path)

    if err != nil {
        panic("failed to read file: " + path)
    }

    parts := strings.Split(strings.TrimSpace(string(data)), "\n\n")
    return NewRulesFromStr(parts[0]), NewUpdatesFromStr(parts[1])
}

func main() {

    path := os.Args[1]
    rules, updates := getInput(path)

    var result uint32
    var wg sync.WaitGroup
    for _, update := range updates {
        wg.Add(1)
        go func() {
            defer wg.Done()
            var i uint32
            for update.ApplyRules(rules) > 0 {
                i++
            }
            //if i == 0 { // part #1
            if i > 0 { // part #2
                atomic.AddUint32(&result, update.Middle())
            }
        }()
    }
    wg.Wait()

    log.Printf("Result: %d\n", result)
}
