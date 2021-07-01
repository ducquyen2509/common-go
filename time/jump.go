package time

import (
	"errors"
	"time"
)

func DayOffset(offset int) Options {
	return func(t *time.Time) error {
		*t = t.AddDate(0, 0, offset)
		return nil
	}
}

func MonthOffset(offset int) Options {
	return func(t *time.Time) error {
		*t = t.AddDate(0, offset, 0)
		return nil
	}
}

func YearOffset(offset int) Options {
	return func(t *time.Time) error {
		*t = t.AddDate(offset, 0, 0)
		return nil
	}
}

func TimeOffset(d time.Duration) Options {
	return func(t *time.Time) error {
		*t = t.Add(d)
		return nil
	}
}

func SetDay(d int) Options {
	return func(t *time.Time) error {
		*t = time.Date(t.Year(), t.Month(), d, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
		return nil
	}
}

func SetMonth(m time.Month) Options {
	return func(t *time.Time) error {
		*t = time.Date(t.Year(), m, t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
		return nil
	}
}

func SetYear(y int) Options {
	return func(t *time.Time) error {
		*t = time.Date(y, t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
		return nil
	}
}

func SetLocation(l *time.Location) Options {
	return func(t *time.Time) error {
		*t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), l)
		return nil
	}
}

func Offset(f Field, val int) Options {
	return func(t *time.Time) error {
		switch f {
		case Year:
			*t = t.AddDate(val, 0, 0)
		case Month:
			*t = t.AddDate(0, val, 0)
		case DayOfYear:
			*t = t.AddDate(0, 0, val)
		case DayOfMonth:
			maxDay := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location()).Add(-time.Nanosecond).Day()
			newDay := t.Day() + val
			for newDay < 0 {
				newDay += maxDay
			}
			newDay = newDay % maxDay
			if newDay == 0 {
				return SetDay(maxDay)(t)
			} else {
				return SetDay(newDay)(t)
			}
		default:
			return errors.New("field not found")
		}
		return nil
	}
}

func Set(f Field, val interface{}) Options {
	return func(t *time.Time) error {
		ok := false
		switch val.(type) {
		case int:
			ok = true
		case time.Month:
			if f == Month {
				ok = true
			}
			ok = false
		}
		if !ok {
			return errors.New("args not valid")
		}

		switch f {
		case Year:
			*t = time.Date(val.(int), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
		case Month:
			*t = time.Date(t.Year(), val.(time.Month), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
		case DayOfMonth:
			*t = time.Date(t.Year(), t.Month(), val.(int), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())

		case Hour:
			*t = time.Date(t.Year(), t.Month(), t.Day(), val.(int), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
		case HourOfDay:
			*t = time.Date(t.Year(), t.Month(), t.Day(), val.(int), t.Minute(), t.Second(), t.Nanosecond(), t.Location())

		case Minute:
			*t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), val.(int), t.Second(), t.Nanosecond(), t.Location())
		case Second:
			*t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), val.(int), t.Nanosecond(), t.Location())
		case Nanosecond:
			*t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), val.(int), t.Location())

		case Millisecond:
			return errors.New("field not support")
		case DayOfYear:
			return errors.New("field not support")
		default:
			return errors.New("field not found")
		}

		return nil
	}
}

func ResetTimeToZero() Options {
	return func(t *time.Time) error {
		//*t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		*t = t.Truncate(24 * time.Hour)
		return nil
	}
}
