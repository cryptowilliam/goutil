package gexcel

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/extrame/xls"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/tealeg/xlsx"
	"strconv"
	"strings"
)

type (
	XlsDoc struct {
		xlsx *xlsx.File
		xls  *xls.WorkBook
	}

	XlsSheet struct {
		xlsxSheet *xlsx.Sheet
		xlsSheet  *xls.WorkSheet
	}

	Col struct {
		xlsxCol *xlsx.Col
		xlsCol  *xls.Col
	}

	Row struct {
		xlsxCol *xlsx.Row
		xlsCol  *xls.Row
	}

	Cell struct {
		xlsxCol *xlsx.Cell
		xlsCol  *xls.CellRange
	}

	MemSheet struct {
		Name    string
		Content map[string]string // map["RowId,CollId"}"CellContent"
	}

	MemDoc struct {
		Sheets map[uint32]*MemSheet
	}
)

func OpenPath(path string) (*XlsDoc, error) {
	buf, err := gfs.FileToBytes(path)
	if err != nil {
		return nil, err
	}
	return OpenBytes(buf)
}

func OpenBytes(b []byte) (*XlsDoc, error) {
	dx, err := xlsx.OpenBinary(b)
	if err == nil {
		return &XlsDoc{xlsx: dx}, nil
	}

	reader, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, err
	}
	dx, err = xlsx.ReadZipReader(reader)
	if err != nil {
		return nil, err
	}
	return &XlsDoc{xlsx: dx}, nil
}

func (d *XlsDoc) Sheets(idx int) *XlsSheet {
	r := XlsSheet{}
	if d.xlsx != nil {
		r.xlsxSheet = d.xlsx.Sheets[idx]
	} else {
		r.xlsSheet = d.xls.GetSheet(idx)
	}
	return &r
}

func (d *XlsDoc) SheetName(idx int) string {
	if d.xlsx != nil {
		return d.xlsx.Sheets[idx].Name
	} else {
		return d.xls.GetSheet(idx).Name
	}
}

func (d *XlsDoc) SheetCount() int {
	if d.xlsx != nil {
		return len(d.xlsx.Sheets)
	} else {
		return d.xls.NumSheets()
	}
}

func (d *XlsDoc) RowCount(sheetIdx int) int {
	if sheetIdx >= d.SheetCount() {
		return 0
	}

	if d.xlsx != nil {
		return d.xlsx.Sheets[sheetIdx].MaxRow + 1
	} else {
		return int(d.xls.GetSheet(sheetIdx).MaxRow) + 1
	}
}

func (d *XlsDoc) CollCount(sheetIdx, rowIdx int) int {
	if sheetIdx >= d.SheetCount() {
		return 0
	}
	if rowIdx >= d.RowCount(sheetIdx) {
		return 0
	}

	if d.xlsx != nil {
		r, err := d.xlsx.Sheets[sheetIdx].Row(rowIdx)
		if err != nil {
			return 0
		}
		cellCount := 0
		for i := 0; ; i++ {
			if r.GetCell(i) != nil {
				cellCount++
			}
		}
		return cellCount
	} else {
		return int(d.xls.GetSheet(sheetIdx).Row(rowIdx).LastCol()) + 1
	}
}

func (d *XlsDoc) GetCell(sheetIdx, rowIdx, colIdx int) (string, bool) {
	if sheetIdx < 0 || rowIdx < 0 || colIdx < 0 {
		return "", false
	}
	if colIdx >= d.CollCount(sheetIdx, rowIdx) {
		return "", false
	}
	if d.xlsx != nil {
		r, err := d.xlsx.Sheets[sheetIdx].Row(rowIdx)
		if err != nil {
			return "", false
		}
		c := r.GetCell(colIdx)
		if c == nil {
			return "", false
		}
		return c.String(), true
	} else {
		return d.xls.GetSheet(sheetIdx).Row(rowIdx).Col(colIdx), true
	}
}

func (d *XlsDoc) ToMemDoc(useFirstCollCount bool) *MemDoc {
	res := &MemDoc{Sheets: map[uint32]*MemSheet{}}

	_patch_coll_count_of_first_row_ := 0

	sheetCount := d.SheetCount()
	for sheetIdx := 0; sheetIdx < sheetCount; sheetIdx++ {
		res.SetSheet(uint32(sheetIdx), d.SheetName(sheetIdx))
		rowCount := d.RowCount(sheetIdx)
		for rowId := 0; rowId < rowCount; rowId++ {
			collCount := d.CollCount(sheetIdx, rowId)
			if rowId == 0 {
				if sheetIdx == 0 {
					_patch_coll_count_of_first_row_ = collCount
				}
			}
			if useFirstCollCount {
				collCount = _patch_coll_count_of_first_row_
			}
			for collIdx := 0; collIdx < collCount; collIdx++ {
				cellStr, ok := d.GetCell(sheetIdx, rowId, collIdx)
				if ok {
					res.SetCell(uint32(sheetIdx), uint32(rowId), uint32(collIdx), cellStr)
				}
			}
		}
	}
	return res
}

func NewEmptyMemDoc() *MemDoc {
	return &MemDoc{Sheets: map[uint32]*MemSheet{}}
}

func (d *MemDoc) GetSheetName(idx uint32) string {
	if _, ok := d.Sheets[idx]; !ok {
		return ""
	}
	return d.Sheets[idx].Name
}

func (d *MemDoc) GetSheetCount() uint32 {
	return uint32(len(d.Sheets))
}

func (d *MemDoc) GetRowCount(sheetIdx uint32) uint32 {
	maxRowId := uint32(0)
	for rowCollStr := range d.Sheets[sheetIdx].Content {
		ss := strings.Split(rowCollStr, ",")
		if len(ss) != 2 {
			return 0
		}
		rowIdx, err := strconv.ParseInt(ss[0], 10, 64)
		if err != nil {
			return 0
		}
		if uint32(rowIdx) > maxRowId {
			maxRowId = uint32(rowIdx)
		}
	}

	return maxRowId + 1
}

func (d *MemDoc) GetCollCount(sheetIdx uint32) uint32 {
	maxCollId := uint32(0)
	for rowCollStr := range d.Sheets[sheetIdx].Content {
		ss := strings.Split(rowCollStr, ",")
		if len(ss) != 2 {
			return 0
		}
		colIdx, err := strconv.ParseInt(ss[1], 10, 64)
		if err != nil {
			return 0
		}
		if uint32(colIdx) > maxCollId {
			maxCollId = uint32(colIdx)
		}
	}

	return maxCollId + 1
}

func (d *MemDoc) GetCell(sheetIdx, rowIdx, colIdx uint32) string {
	return d.Sheets[sheetIdx].Content[fmt.Sprintf("%d,%d", rowIdx, colIdx)]
}

func (d *MemDoc) SetCell(sheetIdx, rowIdx, colIdx uint32, content string) {
	_, ok := d.Sheets[sheetIdx]
	if !ok {
		d.Sheets[sheetIdx] = &MemSheet{Content: map[string]string{}}
	}

	d.Sheets[sheetIdx].Content[fmt.Sprintf("%d,%d", rowIdx, colIdx)] = content
}

func (d *MemDoc) SetSheet(sheetIdx uint32, name string) {
	_, ok := d.Sheets[sheetIdx]
	if !ok {
		d.Sheets[sheetIdx] = &MemSheet{Content: map[string]string{}}
	}
	d.Sheets[sheetIdx].Name = name
}

func (d *MemDoc) ToXlsx() ([]byte, error) {
	f := xlsx.NewFile()
	sheetCount := d.GetSheetCount()
	for sheetIdx := uint32(0); sheetIdx < sheetCount; sheetIdx++ {
		f.AddSheet(d.GetSheetName(sheetIdx))
		rowCount := d.GetRowCount(sheetIdx)
		collCount := d.GetCollCount(sheetIdx)
		for rowIdx := uint32(0); rowIdx < rowCount; rowIdx++ {
			f.Sheets[sheetIdx].AddRow()
			for collIdx := uint32(0); collIdx < collCount; collIdx++ {
				r, err := f.Sheets[sheetIdx].Row(int(rowIdx))
				if err != nil {
					return nil, err
				}
				r.AddCell()
				s := d.GetCell(sheetIdx, rowIdx, collIdx)
				c := r.GetCell(int(collIdx))
				if c != nil {
					c.SetString(s)
				}
			}
		}
	}
	res := bytes.NewBuffer(nil)
	if err := f.Write(res); err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}
