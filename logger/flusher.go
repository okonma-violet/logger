package logger

import (
	"time"
)

type LocalFlusher struct {
	ch         chan []locallogframe
	flushlvl   LogsFlushLevel
	cancel     chan struct{}
	allflushed chan struct{}
}

var flushertags []byte

const chanlen = 4

// TODO: make logsserver
func NewLocalFlusher(logsflushlvl LogsFlushLevel) LogsFlusher {
	f := &LocalFlusher{
		ch:         make(chan []locallogframe, chanlen),
		flushlvl:   logsflushlvl,
		cancel:     make(chan struct{}),
		allflushed: make(chan struct{}),
	}
	go f.flushWorker()
	return f
}

func (f *LocalFlusher) flushWorker() {
	for {
		select {
		case logslist := <-f.ch:
			for _, frame := range logslist {
				if LogsFlushLevel(frame.lt) >= f.flushlvl {
					println(frame.body)
				}
			}
		case <-f.cancel:
			for {
				select {
				case logslist := <-f.ch:
					for _, frame := range logslist {
						if LogsFlushLevel(frame.lt) >= f.flushlvl {
							println(frame.body)
						}
					}
				default:
					close(f.allflushed)
					return
				}
			}
		}
	}
}

func (f *LocalFlusher) Done() <-chan struct{} {
	return f.allflushed
}

func (f *LocalFlusher) DoneWithTimeout(timeout time.Duration) {
	t := time.NewTimer(timeout)
	select {
	case <-f.allflushed:
		return
	case <-t.C:
		PrintLn(Error, "Flusher", "bone by reached timeout, don't wait last flush")
		return
	}
}

func (f *LocalFlusher) Close() {
	close(f.cancel)
}

type NetFlusher struct {
	localflusher *LocalFlusher
	ch           chan []netlogframe
	//flushlvl   LogsFlushLevel
	cancel     chan struct{}
	allflushed chan struct{}

	servers []*nonEpollReConnector
}

func NewNetFlusher(local_logsflushlvl, net_logsflushlvl LogsFlushLevel) LogsFlusher {
	f := &NetFlusher{
		ch: make(chan []netlogframe, chanlen),
		//flushlvl:   logsflushlvl,
		cancel:     make(chan struct{}),
		allflushed: make(chan struct{}),
	}
	go f.flushWorker()
	return f
}
