package types

import (
	"fmt"
	"strconv"
	"strings"
)

type BandwidthWithQuantity struct {
	s string // MB or KB
	i int64  // bytes
}

type PortsRange struct {
	Start  int `json:"start,omitempty"`
	End    int `json:"end,omitempty"`
	Single int `json:"single,omitempty"`
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
