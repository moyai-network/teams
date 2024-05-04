package sotw

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// logFile will temporarily store the timestamp of the SOTW
const logFile = "assets/sotw.log"

// sotw represents the private global object of the SOTW time
var sotw time.Time

func init() {
	if _, err := os.ReadFile(logFile); err != nil {
		_ = os.WriteFile(logFile, []byte(fmt.Sprint(time.Now().Unix())), 0777)
		return
	}

	d, _ := os.ReadFile(logFile)
	i, _ := strconv.ParseInt(string(d), 10, 64)
	sotw = time.Now().Add(time.Duration(i))
}

// Start starts the SOTW timer
func Start() {
	sotw = time.Now().Add(1 * time.Hour)
}

// End ends the SOTW timer
func End() {
	sotw = time.Time{}
}

// Save will save the SOTW time to the log file
func Save() {
	s := fmt.Sprint(int64(time.Until(sotw)))
	_ = os.WriteFile(logFile, []byte(s), 0777)
}

// Running returns the SOTW time and if it is running
func Running() (time.Time, bool) {
	return sotw, !sotw.IsZero() && time.Now().Before(sotw)
}
