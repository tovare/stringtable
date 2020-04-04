package stringtable

import (
	"os"
	"testing"
)

// TestReadCSV ... test that we can read the file correctly.
func TestWriteCSV(t *testing.T) {
	a := Table{{"one", "two"}, {"1", "2"}}
	err := a.WriteCSV("test_tmp.csv")
	if err != nil {
		t.Errorf("Failed to save %v", err)
	}
	_, err = ReadCSV("test_tmp.csv")
	if err != nil {
		t.Errorf("Failed rountrip %v", err)
	}
	err = os.Remove("test_tmp.csv")
	if err != nil {
		t.Errorf("Failed to clean up test-file.")
	}

}

// TestAppend ... test catentation of two tables.
func TestAppend(t *testing.T) {

	a := Table{{"one", "two"}, {"1", "2"}}
	b := Table{{"one", "two"}, {"1", "2"}}

	res := Table{{"one", "two"}, {"1", "2"}, {"1", "2"}}
	c := a.Append(b)
	for rndx, row := range c {
		for cndx, v := range row {
			if v != res[rndx][cndx] {
				t.Logf("Append equality test failed.")
			}
		}
	}

	// Empty table
	d := Table{}
	d = d.Append(c)
	for rndx, row := range d {
		for cndx, v := range row {
			if v != res[rndx][cndx] {
				t.Logf("Append equality test for empty table failed.")
			}
		}
	}
}

func TestColmap(t *testing.T) {
	testTable := Table{{"a", "b"}, {"1", "2"}, {"3", "4"}, {"1", "5"}}
	c := testTable.Colmap()
	if c["a"] != 0 {
		t.Logf("Failed identifying a on index 0")
	}
	if c["b"] != 1 {
		t.Logf("Failed to identify b on index 1")
	}
}

func TestFilter(t *testing.T) {
	testTable := Table{{"a", "b"}, {"1", "2"}, {"3", "4"}, {"1", "5"}}
	targetTable := Table{{"a", "b"}, {"1", "2"}, {"1", "5"}}
	newTable := testTable.Filter(func(row []string, c ColumnIndices) bool {
		if row[c["a"]] == "1" {
			return true
		}
		return false
	})
	for ri, row := range newTable {
		for ci, s := range row {
			if targetTable[ri][ci] != s {
				t.Logf("Failed comparison %v in result != %v in target", s, targetTable[ri][ci])
			}
		}
	}
}

func TestLeftJoin(t *testing.T) {
	a := Table{
		{"queue", "data1"},
		{"one", "1"},
		{"two", "2"},
		{"one", "11"},
		{"one", "111"},
		{"two", "22"},
		{"three", "3"},
	}

	b := Table{
		{"queue", "data2", "data3"},
		{"one", "d2", "d3"},
		{"two", "t2", "t3"},
	}

	expected := Table{
		{"queue", "data1", "data2", "data3"},
		{"one", "1", "d2", "d3"},
		{"two", "2", "t2", "t3"},
		{"one", "11", "d2", "d3"},
		{"one", "111", "d2", "d3"},
		{"two", "22", "t2", "t3"},
		{"three", "3", "", ""},
	}

	c := a.LeftJoin(b)

	for ri, row := range c {
		for rc, col := range row {
			if r := expected[ri][rc]; col != r {
				t.Errorf("Unexpected result during LeftJoin found %v expected %v", col, r)
			}
		}
	}
	// TODO: Trigger error instead.
	t.Log(a, b, expected)

}

func TestAddColumn(t *testing.T) {
	test1 := Table{
		{"AA", "BB"},
		{"1a", "1b"},
		{"2a", "2b"},
	}
	res1 := Table{
		{"AA", "BB", "CC"},
		{"1a", "1b", "1c"},
		{"2a", "2b", "2c"},
	}

	if !test1.Equal(test1) {
		t.Errorf("Failed initial equality test.")
	}

	if test1.Equal(res1) {
		t.Errorf("Failed initial inequality test.")
	}

	testResult := test1.AddColumns(Table{{"CC"}, {"1c"}, {"2c"}})
	if !testResult.Equal(res1) {
		t.Errorf("Failed first test, not equeal")
	}

	res2 := Table{
		{"AA", "BB", "CC"},
		{"1a", "1b", "1b"},
		{"2a", "2b", ""},
	}

	testResult2 := test1.AddColumns(Table{{"CC"}, {"1b"}})
	if !testResult2.Equal(res2) {
		t.Error("Failed filler value test.")
	}

}

func TestSelect(t *testing.T) {

	test1 := Table{
		{"queue", "data2", "data3"},
		{"one", "d2", "d3"},
		{"two", "t2", "t3"},
	}

	// Basic selection

	expected1 := Table{
		{"data2", "data3"},
		{"d2", "d3"},
		{"t2", "t3"},
	}

	result, err := test1.Select("data2", "data3")
	if err != nil {
		t.Error(err)
	}

	for ri, row := range result {
		for ci, col := range row {
			if rv := expected1[ri][ci]; col != rv {
				t.Errorf("T2: Unexpected value in select result. Expected %v but got %v in result row %v", rv, col, ri)
			}
		}
	}

	// Reversed column order.

	expected2 := Table{
		{"data3", "data2"},
		{"d3", "d2"},
		{"t3", "t2"},
	}

	result, err = test1.Select("data3", "data2")
	if err != nil {
		t.Error(err)
	}

	for ri, row := range result {
		for ci, col := range row {
			if rv := expected2[ri][ci]; col != rv {
				t.Errorf("Unexpected value in reversed result. Expected %v but got %v in result row %v", rv, col, ri)
			}
		}
	}

}
