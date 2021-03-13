package xelf

import (
    "debug/elf"
    "fmt"
    "runtime"
)

type RelResolver interface {
    String() string
}

type Arch struct {
    class int
}

var name2Arch = map[string]func(int) RelResolver{
    "amd64": newArchX8664,
    "386":   newArch386,
    "arm64": newArchAArch64,
    "arm":   newArchARM,
}

// GetRelType get relocation type
func GetRelType(code int) string {
    resolver, ok := name2Arch[runtime.GOARCH]
    if !ok {
        panic(fmt.Sprintf("Unsupported architecture: %q", runtime.GOARCH))
    }
    return resolver(code).String()
}

// archX8664  x86-64
type archX8664 struct{ Arch }

func newArchX8664(code int) RelResolver {
    return &archX8664{Arch{class: code}}
}

func (a archX8664) String() string {
    return elf.R_X86_64(a.class).String()
}

// arch386 386
type arch386 struct{ Arch }

func newArch386(code int) RelResolver {
    return &arch386{Arch{class: code}}
}

func (a arch386) String() string {
    return elf.R_386(a.class).String()
}

// archAArch64 AArch64
type archAArch64 struct{ Arch }

func newArchAArch64(code int) RelResolver {
    return &archAArch64{Arch{class: code}}
}

func (a archAArch64) String() string {
    return elf.R_AARCH64(a.class).String()
}

// archARM ARM
type archARM struct{ Arch }

func newArchARM(code int) RelResolver {
    return &archARM{Arch{class: code}}
}

func (a archARM) String() string {
    return elf.R_ARM(a.class).String()
}
