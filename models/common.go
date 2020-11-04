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
type Date time.Time

const (
	UserLeader  = 10
	UserHR      = 6
	UserIT      = 7
	UserFinance = 8
	UserFront   = 9

	TimeFormat = "2006-01-02 15:04:05"
	DateFormat = "2006-01-02"
	PosFormat = "2006/01/02 15:04:05"

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

	//flow status
	FlowNA         = "NA"
	FlowProcessing = "Processing"
	FlowCompleted  = "Completed"
	FlowApproved   = "Approved"
	FlowRejected   = "Rejected"

	LeaveSick   = "Sick"
	LeaveShift  = "Shift"
	LeaveAnnual = "Annual"
)

var Months = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

var EmpStatus = map[int]string{
	0: "未入职",
	1: "拟入职",
	2: "在职",
	3: "离职",
	4: "已解约",
}

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
	if err != nil {
		now, err = time.Parse(TimeFormat, "0000-00-00 00:00:00")
		err = nil
	}
	*t = Time(now)
	return
}
func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	if !(time.Time(t)).IsZero() {
		b = time.Time(t).AppendFormat(b, TimeFormat)
	}
	b = append(b, '"')
	return b, nil
}
func (t Time) String() string {
	return time.Time(t).Format(TimeFormat)
}

func (t Time) PosFormat() string {
	return time.Time(t).Format(PosFormat)
}

func (t Time) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t Time) SubToHour(t1 Time) int {
	return int(time.Time(t).Sub(time.Time(t1)) / time.Hour)
}

// Value insert timestamp into mysql need this function.
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if time.Time(t).UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return time.Time(t), nil
}

// Scan value of time.Time
func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

//====date type
func (d *Date) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+DateFormat+`"`, string(data), time.Local)
	*d = Date(now)
	return
}
func (d Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(DateFormat)+2)
	b = append(b, '"')
	b = time.Time(d).AppendFormat(b, DateFormat)
	b = append(b, '"')
	return b, nil
}
func (d Date) String() string {
	return time.Time(d).Format(DateFormat)
}

// Value insert timestamp into mysql need this function.
func (d Date) Value() (driver.Value, error) {
	var zeroTime time.Time
	if time.Time(d).UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return time.Time(d), nil
}

// Scan value of time.Time
func (d *Date) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*d = Date(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
