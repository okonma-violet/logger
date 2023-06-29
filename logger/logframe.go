package logger

import (
	"encoding/binary"
	"time"
	"unsafe"
)

const (
	TagStartSep byte = 91 // "["
	TagEndSep   byte = 93 // "]"
	TagDelim    byte = 32 // " "
	Time_layout      = " [01/02 15:04:05.000000]"
)

var NetTagsDelim = [2]byte{255, 255}
var ByteOrder binary.ByteOrder = binary.LittleEndian

type netlogframe struct {
	lt        LogType
	localbody string
	netbody   []byte
}

func newNetLogFrame(lt LogType, t time.Time, loctags []byte, nettags []byte, tagsnum byte, lasttag, text string) netlogframe {
	return netlogframe{
		lt:        lt,
		localbody: concatLocalLogFrameBody(lt, t, loctags, lasttag, text),
		netbody:   concatNetLogFrameBody(lt, t, nettags, tagsnum, lasttag, text),
	}
}

func concatNetLogFrameBody(lt LogType, t time.Time, tags []byte, tagsnum byte, lasttag, text string) []byte {
	result := make([]byte, 14+len(tags)+len(lasttag)+len(text))
	result[0] = lt.Byte()
	ByteOrder.PutUint64(result[1:9], uint64(t.UnixMicro()))
	result[9] = tagsnum
	copy(result[10:], tags)
	copy(result[10+len(tags):], NetTagsDelim[:])
	copy(result[12+len(tags):], []byte(lasttag))
	copy(result[12+len(tags)+len(lasttag):], NetTagsDelim[:])
	copy(result[14+len(tags)+len(lasttag):], []byte(text))
	return result
}

func appendNetTags(tags []byte, toappend []string) []byte {
	n := (2 * len(toappend)) + len(tags)
	for _, a := range toappend {
		n += len(a)
	}
	result := make([]byte, len(tags), n)
	copy(result, tags)
	for _, a := range toappend {
		result = append(result, []byte(a)...)
		result = append(result, NetTagsDelim[:]...)
	}
	return result
}

//////////////////////////////////////////////////

type locallogframe struct {
	lt   LogType
	body string
}

func newLocalLogFrame(lt LogType, t time.Time, tags []byte, lasttag, text string) locallogframe {
	return locallogframe{lt: lt, body: concatLocalLogFrameBody(lt, t, tags, lasttag, text)}
}

func concatLocalLogFrameBody(lt LogType, t time.Time, tags []byte, lasttag, text string) string {
	result := make([]byte, 43+len(tags)+len(lasttag)+len(text))

	copy(result, []byte(lt.Colorize()))
	result[5] = TagStartSep
	copy(result[6:], lt.ByteStr())
	result[9] = TagEndSep
	copy(result[10:], []byte(ColorWhite))
	copy(result[15:], []byte(t.Format(Time_layout)))
	copy(result[39:], tags)
	result[39+len(tags)] = TagDelim
	result[40+len(tags)] = TagStartSep
	copy(result[41+len(tags):], []byte(lasttag))
	result[41+len(tags)+len(lasttag)] = TagEndSep
	result[42+len(tags)+len(lasttag)] = TagDelim
	copy(result[43+len(tags)+len(lasttag):], []byte(text))
	return *(*string)(unsafe.Pointer(&result))
}

func appendLocalTags(tags []byte, toappend []string) []byte {
	n := (3 * len(toappend)) + len(tags)
	for _, a := range toappend {
		n += len(a)
	}
	result := make([]byte, len(tags), n)
	copy(result, tags)
	for _, a := range toappend {
		result = append(result, TagDelim, TagStartSep)
		result = append(result, []byte(a)...)
		result = append(result, TagEndSep)
	}
	return result
}
