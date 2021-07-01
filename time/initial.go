package time

import "time"

func FromMillisecond(m int64) Options {
	return func(t *time.Time) error {
		*t = time.Unix(0, m*int64(time.Millisecond))
		return nil
	}
}

func InitValue(layout, val string) Options {
	return func(t *time.Time) error {
		if tmp, err := time.Parse(layout, val); err != nil {
			return err
		} else {
			*t = tmp
			return nil
		}
	}
}

func CopyOf(sour time.Time) Options {
	return func(t *time.Time) error {
		*t = sour
		return nil
	}
}
