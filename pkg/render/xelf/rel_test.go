package xelf

import (
    "debug/elf"
    "fmt"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

func TestBuildRela(t *testing.T) {
    f, err := getELFFile("./testdata/rel.o")
    if err != nil {
        log.Error(err)
    }
    sra, err := BuildRelIf(f, false)
    if err != nil {
        log.Error(err)
    }
    assert.Equal(t, 2, len(sra))
    textSra := sra[0]
    assert.Equal(t, ".rela.text", textSra.Section.Name)
    assert.Equal(t, 3, len(textSra.Rels))
    assert.Equal(t, uint64(0x25), textSra.Rels[0].Offset)
    assert.Equal(t, "shared", textSra.Rels[0].Value)
    f.Close()
}

func TestBuildRela_dyn(t *testing.T) {
    f, err := getELFFile("./testdata/lib.so")
    if err != nil {
        log.Error(err)
    }
    sr, err := BuildRelIf(f, true)
    if err != nil {
        log.Error(err)
    }
    fmt.Println(sr)
}

func getELFFile(path string) (*elf.File, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    eFile, err := elf.NewFile(f)
    if err != nil {
        return nil, err
    }
    return eFile, nil
}
