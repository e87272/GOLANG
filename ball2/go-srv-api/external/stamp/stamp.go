package stamp

import "time"

const (
	Second   = int64(1000)
	Minute   = Second * 60
	Hour     = Minute * 60
	Day      = Hour * 24
	timeZone = Hour * +8
)

func Now() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func Date(timestamp int64) int64 {
	return timestamp - (timestamp+timeZone)%Day
}

func Today() int64 {
	return Date(Now())
}
