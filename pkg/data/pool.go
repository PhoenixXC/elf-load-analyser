package data

type Pool struct {
    ch   chan AnalyseData
    exit chan struct{}
    data []AnalyseData
}

func NewPool() *Pool {
    return &Pool{ch: make(chan AnalyseData), exit: make(chan struct{}), data: make([]AnalyseData, 0)}
}

func (p *Pool) Chan() chan<- AnalyseData {
    return p.ch
}

func (p *Pool) Data() []AnalyseData {
    // 当收集进程中止才可以获取数据，提前获取意义不大
    // TODO 实时获取数据，添加分页支持，需要注意线程安全
    <-p.exit
    return p.data
}

func (p *Pool) Close() {
    close(p.ch)
    close(p.exit)
}

func (p *Pool) Init() {
    go func() {
        for {
            select {
            case <-p.exit:
                break
            case data := <-p.ch:
                p.data = append(p.data, data)
            }
        }
    }()
}
