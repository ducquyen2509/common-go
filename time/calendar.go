package time

import "time"

const (
	DateFormat              = "2006-01-02"
	DateWithoutHyphenFormat = "20060102"
	DateTime12hFormat       = "2006-01-02 03:04:05"
	DateTime24hFormat       = "2006-01-02 15:04:05"

	Year       = Field(1)
	Month      = Field(2)
	DayOfMonth = Field(5)
	DayOfYear  = Field(6)

	Hour        = Field(10)
	HourOfDay   = Field(11)
	Minute      = Field(12)
	Second      = Field(13)
	Millisecond = Field(14)
	Nanosecond  = Field(15)
)

type Options func(t *time.Time) error
type Field int

func TimestampMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetTime(ops ...Options) (time.Time, error) {

	t := time.Now()
	for _, op := range ops {
		if err := op(&t); err != nil {
			return t, err
		}
	}
	return t, nil
}
