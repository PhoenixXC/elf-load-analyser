package render

import (
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "testing"
)

func TestElfRender_buildStaticData(t *testing.T) {
    r, e := NewElfRender("/tmp/test/rel.o")
    if e != nil {
        log.Error(e)
    }
    r.buildStaticData()
}
