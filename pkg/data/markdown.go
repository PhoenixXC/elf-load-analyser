package data

import (
    "bytes"
    "strings"
)

type ItemLevel uint8

const (
    H1Level = ItemLevel(1)
    H2Level = ItemLevel(2)
    H3Level = ItemLevel(3)
    H4Level = ItemLevel(4)
)

type Markdown interface {
    String() string
}

// Item Markdown simple context
type Item struct {
    level   ItemLevel
    title   string
    content string
}

func NewItem(level ItemLevel, title string, content ...string) *Item {
    return &Item{level: level, title: title, content: strings.Join(content, "\n")}
}

func (i *Item) String() string {
    var buf bytes.Buffer
    buf.Grow(int(i.level) + len(i.title) + len(i.content) + 5)
    for k := 0; k < int(i.level); k++ {
        buf.WriteRune('#')
    }
    buf.WriteRune(' ')
    buf.WriteString(i.title)
    buf.WriteRune('\n')
    buf.WriteRune('\n')
    buf.WriteString(i.content)
    buf.WriteRune('\n')
    buf.WriteRune('\n')
    return buf.String()
}

// Table Markdown
type Table struct {
    col     int
    desc    string
    head    []string
    content [][]string
}

func NewTable(head ...string) *Table {
    return &Table{head: head, col: len(head), content: [][]string{}}
}

func (t *Table) WithDesc(desc string) *Table {
    t.desc = desc
    return t
}

func (t *Table) AddRow(rows ...string) *Table {
    if len(rows) < t.col {
        for len(rows) < t.col {
            rows = append(rows, "—")
        }
    }
    r := rows[:t.col]
    t.content = append(t.content, r)
    return t
}

func (t *Table) String() string {
    var buf, line bytes.Buffer
    ch := '|'
    for _, h := range t.head {
        buf.WriteRune(ch)
        buf.WriteString(h)
        line.WriteString("|---")
    }
    buf.WriteString("|\n")
    line.WriteString("|\n")
    buf.WriteString(line.String())

    for _, row := range t.content {
        cnt := 0
        for _, c := range row {
            cnt++
            buf.WriteRune(ch)
            buf.WriteString(c)
        }
        buf.WriteString("|\n")
    }
    if len(t.desc) > 0 {
        buf.WriteString("\n<small>" + t.desc + "</small>")
    }
    buf.WriteString("\n\n")
    return buf.String()
}
