package gtime

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"time"
)

func HasDuplicated(times []time.Time, accuracy time.Duration) (bool, error) {
	switch accuracy {
	case time.Nanosecond:
		//truncateAccuracy = time.Duration(0)
	case time.Microsecond:
		//truncateAccuracy = time.Nanosecond
	case time.Millisecond:
		//truncateAccuracy = time.Microsecond
	case time.Second:
		//truncateAccuracy = time.Millisecond
	case time.Minute:
		//truncateAccuracy = time.Second
	case time.Hour:
	//truncateAccuracy = time.Minute
	default:
		return false, gerrors.Errorf("unsupported verify duration %s", accuracy.String())
	}
	/*if truncateAccuracy >= Day {
		return false, gerrors.Errorf("unsupported verify duration %s", accuracy.String())
	}*/

	cmpMap := make(map[int64]bool)
	for _, v := range times {
		v = v.Truncate(accuracy)
		if cmpMap[v.UnixNano()] {
			return true, nil
		} else {
			cmpMap[v.UnixNano()] = true
		}
	}
	return false, nil
}
