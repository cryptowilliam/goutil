package gscore

import "github.com/cryptowilliam/goutil/basic/gerrors"

type sub struct {
	scaleInTotal                              int // 此sub项目的占比，所有sub的这个值的总和必须等于100
	value, valueOfZeroScore, valueOfFullScore float64
}

// 按比例计算的比分计算器
type ScoreCalc struct {
	zeroScore float64 // 最低起点分
	fullScore float64 // 最高顶点分
	sublist   []sub
}

func NewScoreCalc(zeroScore, fullScore float64) (*ScoreCalc, error) {
	if zeroScore >= fullScore {
		return nil, gerrors.Errorf("Invalid valueOfZeroScore(%f) and valueOfFullScore(%f)", zeroScore, fullScore)
	}
	return &ScoreCalc{zeroScore: zeroScore, fullScore: fullScore}, nil
}

func (s *ScoreCalc) AddSub(scaleInTotal int, value, valueOfZeroScore, valueOfFullScore float64) {
	it := sub{scaleInTotal: scaleInTotal, value: value, valueOfZeroScore: valueOfZeroScore, valueOfFullScore: valueOfFullScore}
	if it.value < it.valueOfZeroScore {
		it.value = it.valueOfZeroScore
	}
	if it.value > it.valueOfFullScore {
		it.value = it.valueOfFullScore
	}
	s.sublist = append(s.sublist, it)
}

func (s *ScoreCalc) GetScore() (float64, error) {
	totalScale := 0
	for _, v := range s.sublist {
		totalScale += v.scaleInTotal
	}
	if totalScale != 100 {
		return 0, gerrors.Errorf("Correct total scale 100, but get %d", totalScale)
	}

	totalScore := float64(0)
	for _, v := range s.sublist {
		totalScore += float64(v.scaleInTotal) * ((v.value - v.valueOfZeroScore) / (v.valueOfFullScore - v.valueOfZeroScore))
	}

	return s.zeroScore + ((s.fullScore - s.zeroScore) * (float64(totalScore) / float64(100))), nil
}
