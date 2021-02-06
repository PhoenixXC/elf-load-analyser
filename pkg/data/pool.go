package data

type Pool struct {
    ch   chan *AnalyseData
    exit chan struct{}
    data []*AnalyseData
}

func NewPool() *Pool {
    return &Pool{ch: make(chan *AnalyseData), exit: make(chan struct{}), data: make([]*AnalyseData, 0)}
}

func (p *Pool) Chan() chan<- *AnalyseData {
    return p.ch
}

func (p *Pool) Data() []*AnalyseData {
    <-p.exit
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
