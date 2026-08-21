package tsrollup

import "errors"

var (
	ErrClosed          = errors.New("tsrollup: engine closed")
	ErrInvalid         = errors.New("tsrollup: invalid argument")
	ErrNotFound        = errors.New("tsrollup: series not found")
	ErrNoAgg           = errors.New("tsrollup: no aggregator factory")
	ErrCorruptSegment  = errors.New("tsrollup: corrupt segment")
	ErrPersist         = errors.New("tsrollup: persist failed")
	ErrCanceled        = errors.New("tsrollup: canceled")
	ErrWAL             = errors.New("tsrollup: wal failure")
	ErrWindowOverlap   = errors.New("tsrollup: window overlap")
	ErrEmpty           = errors.New("tsrollup: empty result")
)
