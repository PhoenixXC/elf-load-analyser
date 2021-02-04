package data

const (
    MarkdownType = iota
    GraphvizType
)

type Data struct {
    Class int // Data type, MarkdownType or GraphvizType
    Data  string
    Style string // css style
}

func NewData(class int, data string) *Data {
    return &Data{Class: class, Data: data}
}

func (d *Data) WithStyle(css string)  {
    d.Style = css
}
