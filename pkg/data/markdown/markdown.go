package markdown

import (
    "bytes"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "strings"
)

type ItemLevel uint8

const (
    H1 = ItemLevel(iota + 1)
    H2
    H3
    H4
)

type Interface interface {
    ToMarkdown() string
}

type Markdown struct {
    Interface
}

func (m Markdown) Class() data.Type {
    return data.MarkdownType
}

func (m Markdown) Data() string {
    return m.Interface.ToMarkdown()
}

// Content simple context
type Content struct {
    Markdown
    level   ItemLevel
    title   string
    content string
}

// NewContent 创建 markdown内容块
func NewContent() *Content {
    c := &Content{}
    c.Markdown.Interface = c
    return c
}

func NewTitleContents(level ItemLevel, title string) *Content {
    return NewContent().WithTitle(level, title)
}

func NewTextContent(content ...string) *Content {
    return NewContent().WithContents(content...)
}

func (m *Content) WithTitle(level ItemLevel, title string) *Content {
    m.level = level
    m.title = title
    return m
}

func (m *Content) WithContents(content ...string) *Content {
    m.content += strings.Join(content, "\n\n")
    return m
}

func (m *Content) Append(am Interface) *Content {
    appendContent := am.ToMarkdown()
    var buf bytes.Buffer
    buf.Grow(len(m.content) + len(appendContent))
    buf.WriteString(m.content)
    buf.WriteString(appendContent)

    return NewContent().WithTitle(m.level, m.title).WithContents(buf.String())
}

func (m *Content) ToMarkdown() string {
    var buf bytes.Buffer
    buf.Grow(int(m.level) + len(m.title) + len(m.content) + 5)
    if len(m.title) > 0 && m.level != 0 {
        for k := 0; k < int(m.level); k++ {
            buf.WriteRune('#')
        }
        buf.WriteRune(' ')
        buf.WriteString(m.title)
        buf.WriteString("\n\n")
    }
    buf.WriteString(m.content)
    buf.WriteString("\n\n")
    return buf.String()
}

// Table table format
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

func (t *Table) ToMarkdown() string {
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

// List Markdown 列表
type List struct {
    list   []string
    length int
}

func NewList(list ...string) *List {
    return &List{list: list}
}

func NewListFromContent(contents ...Markdown) *List {
    l, conLen := &List{}, len(contents)
    l.list = make([]string, conLen)
    l.length = conLen*3 + 2
    for _, content := range contents {
        c := content.ToMarkdown()
        l.length += len(c)
        l.list = append(l.list, c)
    }
    return l
}

func (l List) ToMarkdown() string {
    var buf bytes.Buffer
    if l.length != 0 {
        buf.Grow(l.length)
    }
    for _, s := range l.list {
        buf.WriteString("- ")
        buf.WriteString(s)
        buf.WriteRune('\n')
    }
    buf.WriteString("\n\n")
    return buf.String()
}
