package domain

import "fmt"

type MonthYearObject struct {
	month int
	year  int
}

func NewMonthYear(month int, year int) (MonthYearObject, error) {
	if month < 1 || month > 12 {
		return MonthYearObject{}, fmt.Errorf("invalid month: %d", month)
	}
	if year < 1 {
		return MonthYearObject{}, fmt.Errorf("invalid year: %d", year)
	}
	return MonthYearObject{
		month: month,
		year:  year,
	}, nil
}

func (m MonthYearObject) isCurrentMonth() bool {
	return false
}

func (m MonthYearObject) isPast() bool {
	return false
}

func (m MonthYearObject) isFuture() bool {
	return false
}

func (m MonthYearObject) Month() int {
	return m.month
}

func (m MonthYearObject) Year() int {
	return m.year
}
