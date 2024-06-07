package util

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

func GetAuthKey(token string, timestamp int64) (key string) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(token))
	md5Ctx.Write([]byte(strconv.FormatInt(timestamp, 10)))
	data := md5Ctx.Sum(nil)
	return hex.EncodeToString(data)
}

func ParseRangeNumbers(rangeStr string) (numbers []int64, err error) {
	rangeStr = strings.TrimSpace(rangeStr)
	numbers = make([]int64, 0)
	// 例如 1000-2000,2001,2002,3000-4000
	numRange := strings.Split(rangeStr, ",")
	for _, numRangeStr := range numRange {
		// 1000-2000 or 2001
		numArray := strings.Split(numRangeStr, "-")
		// 长度是1或2才正确
		rangeType := len(numArray)
		switch rangeType {
		case 1:
			// 单独数字
			singleNum, errRet := strconv.ParseInt(strings.TrimSpace(numArray[0]), 10, 64)
			if errRet != nil {
				err = fmt.Errorf("range number is invalid, %v", errRet)
				return
			}
			numbers = append(numbers, singleNum)
		case 2:
			// 遍历数字
			minN, errRet := strconv.ParseInt(strings.TrimSpace(numArray[0]), 10, 64)
			if errRet != nil {
				err = fmt.Errorf("invalid number is invalid, %v", errRet)
				return
			}
			maxN, errRet := strconv.ParseInt(strings.TrimSpace(numArray[1]), 10, 64)
			if errRet != nil {
				err = fmt.Errorf("range number is invalid, %v", errRet)
				return
			}
			if maxN < minN {
				err = fmt.Errorf("range number is invalid")
				return
			}
			for i := minN; i <= maxN; i++ {
				numbers = append(numbers, i)
			}
		default:
			err = fmt.Errorf("range number is invalid ")
			return
		}
	}
	return
}

func ConstantTimeEqString(a, b string) bool {
	// todo 学习
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
