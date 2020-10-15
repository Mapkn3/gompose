package util

import "log"

// Check checks error and log message if not nill
func Check(err error) {
	if err != nil {
		log.Panicf("%v", err)
	}
}
