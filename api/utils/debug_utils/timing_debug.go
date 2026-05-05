package debug_utils

import (
	"fmt"
	"time"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s", name, elapsed)
}
