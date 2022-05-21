package gany

import "encoding/json"

func ConvertJsonNumberToFloat64Array(s []interface{}) ([]float64, error) {
	var result []float64

	for i := range s {
		v, err := s[i].(json.Number).Float64()
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}
