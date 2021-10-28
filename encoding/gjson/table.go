package gjson

// 矩形表格，类似excel中的一个sheet

type Row []string

type Table struct {
	Rows []Row
}

/*
{
   "Rows":[
     ["","",""], // Row
     ["","",""], // Row
     ["","",""] // Row
   ]
}
*/

type Column struct {
	Cells []Cell
}

/*
type Row struct {
	Cells []Cell
}*/

type Cell struct {
	columnIndex int
	rowIndex    int
	Data        string
}

func NewTableFromString(json string) (*Table, error) {
	return nil, nil
}

func NewTableFromExcel(filename string) (*Table, error) {
	return nil, nil
}

func (d *Table) GetColumnsCount() (int, error) {
	return 0, nil
}

func (d *Table) GetColumn(columnIndex int) (Column, error) {
	return Column{}, nil
}

func (d *Table) AddColumn(afterColumnIndex int, newColumn Column) error {
	return nil
}

func (d *Table) RemoveColumn(columnIndex int) error {
	return nil
}

func (d *Table) GetRowsCount() (int, error) {
	return 0, nil
}

func (d *Table) GetRow(rowIndex int) (Row, error) {
	return Row{}, nil
}

func (d *Table) AddRow(afterRowIndex int, newRow Row) error {
	return nil
}
