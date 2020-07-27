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
	TimeFormat = "2006-01-02 15:04:05"
	DateFormat = "2006-01-02"

	Amount_Use        = "use"
	Amount_Cancel     = "cancel"
	Amount_ConvIn     = "convert_in"
	Amount_ConvOut    = "convert_out"
	Amount_Frozen_In  = "frozen_in"
	Amount_Frozen_Out = "frozen_out"
	Amount_Buy        = "buy"
	Amount_Back       = "back"
	Amount_Delay_In   = "delay_in"
	Amount_Delay_Out  = "delay_out"
)

var AmountChange = map[string]int{
	Amount_Use:        -1,
	Amount_Cancel:     1,
	Amount_ConvIn:     1,
	Amount_ConvOut:    -1,
	Amount_Buy:        1,
	Amount_Delay_In:   1,
	Amount_Delay_Out:  -1,
	Amount_Frozen_In:  1,
	Amount_Frozen_Out: -1,
	Amount_Back:       1,
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return
}
func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}
func (t Time) String() string {
	return time.Time(t).Format(TimeFormat)
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
