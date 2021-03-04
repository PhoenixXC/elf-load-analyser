package enhance

import (
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "github.com/phoenixxc/elf-load-analyser/pkg/modules"
    "strconv"
)

const (
    StartMark       = "_START_"
    kernelBootNsKey = "_ns_"
    enhancerName    = "time_amend"
)

func init() {
    modules.RegisteredEnhancer(enhancerName, &timeEnhance{})
}

type TimeEvent interface {
    Ns() uint64
}

// TimeEventResult EventResult with env boot ns, coordinate with enhance.timeEnhance
type TimeEventResult struct {
    TS uint64 `enhance:"_ns_"`
}

func (t *TimeEventResult) Ns() uint64 {
    return t.TS
}

// timeEnhance 时间修正处理器
type timeEnhance struct{}

func (t timeEnhance) PreHandle(tCtx *modules.TableCtx) {
    if !tCtx.IsMark(StartMark) {
        WaitNs()
    }
}

func (t timeEnhance) AfterHandle(tCtx *modules.TableCtx,
    aData *data.AnalyseData, err error) (*data.AnalyseData, error) {
    if err != nil {
        return aData, err
    }

    if v, ok := aData.Extra(kernelBootNsKey); ok {
        ns, e := strconv.ParseUint(v, 10, 64)
        if e != nil {
            log.Warnf("%s can not convert %q data, %v", enhancerName, kernelBootNsKey, e)
            return aData, err
        }
        if tCtx.IsMark(StartMark) {
            SendNs(NewNsMap(ns))
        } else {
            aData.Timestamp = AmendTime(ns)
        }
    }
    return aData, err
}
