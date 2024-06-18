package types

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	MB = 1024 * 1024
	KB = 1024

	BandwidthLimitModeClient = "client"
	BandwidthLimitModeServer = "server"
)

type BandwidthQuantity struct {
	s string // MB or KB

	i int64 // bytes
}

func NewBandwidthQuantity(s string) (BandwidthQuantity, error) {
	q := BandwidthQuantity{}
	err := q.UnmarshalString(s)
	if err != nil {
		return q, err
	}
	return q, nil
}

func (q *BandwidthQuantity) String() string {
	return q.s
}

func (q *BandwidthQuantity) UnmarshalString(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	var (
		base int64
		f    float64
		err  error
	)
	switch {
	case strings.HasSuffix(s, "MB"):
		base = MB
		fstr := strings.TrimSuffix(s, "MB")
		f, err = strconv.ParseFloat(fstr, 64)
		if err != nil {
			return err
		}
	case strings.HasSuffix(s, "KB"):
		base = KB
		fstr := strings.TrimSuffix(s, "KB")
		f, err = strconv.ParseFloat(fstr, 64)
		if err != nil {
			return err
		}
	default:
		return errors.New("unit not support")
	}

	q.s = s
	q.i = int64(f * float64(base))
	return nil
}

func (q *BandwidthQuantity) Bytes() int64 {
	return q.i
}

type PortsRange struct {
	Start  int `json:"start,omitempty"`
	End    int `json:"end,omitempty"`
	Single int `json:"single,omitempty"`
}

type PortsRangeSlice []PortsRange

func (p PortsRangeSlice) String() string {
	if len(p) == 0 {
		return ""
	}
	strs := make([]string, 0)
	for _, v := range p {
		if v.Single > 0 {
			strs = append(strs, strconv.Itoa(v.Single))
		} else {
			strs = append(strs, strconv.Itoa(v.Start)+"-"+strconv.Itoa(v.End))
		}
	}
	return strings.Join(strs, ",")
}

// NewPortsRangeSliceFromString the format of str is like "1000-2000,3000,4000-5000"
func NewPortsRangeSliceFromString(str string) ([]PortsRange, error) {
	str = strings.TrimSpace(str)
	numRanges := strings.Split(str, "")
	out := make([]PortsRange, 0)
	for _, numRangeStr := range numRanges {
		// 1000-2000 or 2001
		numArray := strings.Split(numRangeStr, "-")
		// length: 只有1或2是正确的
		rangType := len(numArray)
		switch rangType {
		case 1:
			// 单独数字
			singleNum, err := strconv.ParseInt(strings.TrimSpace(numArray[0]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("range number is invalid, %v", err)
			}
			out = append(out, PortsRange{Single: int(singleNum)})
		case 2:
			// range numbers
			minP, err := strconv.ParseInt(strings.TrimSpace(numArray[0]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("range number is invalid, %v", err)
			}
			maxP, err := strconv.ParseInt(strings.TrimSpace(numArray[1]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("range number is invalid, %v", err)
			}
			if minP < maxP {
				return nil, fmt.Errorf("range number is invalid")
			}
			out = append(out, PortsRange{Start: int(minP), End: int(maxP)})
		default:
			return nil, fmt.Errorf("range number is invalid")
		}
	}
	return out, nil
}
