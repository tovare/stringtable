package stringtable

/*
	package stringtable implements methods to manipulate a rectangular array of strings where
	the convention is that the first element of each column is a header value.

	Usage

	Use the Table type to signify that the [][]string data is rectangular and with
	the first line as a header. All operations are immutable unless documented.

	You can declare a new table like this:

		mytable := Table{  {"column1", "column2"} {"row1", "row1"}, { "row2", row2} }

	If you convert data in this format from another source you can assert the type like this:

		mytable := Table(resultFromOtherAPI)

	Utility features:

	* ReadCSV .. read a semicolon separated CSV
    * WriteCSV writes the table to a file

	The package currently support a number of methods to manipulate a table

	* Filter	keep all the elements which satisfy your filter function.
	* Append	concatenaate two tables with the same rows.
	* LeftJoin	Join tables on common key.
	* AddColumn	Add a column with name and the correct length.
	* Select	Select columns by name in preferred order.
	* Colmap    Get a map of column names with their index.

	Future extensions:

	* SplitN	TODO split the table into n equal parts.
	* GetRandomRows TODO get a random selction of rows.
	* Apply		MAYBE mutate a column (Just iterate for now starting with the first row of data)
	* Reduce	MAYBE NOT reduce to a target table
	* Head		MAYBE get the first 10 lines.

*/

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

// Table defines a table of strings where the first row contains the header.
type Table [][]string

// ColumnIndices maps the name of a column to it´s index.
type ColumnIndices map[string]int

/************************************************************************
 * Utility functions
 ************************************************************************/

// ReadCSVTrimmed reads a CSV file and trim whitespace in each field.
func ReadCSVTrim(filename string) (table Table, err error) {
	table, err = ReadCSV(filename)
	for i, line := range table {
		for j, col := range line {
			table[i][j] = strings.TrimSpace(col)
		}
	}
	return
}

// ReadCSV reads a CSV file and return a table structure delimited with semicolons.
func ReadCSV(filename string) (table Table, err error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return
	}
	csv := csv.NewReader(f)
	csv.Comma = ';'
	csv.ReuseRecord = false
	t, err := csv.ReadAll()
	if err != nil {
		return
	}

	table = Table(t)
	return
}

// WriteCSV writes the table to a Semicolon separated file.
func (m Table) WriteCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	for _, value := range m {
		err := writer.Write(value)
		if err != nil {
			return err
		}
	}
	return err
}

/************************************************************************
 * Table operations
 ************************************************************************/

// Append two tables with equal number of columns. The initial table may be empty.
// This function will terminate with fatal if the number of columns varies.
func (m Table) Append(a Table) Table {
	if len(m) > 0 {
		if len(m[0]) != len(a[0]) {
			log.Fatalf("Columns are not the same this column has %v columns, but the other table has %v columns."+
				" Current table has %v columns but joining table has %v columns", len(m[0]), len(a[0]), m[0], a[0])
		}
		m = append(m, a[1:]...)
	} else {
		m = a
	}
	return m
}

// AddColumns to the table. The Supplied columns must have the first row as a header.
// If the columns are too short, empty filler values are added
// to preserve the rectangular shape of the table.
func (m Table) AddColumns(columns Table) (newTable Table) {

	filler := make([]string, len(columns[0]))

	newTable = make(Table, len(m))

	for i, row := range m {
		if i < len(columns) {
			newTable[i] = append(row, columns[i]...)
		} else {
			newTable[i] = append(row, filler...)
		}
	}

	return
}

// Colmap returns map of all the column names with which index they have in the table,
// this is to make it easier to intuitively write filter-functions and build this only
// once when using a filter.
func (m Table) Colmap() (columns ColumnIndices) {
	columns = make(ColumnIndices, len(m[0]))
	for i, col := range m[0] {
		columns[col] = i
	}
	return
}

// Filter .. keep lines satisfying the formula into a new table. Header row is skipped.
func (m Table) Filter(f func(row []string, c ColumnIndices) bool) (result Table) {
	result = Table{}
	result = append(result, m[0])
	cmap := m.Colmap()
	for _, line := range m[1:] {
		if f(line, cmap) {
			result = append(result, line)
		}
	}
	return result
}

// LeftJoin will copy it´s values for lines with matching keys on the left side, will mutate
// left side. Right side needs to bee unique. The key must be the first column. *This function
// is mutable* If the right side of the join doesn´t have the key, an empty filler will be
// added to keep the rectangular shape of the table.
func (m Table) LeftJoin(right Table) Table {
	// Append new headers
	newcols := right[0][1:]
	m[0] = append(m[0], newcols...)
	// I assume that there may be lines that do not match and
	// multiple instances on the left side or right side. I need to create a map of string arrays
	// missing.
	k := make(map[string][]string)
	for _, v := range right[1:] {
		k[v[0]] = v[1:]
	}
	// The empty slice of strings
	filler := make([]string, len(newcols))
	for i := range m[1:] {
		sl := k[m[i+1][0]]
		if sl == nil {
			// No key found, add filler values.
			// TODO: Consider deleting these lines in the future.
			m[i+1] = append(m[i+1], filler...)
		} else {
			// Append thee columns.
			m[i+1] = append(m[i+1], sl...)
		}
	}
	return m
}

// Select creates a table with the column names supplied and in that order.
func (m Table) Select(columns ...string) (table Table, err error) {
	// Iterate the header and create a list of indices.
	indices := make([]int, len(columns))
	for i, colname := range columns {
		found := false
		foundNdx := 0
		for j, target := range m[0] {
			found = target == colname
			if found {
				foundNdx = j
				break
			}
		}
		indices[i] = foundNdx
		if !found {
			return table, fmt.Errorf("Column %v not found in table", colname)
		}
	}
	// Create a table with the selected columns.
	table = Table{}
	table = append(table, columns)
	for _, row := range m[1:] {
		tmp := make([]string, len(indices))
		for i := range indices {
			tmp[i] = row[indices[i]]
		}
		table = append(table, tmp)
	}
	return
}

//******************************************************
// Logical tests
//*******************************************************

func (m Table) Equal(right Table) (ok bool) {

	if len(m) != len(right) {
		return false
	}
	if len(m[0]) != len(right[0]) {
		return false
	}

	// If the tables are of eaqual shape
	// we need to test the contents.
	for i, row := range m {
		for j, col := range row {
			ok = col == right[i][j]
			if !ok {
				return
			}
		}
	}
	return
}
