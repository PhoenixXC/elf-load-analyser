package render

import (
    "debug/elf"
    "fmt"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/data/markdown"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "github.com/phoenixxc/elf-load-analyser/pkg/render/xelf"
    "os"
    "strconv"
)

type ElfRender struct {
    filepath string
    f        *elf.File
}

func NewElfRender(filepath string) (*ElfRender, error) {
    // file header
    f, err := os.Open(filepath)
    if err != nil {
        return nil, err
    }
    eFile, err := elf.NewFile(f)
    if err != nil {
        return nil, err
    }
    return &ElfRender{filepath: filepath, f: eFile}, nil
}

func (e *ElfRender) Render() (d *Data, err error) {
    fHeader := markdown.NewTitleContents(markdown.H2, "ELF File Header").Append(e.buildHeader())
    fProgHeader := markdown.NewTitleContents(markdown.H2, "ELF Prog Header").Append(e.buildProgHeader())
    content := fHeader.Append(fProgHeader).
        // exec file no static rel section
        // Append(e.buildStaticData()).
        Append(e.buildDynamicData())

    t := e.Type()
    d = NewData(data.NewAnalyseData(t.Name, content).WithID(t.ID))
    return
}

func (e *ElfRender) Type() Type {
    return ElfType
}

func (e *ElfRender) Release() {
    _ = e.f.Close()
}

func (e *ElfRender) buildHeader() markdown.Interface {
    header := e.f.FileHeader

    table := markdown.NewTable("MEMBER", "VALUE").
        WithDesc(fmt.Sprintf("table 1: file %q header, for more information, see: %q", e.filepath, "readelf -h ..."))
    table.AddRow("Class", header.Class.String()).
        AddRow("data", header.Data.String()).
        AddRow("ByteOrder", header.ByteOrder.String()).
        AddRow("Version", header.Version.String()).
        AddRow("Os/ABI", header.OSABI.String()).
        AddRow("ABI Version", strconv.Itoa(int(header.ABIVersion))).
        AddRow("Type", header.Type.String()).
        AddRow("Machine", header.Machine.String()).
        AddRow("Version", header.Version.String()).
        AddRow("Entry", convertAddr(header.Entry))
    return table
}

func (e *ElfRender) buildProgHeader() markdown.Interface {
    ph := e.f.Progs

    table := markdown.NewTable("Type", "Offset", "FileSize", "VirtAddr", "MemSize", "PhysAddr", "Flags", "Align").
        WithDesc(fmt.Sprintf("table 2: file %q program headers, for more information, see: %q", e.filepath, "readelf -l ..."))
    for _, prog := range ph {
        row := make([]string, table.Col())
        convertRow(row, prog.Type, prog.Off, prog.Filesz, prog.Vaddr, prog.Memsz, prog.Paddr, prog.Flags, prog.Align)
        table.AddRow(row...)
    }

    return table
}

func (e *ElfRender) buildStaticData() markdown.Interface {
    sectionRels, err := xelf.BuildRelIf(e.f, false)
    if err != nil {
        log.Warnf("Get static rel data from elf file failed: %v", err)
        return markdown.EmptyIf
    }
    if len(sectionRels) == 0 {
        return markdown.EmptyIf
    }

    mk := markdown.NewTitleContents(markdown.H3, "Static relocation info")
    for _, rel := range sectionRels {
        mk.Append(relSecToMarkdown(rel))
    }
    return mk
}

func (e *ElfRender) buildDynamicData() markdown.Interface {
    dynInfo, err := xelf.BuildDynamicInfo(e.f)
    if err != nil {
        log.Warnf("Get static rel data from elf file failed: %v", err)
        return markdown.EmptyIf
    }

    mk := markdown.NewTitleContents(markdown.H2, "Dynamic info").
        Append(markdown.NewTitleContents(markdown.H3, "interp").WithContents(dynInfo.Interp))

    symTableContent := markdown.NewTitleContents(markdown.H3, "dynamic symbols")
    if syms := dynInfo.Symbols; len(syms) == 0 {
        symTableContent.WithContents("no data")
    } else {
        symTable := markdown.NewTable("Name", "Section", "Value", "Size", "Library", "Version")
        for _, symbol := range syms {
            row := make([]string, symTable.Col())
            convertRow(row, symbol.Name, symbol.Section, symbol.Value, symbol.Size, symbol.Library, symbol.Version)
            symTable.AddRow(row...)
        }
        symTableContent.Append(symTable)
    }

    tag2DynContent := markdown.NewTitleContents(markdown.H3, "dyn string")
    if t2d := dynInfo.Tag2Dyn; len(dynInfo.Tag2Dyn) == 0 {
        tag2DynContent.WithContents("No data")
    } else {
        tag2DynTable := markdown.NewTable("Tag", "String")
        for tag, strList := range t2d {
            if len(strList) == 0 {
                continue
            }
            tagStr := tag.String()
            for _, str := range strList {
                tag2DynTable.AddRow(tagStr, str)
            }
        }
        tag2DynContent.Append(tag2DynTable)
    }

    importSymsContent := markdown.NewTitleContents(markdown.H3, "import symbols")
    if iSym := dynInfo.ImportedSymbols; len(iSym) == 0 {
        importSymsContent.WithContents("no data")
    } else {
        importSymsTable := markdown.NewTable("Name", "Version", "Library")
        for _, symbol := range iSym {
            importSymsTable.AddRow(symbol.Name, symbol.Version, symbol.Library)
        }
        importSymsContent.Append(importSymsTable)
    }

    relSecs := markdown.NewTitleContents(markdown.H3, "Dynamic relocation sections")
    if rss := dynInfo.RelSections; len(rss) == 0 {
        relSecs.WithContents("No data")
    } else {
        for _, rel := range rss {
            relSecs.Append(relSecToMarkdown(rel))
        }
    }

    mk.Append(symTableContent)
    mk.Append(tag2DynContent)
    mk.Append(importSymsContent)
    mk.Append(relSecs)
    return mk
}

func relSecToMarkdown(rSec xelf.RelSection) markdown.Interface {
    sec := rSec.Section
    mk := markdown.NewTitleContents(markdown.H4, "Relocation Section ["+sec.Name+"]")
    secTable := markdown.NewTable("Name", "Type", "Flags", "Addr", "Offset",
        "Size", "Link", "Info", "Addralign", "Entsize", "FileSize")
    row := make([]string, secTable.Col())
    convertRow(row, sec.Name, sec.Type, sec.Flags, sec.Addr, sec.Offset, sec.Size, sec.Link, sec.Info,
        sec.Addralign, sec.Entsize, sec.FileSize)
    secTable.AddRow(row...)
    mk.Append(secTable).
        Append(relsToMarkdown(rSec.Rels))
    return mk
}

func relsToMarkdown(rels []xelf.RelDecoded) markdown.Interface {
    if len(rels) == 0 {
        return markdown.EmptyIf
    }
    relTable := markdown.NewTable("Offset", "Type", "Value")
    for _, rel := range rels {
        relTable.AddRow(convertAddr(rel.Offset), rel.XType, rel.Value)
    }
    return relTable
}

func convertRow(row []string, d ...interface{}) {
    for i := range d {
        val := d[i]
        switch val := val.(type) {
        case uint64:
            row[i] = convertAddr(val)
        case uint32:
            row[i] = convertAddr(uint64(val))
        case string:
            row[i] = val
        case fmt.Stringer:
            row[i] = val.String()
        }
    }
}

func convertAddr(addr uint64) string {
    return "0X" + strconv.FormatUint(addr, 16)
}
