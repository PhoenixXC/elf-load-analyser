package bcc

import bpf "github.com/iovisor/gobpf/bcc"

type KretprobeEvent struct {
    KprobeEvent
}

func NewKretprobeEvent(name string, fnName string, maxActive int) *Event {
    e := NewEvent(KretprobesEvent, name, fnName)
    ke := KretprobeEvent{
        KprobeEvent: KprobeEvent{
            Event:     e,
            maxActive: maxActive,
        },
    }
    e.action = &ke
    return e
}

func (e KretprobeEvent) Attach(m *bpf.Module, fd int) error {
    return m.AttachKretprobe(e.FnName, fd, e.maxActive)
}