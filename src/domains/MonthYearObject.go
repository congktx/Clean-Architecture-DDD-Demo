package domains

type MonthYearObject struct {
	month int
	year  int
}

func NewMonthYear(month int, year int) MonthYearObject {
	return MonthYearObject{
		month: month,
		year:  year,
	}
}

func (m *MonthYearObject) isCurrentMonth() bool {
	return false
}

func (m *MonthYearObject) isPast() bool {
	return false
}

func (m *MonthYearObject) isFuture() bool {
	return false
}
