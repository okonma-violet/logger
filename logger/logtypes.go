package logger

type LogType byte

const (
	Debug   LogType = 1
	Info    LogType = 2
	Warning LogType = 3
	Error   LogType = 4
)

var (
	debugStr   = append(make([]byte, 3), "DBG"...)
	infoStr    = append(make([]byte, 3), "INF"...)
	warningStr = append(make([]byte, 3), "WRN"...)
	errorStr   = append(make([]byte, 3), "ERR"...)
	unknownStr = append(make([]byte, 3), "UNK"...)
)

func (lt LogType) String() string {
	switch lt {
	case Debug:
		return "DBG"
	case Info:
		return "INF"
	case Warning:
		return "WRN"
	case Error:
		return "ERR"
	}
	return "UNK"
}

func (lt LogType) Byte() byte {
	return byte(lt)
}

// do not change result array, only copy it
func (lt LogType) ByteStr() []byte {
	switch lt {
	case Debug:
		return debugStr
	case Info:
		return infoStr
	case Warning:
		return warningStr
	case Error:
		return errorStr
	}
	return unknownStr
}

const (
	ColorReset = "\033[0m"
	ColorWhite = "\033[97m"
)

func (lt LogType) Colorize() string {
	switch lt {
	case Debug:
		return "\033[32m" // Green
	case Info:
		return "\033[36m" // Gay
	case Warning:
		return "\033[33m" // Yellow
	case Error:
		return "\033[31m" // Red
	}
	return "\033[35m" // Purple
}

type LogsFlushLevel byte

const (
	ZeroLevel    LogsFlushLevel = 0 // none
	DebugLevel   LogsFlushLevel = 1 // all
	InfoLevel    LogsFlushLevel = 2 // infos, warnings, errors only
	WarningLevel LogsFlushLevel = 3 // warnings, errors only
	ErrorLevel   LogsFlushLevel = 4 // errors only
)
