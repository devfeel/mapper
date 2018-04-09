package mapper

import "time"

type JSONTime time.Time

var (
	timeJSONFormat = "2006-01-02 15:04:05"
)

func SetTimeJSONFormat(format string) {
	timeJSONFormat = format
}

func GetTimeJSONFormat() string {
	return timeJSONFormat
}

func (t *JSONTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeJSONFormat+`"`, string(data), time.Local)
	*t = JSONTime(now)
	return
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeJSONFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeJSONFormat)
	b = append(b, '"')
	return b, nil
}

func (t JSONTime) String() string {
	return time.Time(t).Format(timeJSONFormat)
}
