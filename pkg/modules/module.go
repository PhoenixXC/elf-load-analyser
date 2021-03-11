package modules

import (
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
    bpf "github.com/iovisor/gobpf/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/factory"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "reflect"
    "strconv"
)

var (
    defaultFlags []string
)

const (
    EnhanceTag = "enhance"
)

// EventResult 事件结果接口，实现中的类型大小在编译时必须是已知的
type EventResult interface {
    Render() *data.AnalyseData
}

// MonitorModule 模块抽象接口
type MonitorModule interface {
    // Monitor 返回模块的名称，以及是否作为结束标志
    Monitor() string
    // Source 返回注入的 bcc 源码
    Source() string
    // Events 返回要注册的事件
    Events() []*bcc.Event
    // Resolve 解析、发送处理结果
    Resolve(m *bpf.Module, ch chan<- *data.AnalyseData, ready chan<- struct{}, stop <-chan struct{})
}

// ModuleDefaultInit 注册 Module
func ModuleDefaultInit(mm MonitorModule) {
    ModuleInit(mm, false)
}

// ModuleInit 注册 Module
func ModuleInit(mm MonitorModule, end bool) {
    // PerfResolveMonitorModule stop handler check
    if end {
        if mm, ok := mm.(*PerfResolveMonitorModule); ok {
            mm.stopHandler = nil
        }
    }

    m := bcc.NewMonitor(mm.Monitor(), mm.Source(), defaultFlags, mm.Resolve)
    for _, event := range mm.Events() {
        m.AddEvent(event)
    }
    if end {
        m.SetEnd()
    }
    factory.Register(m)
}

func Render(d []byte, event EventResult, enhance bool) (*data.AnalyseData, error) {
    err := binary.Read(bytes.NewBuffer(d), bpf.GetHostByteOrder(), event)
    if err != nil {
        return nil, fmt.Errorf("Failed to decode received data to %q, %v\n",
            reflect.TypeOf(event).Name(), err)
    }
    aData := event.Render()
    if enhance {
        enhanceStructField(event, aData)
    }
    return aData, nil
}

func enhanceStructField(ptr interface{}, d *data.AnalyseData) {
    reType := reflect.TypeOf(ptr)
    var v = reflect.ValueOf(ptr)
    if reType.Kind() == reflect.Ptr {
        reType = reflect.ValueOf(ptr).Elem().Type()
        v = reflect.ValueOf(ptr).Elem()
    }
    if reType.Kind() != reflect.Struct {
        panic("Invalid type" + reType.Kind().String())
    }
    parseField(v, d)
}

func parseField(v reflect.Value, d *data.AnalyseData) {
    for i := 0; i < v.NumField(); i++ {
        structField := v.Type().Field(i)
        k := v.Field(i)
        if !k.CanInterface() {
            continue
        }
        tag := structField.Tag
        label := tag.Get(EnhanceTag)
        if label == "" {
            if structField.Anonymous && structField.Type.Kind() == reflect.Struct {
                parseField(k, d)
            }
            continue
        }
        s, err := toString(k)
        if err == nil {
            d.PutExtra(label, s)
        } else {
            log.Warnf("Parse field(%q) error: %v", k.Type().Name, err)
        }
    }
}

func toString(value reflect.Value) (key string, err error) {
    switch value.Type().Kind() {
    case reflect.Float64:
        key = strconv.FormatFloat(value.Float(), 'f', -1, 64)
    case reflect.Float32:
        key = strconv.FormatFloat(value.Float(), 'f', -1, 64)
    case reflect.Int:
        key = strconv.Itoa(int(value.Int()))
    case reflect.Uint:
        key = strconv.FormatUint(value.Uint(), 10)
    case reflect.Int8:
        key = strconv.Itoa(int(value.Int()))
    case reflect.Uint8:
        key = strconv.Itoa(int(value.Uint()))
    case reflect.Int16:
        key = strconv.Itoa(int(value.Int()))
    case reflect.Uint16:
        key = strconv.Itoa(int(value.Uint()))
    case reflect.Int32:
        key = strconv.Itoa(int(value.Int()))
    case reflect.Uint32:
        key = strconv.Itoa(int(value.Uint()))
    case reflect.Int64:
        key = strconv.FormatInt(value.Int(), 10)
    case reflect.Uint64:
        key = strconv.FormatUint(value.Uint(), 10)
    case reflect.String:
        key = value.String()
    case reflect.Interface:
        i := value.Interface()
        if si, ok := i.(fmt.Stringer); ok {
            key = si.String()
        } else {
            err = errors.New(fmt.Sprintf("Unsupported interface type: %q\n", value.Type().Name()))
        }
    default:
        err = errors.New(fmt.Sprintf("Unsupported type %q\n", value.Type().Name()))
    }

    return
}
