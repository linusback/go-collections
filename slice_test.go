package collections

import (
	"strconv"
	"testing"
)

func Test_Map(t *testing.T) {
	integers := make(Slice[int64], 0, 10)
	for i := int64(0); i < 10; i++ {
		integers = append(integers, i)
		t.Logf("adding %d\n", i)
	}
	tr := func(i int64) string {
		return strconv.FormatInt(i, 10)
	}
	stringArr := integers.Map(tr)
	for i, v := range stringArr {
		t.Logf("strings %d: %s\n", i, v)
	}
}
