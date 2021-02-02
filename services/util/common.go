package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// ComputeMD5 计算文件的md5码
func ComputeMD5(filePath string) (string, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(result)), nil
}

func StringMd5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

const otpChars = "1234567890"

func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, 0)
}

// 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// note: this is off by two days on the real epoch (1/1/1900) because
// - the days are 1 indexed so 1/1/1900 is 1 not 0
// - Excel pretends that Feb 29, 1900 existed even though it did not
// The following function will fail for dates before March 1st 1900
// Before that date the Julian calendar was used so a conversion would be necessary
var excelEpoch = time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)

func ExcelDateToDate(excelDate string) (time.Time, error) {
	var days, err = strconv.Atoi(excelDate)
	if err != nil {
		return excelEpoch, err
	}
	return excelEpoch.Add(time.Second * time.Duration(days*86400)), nil
}
