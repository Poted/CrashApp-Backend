package helperz

import (
	"fmt"
	"runtime"
)

func AvailableGoroutines() int {
	var MaxGRChan = make(chan int, 1)
	go monitorSystemLoad(0, true, MaxGRChan)
	return <-MaxGRChan
}

// 'bufferPercentage' is a parcentage value of whole accesible go routines number that
// are wanted.
func monitorSystemLoad(bufferPercentage int, printLog bool, availableGoRoutines chan (int)) {

	var bp int
	var ret int
	// for {
	switch bp = bufferPercentage; {
	case bp < 0:
		ret = 0
	case bp == 0:
		ret = int(float32(runtime.NumCPU()*2 - runtime.NumGoroutine()))
	case bp > 0:
		if bufferPercentage >= 100 {
			bufferPercentage = 100
		}
		buf := float32(bufferPercentage) / 100
		ret = int(float32((runtime.NumCPU()*2 - runtime.NumGoroutine())) * buf)
	}

	// time.Sleep(1 * time.Second) // Adjust monitoring interval as needed
	if printLog {
		fmt.Printf("\nRunning Goroutines:: %v\n", runtime.NumGoroutine())
		fmt.Printf("Available Goroutines:: %v\n", runtime.NumCPU()*2)
		fmt.Printf("Returned Goroutines:: %v\n", ret)

	}
	availableGoRoutines <- ret
	// return
	// }

}
