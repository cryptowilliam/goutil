package gcompress

import "github.com/cryptowilliam/goutil/basic/gerrors"

type (
	CompParam struct {
		Level int
	}
)

func (p CompParam) Verify(algo Comp) error {
	switch algo {
	// Flate level must in [-2, 9].
	case CompFlate:
		if p.Level < -2 || p.Level > 9 {
			return gerrors.New("compress algorithm %s level %d is out of [-2, 9]", algo, p.Level)
		}
	}

	return nil
}
