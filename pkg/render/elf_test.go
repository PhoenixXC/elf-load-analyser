package render

import (
    "fmt"
    "testing"
)

func TestElfRender(t *testing.T) {
    d, _ := NewElfRender("/bin/ls").Render()
    fmt.Println(d.DataStr())
}
