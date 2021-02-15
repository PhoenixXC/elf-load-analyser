package render

import (
    "fmt"
    "testing"
)

func TestElfRender(t *testing.T) {
    render := NewElfRender("/bin/ls")
    data, _ := render.Render()
    fmt.Println(data.Name())
    fmt.Println(data.Data())
}
