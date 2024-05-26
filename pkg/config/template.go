package config

import (
	"fmt"
	"github.com/sunyihoo/frp/pkg/util/util"
)

type NumberPair struct {
	First  int64
	Second int64
}

func parseNumberRangePair(firstRangeStr, secondRangeStr string) ([]NumberPair, error) {
	firstRangeNumbers, err := util.ParseRangeNumbers(firstRangeStr)
	if err != nil {
		return nil, err
	}
	secondRangeNumbers, err := util.ParseRangeNumbers(secondRangeStr)
	if err != nil {
		return nil, err
	}
	if len(firstRangeNumbers) != len(secondRangeNumbers) {
		return nil, fmt.Errorf("first and second range numbers are in paris")
	}
	pairs := make([]NumberPair, 0, len(firstRangeNumbers))
	for i := 0; i < len(firstRangeNumbers); i++ {
		pairs = append(pairs, NumberPair{
			First:  firstRangeNumbers[i],
			Second: secondRangeNumbers[i],
		})
	}
	return pairs, nil
}

func parseNumberRange(firstRangeStr string) ([]int64, error) {
	return util.ParseRangeNumbers(firstRangeStr)
}
