package logger

import (
	"time"
)

type LocalLogsContainer struct {
	ch   chan []locallogframe
	tags []string
}

func (f *LocalFlusher) NewLogsContainer(tags ...string) Logger {
	return &LocalLogsContainer{ch: f.ch, tags: tags}
}

func (l *LocalLogsContainer) Debug(name, logstr string) {
	l.ch <- []locallogframe{newLocalLogFrame()}
}

func (l *LocalLogsContainer) Info(name, logstr string) {
	l.ch <- []locallogframe{newLogFrame(Debug, time.Now(), l.tags, logstr)}
}

func (l *LocalLogsContainer) Warning(name, logstr string) {
	l.ch <- []locallogframe{newLogFrame(Debug, time.Now(), l.tags, logstr)}
}

func (l *LocalLogsContainer) Error(name string, logerr error) {
	var logstr string
	if logerr != nil {
		logstr = logerr.Error()
	} else {
		logstr = "nil err"
	}
	l.ch <- []logframe{newLogFrame(Debug, time.Now(), l.tags, logstr)}
}

func (l *LocalLogsContainer) ErrorStr(name string, logstr string) {
	if logstr == "" {
		logstr = "nil err"
	}
	l.ch <- []logframe{newLogFrame(Debug, time.Now(), l.tags, logstr)}
}

func (l *LocalLogsContainer) NewSubLogger(tags ...string) Logger {
	return &LocalLogsContainer{ch: l.ch, tags: append(l.tags, tags...)}
}

////////////////////////////////////////////////////////////////

type NetLogsContainer struct {
	ch        chan [][]byte
	tagscount int
	tags      []byte
}

func (f *NetFlusher) NewLogsContainer(tags ...string) Logger {
	return &NetLogsContainer{
		ch:        f.ch,
		tagscount: len(tags),
		tags:      concatToBytesWithDelim(tags),
	}
}

func (l *NetLogsContainer) Debug(name, logstr string) {
	l.ch <- [][]byte{EncodeLog(Debug, time.Now(), l.tags, name, logstr)}
}

func (l *NetLogsContainer) Info(name, logstr string) {
	l.ch <- [][]byte{EncodeLog(Info, time.Now(), l.tags, name, logstr)}
}

func (l *NetLogsContainer) Warning(name, logstr string) {
	l.ch <- [][]byte{EncodeLog(Warning, time.Now(), l.tags, name, logstr)}
}

func (l *NetLogsContainer) Error(name string, logerr error) {
	var logstr string
	if logerr != nil {
		logstr = logerr.Error()
	} else {
		logstr = "nil err"
	}
	l.ch <- [][]byte{EncodeLog(Error, time.Now(), l.tags, name, logstr)}
}

func (l *NetLogsContainer) NewSubLogger(tags ...string) Logger {
	return &LocalLogsContainer{ch: l.ch, tags: AppendTags(l.tags, tags...)}
}
