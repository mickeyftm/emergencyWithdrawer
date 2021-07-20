package utils

import "github.com/jon4hz/emergenyWithdrawer/internal/logging"

func ExitHandler() {

	// stop logging feed
	logging.StopChan <- struct{}{}

	// close all open channels
	close(logging.InfoChan)
	close(logging.WarnChan)
	close(logging.ErrChan)
}
