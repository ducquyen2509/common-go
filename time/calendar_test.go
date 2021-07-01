package time

import (
	"testing"
	"time"
)

func TestTimeUtil(t *testing.T) {
	type args struct {
		init        string
		yearOffset  int
		monthOffset int
		dayOffset   int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Case with 0",
			args: args{init: "2020-11-09", yearOffset: 0, monthOffset: 0, dayOffset: 0},
			want: "2020-11-09",
		},
		{
			name: "Test Case with negative month",
			args: args{init: "2020-11-09", yearOffset: 0, monthOffset: -11, dayOffset: 0},
			want: "2019-12-09",
		},
		{
			name: "Test Case with all negative numbers",
			args: args{init: "2020-11-09", yearOffset: -10, monthOffset: -11, dayOffset: -2},
			want: "2009-12-07",
		},
		{
			name: "Test Case with all positive numbers",
			args: args{init: "2020-11-09", yearOffset: 2, monthOffset: 0, dayOffset: 22},
			want: "2022-12-01",
		},
		{
			name: "Test Case with mix numbers",
			args: args{init: "2020-11-09", yearOffset: 2, monthOffset: -3, dayOffset: 22},
			want: "2022-08-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := GetTime(InitValue(DateFormat, tt.args.init),
				YearOffset(tt.args.yearOffset),
				MonthOffset(tt.args.monthOffset),
				DayOffset(tt.args.dayOffset)); err != nil {

				t.Error(err)
			} else if got.Format(DateFormat) != tt.want {
				t.Errorf("TimeUtilitis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeUtilTimeOffset(t *testing.T) {
	type args struct {
		init     string
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Case with 0",
			args: args{init: "2020-11-09 12:00:00", duration: 0 * time.Hour},
			want: "2020-11-09 12:00:00",
		},
		{
			name: "Test Case with 2 hour",
			args: args{init: "2020-11-09 12:00:00", duration: 2 * time.Hour},
			want: "2020-11-09 14:00:00",
		},
		{
			name: "Test Case with negative time",
			args: args{init: "2020-11-09 12:00:00", duration: -2 * time.Hour},
			want: "2020-11-09 10:00:00",
		},
		{
			name: "Test Case with over second",
			args: args{init: "2020-11-09 12:00:00", duration: 90 * time.Second},
			want: "2020-11-09 12:01:30",
		},
		{
			name: "Test Case with over minute",
			args: args{init: "2020-11-09 12:00:00", duration: 150 * time.Minute},
			want: "2020-11-09 14:30:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := GetTime(InitValue(DateTime24hFormat, tt.args.init),
				TimeOffset(tt.args.duration)); err != nil {

				t.Error(err)
			} else if got.Format(DateTime24hFormat) != tt.want {
				t.Errorf("TimeUtilitis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeUtilSetMethod(t *testing.T) {
	type args struct {
		init  string
		year  int
		month time.Month
		day   int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Case with current date",
			args: args{init: "2020-11-09", year: 2020, month: 11, day: 9},
			want: "2020-11-09",
		},
		{
			name: "Test Case with over day of month",
			args: args{init: "2020-11-09", year: 2019, month: time.February, day: 32},
			want: "2019-03-04",
		},
		{
			name: "Test Case with over day of month leap year",
			args: args{init: "2020-11-09", year: 2020, month: time.February, day: 32},
			want: "2020-03-03",
		},
		{
			name: "Test Case with negative day of month",
			args: args{init: "2020-11-09", year: 2020, month: 11, day: -9},
			want: "2020-10-22",
		},
		{
			name: "Test Case with negative year",
			args: args{init: "2020-11-09", year: -2002, month: 11, day: 9},
			want: "-2002-11-09",
		},
		{
			name: "Test Case with mix numbers",
			args: args{init: "2020-10-09", year: 2020, month: -1, day: 9},
			want: "2019-11-09",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := GetTime(InitValue(DateFormat, tt.args.init),
				SetYear(tt.args.year),
				SetMonth(tt.args.month),
				SetDay(tt.args.day)); err != nil {
				t.Error(err)
			} else if got.Format(DateFormat) != tt.want {
				t.Errorf("TimeUtilitis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeUtilResetTime(t *testing.T) {
	type args struct {
		init string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Case with 0",
			args: args{init: "2020-11-09 00:00:00"},
			want: "2020-11-09 00:00:00",
		},
		{
			name: "Test Case with negative time",
			args: args{init: "2020-11-09 09:32:12"},
			want: "2020-11-09 00:00:00",
		},
		{
			name: "Test Case with negative time",
			args: args{init: "2020-11-09 14:32:12"},
			want: "2020-11-09 00:00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := GetTime(InitValue(DateTime24hFormat, tt.args.init), ResetTimeToZero()); err != nil {

				t.Error(err)
			} else if got.Format(DateTime24hFormat) != tt.want {
				t.Errorf("TimeUtilitis() = %v, want %v", got.Format(DateTime24hFormat), tt.want)
			}
		})
	}
}

func TestOffsetDayOfMonth(t *testing.T) {
	type args struct {
		init   string
		offset int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Case with month day",
			args: args{init: "2020-11-09 00:00:00", offset: 30},
			want: "2020-11-09 00:00:00",
		},
		{
			name: "Test Case with positive value",
			args: args{init: "2020-11-09 00:00:00", offset: 2},
			want: "2020-11-11 00:00:00",
		},
		{
			name: "Test Case with positive value over month last time",
			args: args{init: "2020-11-09 00:00:00", offset: 35},
			want: "2020-11-14 00:00:00",
		},
		{
			name: "Test Case with negative time",
			args: args{init: "2020-11-09 00:00:00", offset: -1},
			want: "2020-11-08 00:00:00",
		},
		{
			name: "Test Case with negative time v2",
			args: args{init: "2020-11-09 00:00:00", offset: -15},
			want: "2020-11-24 00:00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := GetTime(InitValue(DateTime24hFormat, tt.args.init), Offset(DayOfMonth, tt.args.offset)); err != nil {

				t.Error(err)
			} else if got.Format(DateTime24hFormat) != tt.want {
				t.Errorf("TimeUtilitis() = %v, want %v", got.Format(DateTime24hFormat), tt.want)
			}
		})
	}

}

func TestTaoLao(t *testing.T) {
	//loc, err := time.LoadLocation("Local")
	//loc, err := time.LoadLocation("America/New_York")
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	ti := time.Now()
	t.Logf("%v", time.Date(ti.Year(), ti.Month(), ti.Day(), 0, 0, 0, 0, ti.Location()))
}
