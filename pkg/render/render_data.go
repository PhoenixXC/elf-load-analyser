package render

import (
    "debug/elf"
    "errors"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
)

var dataCenter = make([]*Data, 3)

func PreAnalyse(ctx Content) {
    // env
    d, _ := doRender(NewEnvRender())
    dataCenter[0] = d
    // elf
    elfRender, e := NewElfRender(ctx.Filepath)
    if e != nil {
        var formatErr *elf.FormatError
        if ok := errors.As(e, &formatErr); ok {
            log.Errorf("Invalid elf file, %v", e)
        } else {
            log.Errorf("Analyse target binary format error, %v", e)
        }
    }
    d, _ = doRender(elfRender)
    dataCenter[1] = d
}

func DoAnalyse(p *data.Pool) []*Data {
    d, err := doRender(NewAnalyseRender(p.Data()))
    if err != nil {
        log.Errorf("Render analyse data error, %v", err)
    } else {
        dataCenter[2] = d
    }
    // NOTE use chain of responsibility
    return dataCenter
}
