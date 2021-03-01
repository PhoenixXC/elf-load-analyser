package data

import (
    "sort"
    "sync"
)

type Pool struct {
    mu    sync.Mutex
    order bool
    ch    chan *AnalyseData
    exit  chan struct{}
    data  []*AnalyseData
}

func (p *Pool) Len() int {
    return len(p.data)
}

func (p *Pool) Less(i, j int) bool {
    return p.data[i].Timestamp.Before(p.data[j].Timestamp)
}

func (p *Pool) Swap(i, j int) {
    d := p.data
    d[i], d[j] = d[j], d[i]
}

func NewPool() *Pool {
    return &Pool{ch: make(chan *AnalyseData), exit: make(chan struct{}), data: make([]*AnalyseData, 0), order: false}
}

func (p *Pool) Chan() chan<- *AnalyseData {
    return p.ch
}

func (p *Pool) Data() []*AnalyseData {
    <-p.exit
    if !p.order {
        p.mu.Lock()
        if !p.order {
            sort.Sort(p)
            p.order = true
        }
        p.mu.Unlock()
    }
    return p.data
}

func (p *Pool) Close() {
    close(p.ch)
    close(p.exit)
}

func (p *Pool) Init() {
    go func() {
    loop:
        for {
            select {
            case <-p.exit:
                break loop
            case data, ok := <-p.ch:
                if !ok {
                    break loop
                }
                p.data = append(p.data, data)
            }
        }
    }()
}
