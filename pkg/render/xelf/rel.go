package xelf

import (
    "debug/elf"
    "encoding/binary"
    "errors"
    "fmt"
    "io"
    "strings"
    "unsafe"
)

type word uint8

const (
    X86 = word(iota)
    X64
)

type wrapRelocation struct {
    Off       uint64
    Info      uint64
    Addend    int64
    w         word
    IsRela    bool
    IsDynamic bool
}

func newWrapRelocation(is64 bool, d []uint64) (*wrapRelocation, error) {
    if d == nil || len(d) < 2 {
        return nil, fmt.Errorf("[]uint64 data invalid, can't parse")
    }
    wRel := &wrapRelocation{w: X64, IsRela: true}
    if !is64 {
        wRel.w = X86
    }
    wRel.Off = d[0]
    wRel.Info = d[1]
    if len(d) == 2 {
        wRel.IsRela = false
    } else {
        wRel.Addend = int64(d[2])
    }
    return wRel, nil
}

type RelSection struct {
    Section *elf.SectionHeader
    Rels    []RelDecoded
}

type RelDecoded struct {
    XType     string
    Value     string
    Offset    uint64
    Addend    int64
    W         word
    IsAddend  bool
    IsDynamic bool
}

func newRel(f *elf.File, wr wrapRelocation) (*RelDecoded, error) {
    rel := &RelDecoded{
        W:         wr.w,
        Offset:    wr.Off,
        Addend:    wr.Addend,
        IsAddend:  wr.IsRela,
        IsDynamic: wr.IsDynamic,
    }
    var sym, typ uint32
    switch wr.w {
    case X86:
        sym = elf.R_SYM32(uint32(wr.Info))
        typ = elf.R_TYPE32(uint32(wr.Info))
    case X64:
        sym = elf.R_SYM64(wr.Info)
        typ = elf.R_TYPE64(wr.Info)
    default:
        panic("BUG: unknown word size")
    }
    symbol, err := getSymbol(f, sym, wr.IsDynamic)
    if err != nil {
        return nil, err
    }
    rel.XType = GetRelType(int(typ))
    rel.Value = symbol
    return rel, nil
}

func BuildRelIf(f *elf.File, dynamic bool) ([]RelSection, error) {
    if f.Class == elf.ELFCLASSNONE {
        return nil, fmt.Errorf("uknown xelf class")
    }
    var srs []RelSection //nolint:prealloc

    sections := f.Sections
    for i := range sections {
        t := sections[i].Type
        if t != elf.SHT_REL && t != elf.SHT_RELA {
            continue
        }
        sec := sections[i]

        if dynamic != isDynSec(sec.Name) {
            continue
        }

        sr := RelSection{Section: &sec.SectionHeader, Rels: []RelDecoded{}}
        data, err := readBytes(sec.Open(), int(sec.FileSize))
        if err != nil {
            return nil, err
        }

        rels := make([]wrapRelocation, sec.FileSize/sec.Entsize)
        if err := fillRelIf(data, rels, dynamic); err != nil {
            return nil, err
        }

        for j := range rels {
            rel, err := newRel(f, rels[j])
            if err != nil {
                return nil, err
            }
            sr.Rels = append(sr.Rels, *rel)
        }
        srs = append(srs, sr)
    }
    return srs, nil
}

func isDynSec(name string) bool {
    return strings.HasSuffix(name, "dyn") || strings.HasSuffix(name, "plt")
}

// fillRelIf 从字节流中填充对象
func fillRelIf(data []byte, ifs []wrapRelocation, dynamic bool) error {
    itemNum := len(data) / len(ifs)
    for i := range ifs {
        buf := data[i*itemNum : (i+1)*itemNum]
        obj, err := parseBytes(buf)
        if err != nil {
            return err
        }
        obj.IsDynamic = dynamic
        ifs[i] = *obj
    }
    return nil
}

// 将字节流解析为对象
func parseBytes(b []byte) (*wrapRelocation, error) {
    bLen := len(b)
    // 2 * 32 3 * 32 | 4 * 32 6 * 32
    itemSize, is64 := 8, true
    if bLen/4 < 4 {
        is64 = false
        itemSize = 4
    }
    byteOrder := getHostByteOrder()
    var uintsIdx = 0
    var uints = make([]uint64, bLen/itemSize)
    for idx := 0; idx+itemSize < bLen; idx += itemSize {
        var u uint64
        endIdx := idx + itemSize
        if is64 {
            u = byteOrder.Uint64(b[idx:endIdx])
        } else {
            u = uint64(byteOrder.Uint32(b[idx:endIdx]))
        }
        uints[uintsIdx] = u
        uintsIdx++
    }
    return newWrapRelocation(is64, uints)
}

func readBytes(r io.Reader, limit int) ([]byte, error) {
    var objBuf []byte
    buf := make([]byte, limit)
    for {
        l, err := r.Read(buf)
        if err != nil {
            if errors.Is(err, io.EOF) {
                objBuf = append(objBuf, buf[:l]...)
                break
            }
            return nil, err
        }
        objBuf = append(objBuf, buf[:l]...)
    }
    return objBuf, nil
}

func getHostByteOrder() binary.ByteOrder {
    var i int32 = 0x01020304
    u := unsafe.Pointer(&i)
    pb := (*byte)(u)
    b := *pb
    if b == 0x04 {
        return binary.LittleEndian
    }
    return binary.BigEndian
}

func getSymbol(f *elf.File, sym uint32, dynamic bool) (string, error) {
    // undefined
    if sym == 0 {
        return "", nil
    }
    var err error
    var symbols []elf.Symbol
    if dynamic {
        symbols, err = f.DynamicSymbols()
        if err != nil {
            return "", nil
        }
    } else {
        symbols, err = f.Symbols()
        if err != nil {
            return "", nil
        }
    }
    if int(sym) > len(symbols) {
        panic("BUG: rel sym idx out of Symbols")
    }
    return symbols[sym-1].Name, nil
}
