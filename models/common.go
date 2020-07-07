package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type RetValue struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

type Time time.Time

const (
	timeFormat = "2006-01-02 15:04:05"
	DateFormat = "2006-01-02"

	Amount_Use   = "use"
	Amount_Conv  = "convert"
	Amount_Buy   = "buy"
	Amount_Delay = "delay"
)

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return
}
func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}
func (t Time) String() string {
	return time.Time(t).Format(timeFormat)
}

// Value insert timestamp into mysql need this function.
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if time.Time(t).UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return time.Time(t), nil
}

// Scan valueof time.Time
func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
