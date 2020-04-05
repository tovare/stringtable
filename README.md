

![Stringtable](stringtable.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/tovare/stringtable)](https://goreportcard.com/report/github.com/tovare/stringtable)

Stringtable provides a set of utility functions on string-based tables to make data wrangling very easy and fast with a high level of interoperability with other libraries.

For ease of use, we define the type: Table, which is simply a rectangular slice of strings where the first row is column headers.


## get started:

    go get github.com/tovare/stringtable


## Make a Table from scratch

    a := Table{{"one", "two"}, {"1", "2"}}

| one | two |
|-----|-----|
| 1   |  2  |

## Getting data via excelize

    f, err = excelize.OpenFile("data/mydata.xlsx")
    if err != nil {
        return
    }
    mytable := stringtable.Table(f.GetRows(f.GetSheetName(1)))

(github.com/360EntSecGroup-Skylar/excelize)


## Select columns

    a := Table{{"one", "two","three"}, 
                {"1", "2","3"},
                {"1", "2","3"},
                {"1", "2","3"}}
    b := a.Select("three","two")

Table a
| one | two | three |
|-----|-----|-------|
| 1   |  2  |  3    |
| 1   |  2  |  3    |

Table b (result, selected and reordered)
| three | two |
|-------|-----|
|  2    |  3  |
|  2    |  3  |

## Select many columns where columns are in a string array.

    a := Table{{"one", "two","three"}, 
                {"1", "2","3"},
                {"1", "2","3"},
                {"1", "2","3"}}

    mySelection := string{"three", "two", "one"}

    b := a.Select(mySelection ...)

## Rename all columns 

    a := Table{{"one", "two"}, {"1", "2"}}
    a[0] =  []string{"uno", "dos"},

Table a 

| one | two |
|-----|-----|
| 1   |  2  |

Table badge

| uno | dos |
|-----|-----|
| 1   |  2  |

## Rename one column

    a := Table{{"one", "two"}, {"1", "2"}}
    m := a.Colmap()
    a[0][m["two"]] = "dos"

Table a 

| one | two |
|-----|-----|
| 1   |  2  |

Table a (changed)

| one | dos |
|-----|-----|
| 1   |  2  |
