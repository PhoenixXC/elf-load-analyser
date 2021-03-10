package env

import (
    _ "embed" // embed for banner.txt
    "encoding/json"
    "fmt"
    "github.com/phoenixxc/elf-load-analyser/pkg/helper"
    "strings"
)

const (
    kvDelimiter = ": \t"
    bannerLen   = 38 + 10
)

//go:embed banner.txt
var banner string

//go:embed status.json
var status string

type Entry struct {
    Key string `json:"key"`
    Val string `json:"val"`
}

func (e Entry) parse() string {
    return e.Key + kvDelimiter + e.Val
}

type Status []Entry

func EchoBanner() {
    kvList := parseStatus()
    if len(kvList) == 0 {
        fmt.Println(banner)
    }

    bannerArr := strings.Split(banner, "\n")
    kvLen, banLen := len(kvList), len(bannerArr)
    startIdx := helper.IfElse(kvLen >= banLen, 0, (banLen-kvLen)/2).(int)

    for i, s := range bannerArr {
        fmt.Printf("%-*.*s", bannerLen, bannerLen, s)
        if bIdx := i - startIdx; bIdx >= 0 && bIdx < kvLen {
            fmt.Print(kvList[bIdx])
        }
        fmt.Println()
    }
    for i := banLen; i < kvLen; i++ {
        fmt.Printf("%-.*s\n", banLen+len(kvList[i]), kvList[i])
    }
}

func parseStatus() (kv []string) {
    s := Status{}
    err := json.Unmarshal([]byte(status), &s)
    if err != nil {
        return []string{}
    }
    entries := []Entry(s)
    kv = make([]string, len(entries))
    for idx, entry := range entries {
        kv[idx] = entry.parse()
    }
    return
}
