package enhance

import (
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "github.com/phoenixxc/elf-load-analyser/pkg/modules"
    "strconv"
)

const (
    StartMark        = "_START_"
    kernelBootNsKey  = "_ns_"
    timeEnhancerName = "TimeAmend"
)

func init() {
    modules.RegisteredEnhancer(timeEnhancerName, &timeEnhancer{})
}

type TimeEvent interface {
    Ns() uint64
}

// TimeEventResult EventResult with env boot ns, coordinate with enhance.timeEnhancer
type TimeEventResult struct {
    TS uint64 `enhance:"_ns_"`
}

func (t *TimeEventResult) Ns() uint64 {
    return t.TS
}

// timeEnhancer 时间修正处理器
type timeEnhancer struct{}

func (t timeEnhancer) PreHandle(tCtx *modules.TableCtx) {
    if !tCtx.IsMark(StartMark) {
        WaitNs()
    }
}

func (t timeEnhancer) AfterHandle(tCtx *modules.TableCtx,
    aData *data.AnalyseData, err error) (*data.AnalyseData, error) {
    if err != nil {
        return aData, err
    }

    if v, ok := aData.ExtraByKey(kernelBootNsKey); ok {
        ns, e := strconv.ParseUint(v, 10, 64)
        if e != nil {
            log.Warnf("%s ==> can not convert %q data, %v", timeEnhancerName, kernelBootNsKey, e)
            return aData, err
        }
        if tCtx.IsMark(StartMark) {
            log.Debugf("%s ==> %s notify other coordinate", timeEnhancerName, tCtx.Name)
            SendNs(NewNsMap(ns))
        } else {
            log.Debugf("%s ==> %s amend timestamp", timeEnhancerName, tCtx.Name)
            aData.Timestamp = AmendTime(ns)
        }
        aData.RmExtra(kernelBootNsKey)
    }
    return aData, err
}
