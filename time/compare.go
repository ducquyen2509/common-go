package time

import "time"

// CompareTimeFormat compare 2 time with same layout
// return
//	1 if t1 > t2
//	0 if t1 = t2
//	-1 if t1 < t2
func CompareTimeFormat(layout, t1, t2 string) (int, error) {
	time1, err := GetTime(InitValue(layout, t1))
	if err != nil {
		return 0, err
	}

	time2, err := GetTime(InitValue(layout, t2))
	if err != nil {
		return 0, err
	}

	return CompareTime(time1, time2), nil
}

func CompareTime(t1, t2 time.Time) int {
	if t1.Sub(t2) > 0 {
		return 1
	} else if t1.Sub(t2) < 0 {
		return -1
	} else {
		return 0
	}
}
