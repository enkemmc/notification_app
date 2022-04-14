package tools

import (
	"fmt"
	"time"
)

func PrintWithTimestamp(s string) {
	now := time.Now().Format(time.Kitchen)
	fmt.Printf("%s %s\n", now, s)
}
