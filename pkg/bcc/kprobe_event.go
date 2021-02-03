package bcc

import (
    bpf "github.com/iovisor/gobpf/bcc"
)

type KprobeEvent struct {
    *Event
    maxActive int
}

func NewKprobeEvent(name string, fnName string, maxActive int) *Event {
    e := NewEvent(KprobesEvent, name, fnName)
    ke := KprobeEvent{
        Event:     e,
        maxActive: maxActive,
    }
    e.action = &ke
    return e
}

func (e KprobeEvent) Attach(m *bpf.Module, fd int) error {
    return m.AttachKprobe(e.FnName, fd, e.maxActive)
}

func (e *KprobeEvent) Load(m *bpf.Module) (int, error) {
    return m.LoadKprobe(e.Name)
}