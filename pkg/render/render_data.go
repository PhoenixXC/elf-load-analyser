package render

import (
    "debug/elf"
    "encoding/json"
    "errors"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/env"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
)

var dataCenter = make([]*data.AnalyseData, 3)

func PreAnalyse(ctx Content) {
    // pre check
    env.CheckEnv()
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

// VisualAnalyseData 数据展示，若 show 为 true，开启 web 服务展示数据，否则持久化数据到硬盘上
func VisualAnalyseData(p *data.Pool, show bool)  {
    renderedData := doAnalyse(p)
    for _, analyseData := range renderedData {
        d, _ := json.Marshal(analyseData)
        log.Info(string(d))
    }
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