package dates

import (
	"testing"
	"fmt"
)

func TestToString(t *testing.T) {
	s := ToDateStr()
	fmt.Printf(s)
}
