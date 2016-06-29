package logging

import (
	"fmt"
	"time"
)

var Debug bool = true

func DebugLogging(text string) {
	if Debug {
		currentTime := time.Now().Local()
		fmt.Println("[" + currentTime.Format(time.RFC850) + "] " + text)
	}
}
