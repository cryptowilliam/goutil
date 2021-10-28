package gcsv

// note:
// some tools export bad csv file which has invisible characters at the beginning of csv file like 0x(EFBBBF), if you find that some member can't be Unmarshal, check all members of csv file

type CsvDoc struct {
}

func OpenCsv(filename string) (*CsvDoc, error) {
	return nil, nil
}

func (d *CsvDoc) ExportJson() (string, error) {
	return "", nil
}
