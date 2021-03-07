package render

import (
    "fmt"
    "testing"
)

func TestEnvRender(t *testing.T) {
    d, _ := NewEnvRender().Render()
    fmt.Println(d.DataStr())
}
