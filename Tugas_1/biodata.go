package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"

)

func main() {
	filePath := flag.String("file", "", "path to Excel file")
	sheetName := flag.String("sheet", "Sheet1", "Sheet name, default is Sheet1")
	noAbsen := flag.String("absen", "", "Row number to retrieve, default is +1 from the last read row")
	// cariHeader := flag.String("cari", "Nama", "Silakan pilih header, sesuai dengan header yang akan dicari.")

	flag.Parse()
	if *filePath == "" {
		fmt.Println("Please provide a path to the Excel file using the -file flag.")
		return
	}

	excel, err := excelize.OpenFile(*filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *noAbsen == "" {
		rows, err := excel.GetRows(*sheetName)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, row := range rows {
			for _, colCell := range row {
				fmt.Print(colCell, "\t")
			}
			fmt.Println()
		}
	} else {
		rowNum, err := strconv.Atoi(*noAbsen)
		if err != nil {
			fmt.Println(err)
			return
		}

		rowNumber := "C" + strconv.Itoa(rowNum+1)

		absen, err := excel.GetCellValue(*sheetName, rowNumber)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Absen %s: %s\n", *noAbsen, absen)
	}
}