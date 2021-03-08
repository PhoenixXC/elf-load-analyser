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
    d, _ := NewEnvRender().Render()
    dataCenter[0] = d
    // elf
    d, e := NewElfRender(ctx.Filepath).Render()
    if e != nil {
        var formatErr *elf.FormatError
        if ok := errors.As(e, &formatErr); ok {
            log.Errorf("Invalid elf file, %v", e)
        } else {
            log.Errorf("Analyse target binary format error, %v", e)
        }
    }
    dataCenter[1] = d
}

func DoAnalyse(p *data.Pool) []*Data {
    render := NewAnalyseRender(p.Data())
    if d, err := render.Render(); err != nil {
        log.Errorf("Render analyse data error, %v", err)
    } else {
        dataCenter[2] = d
    }
    return dataCenter
}
