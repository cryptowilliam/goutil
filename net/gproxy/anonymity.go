package gproxy

const (
	ANONYMITY_NOA = iota
	ANONYMITY_ANM
	ANONYMITY_HIA
)

type Anonymity int

func (a Anonymity) String() string {
	if a == ANONYMITY_NOA {
		return "NOA"
	}
	if a == ANONYMITY_ANM {
		return "ANM"
	}
	if a == ANONYMITY_HIA {
		return "HIA"
	}
	return "UNKNOWN"
}
