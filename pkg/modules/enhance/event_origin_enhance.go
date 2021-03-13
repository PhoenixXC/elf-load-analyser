package enhance

import (
    "bytes"
    "fmt"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "github.com/phoenixxc/elf-load-analyser/pkg/modules"
)

const (
    eventOriginEnhancerName = "EventOrigin"

    eventOriginPrefix  = "origin_"
    EventOriginBccCode = eventOriginPrefix + "bcc_code"
    EventOriginMonitor = eventOriginPrefix + "monitor"
    EventOriginEvents  = eventOriginPrefix + "events"
)

func init() {
    modules.RegisteredEnhancer(eventOriginEnhancerName, &eventOriginEnhancer{})
}

// 返回 event 来源的 Monitor 信息
type eventOriginEnhancer struct{}

func (e eventOriginEnhancer) PreHandle(_ *modules.TableCtx) {}

func (e eventOriginEnhancer) AfterHandle(tCtx *modules.TableCtx,
    aData *data.AnalyseData, err error) (*data.AnalyseData, error) {
    if err != nil {
        return aData, err
    }
    var buffer bytes.Buffer
    events := tCtx.Monitor.Events()
    for i, event := range events {
        if i > 0 {
            buffer.WriteRune('@')
        }
        buffer.WriteString(fmt.Sprintf("%v|%s|%s", event.Class, event.FnName, event.Name))
    }
    serializeEvents := buffer.String()

    aData.PutExtra(EventOriginBccCode, tCtx.Monitor.Source())
    aData.PutExtra(EventOriginMonitor, tCtx.Monitor.Monitor())
    aData.PutExtra(EventOriginEvents, serializeEvents)

    log.Debugf("%s ==> %s collect monitor data", eventOriginEnhancerName, tCtx.Name)
    return aData, nil
}
