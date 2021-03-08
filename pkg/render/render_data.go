package render

import (
    "debug/elf"
    "errors"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "github.com/phoenixxc/elf-load-analyser/pkg/web"
)

var dataCenter = make([]*data.AnalyseData, 3)

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

// VisualAnalyseData 数据展示
func VisualAnalyseData(p *data.Pool, port uint) {
    renderedData := doAnalyse(p)
    go web.StartWebService(renderedData, port)
}

func doAnalyse(p *data.Pool) []*data.AnalyseData {
    render := NewAnalyseRender(p.Data())
    if d, err := render.Render(); err != nil {
        log.Errorf("Render analyse data error, %v", err)
    } else {
        dataCenter[2] = d
    }
    return dataCenter
}