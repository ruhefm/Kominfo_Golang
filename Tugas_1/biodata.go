package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type kelas struct {
	ID        int
	Nama      string
	Alamat    string
	Pekerjaan string
	Alasan    string
}

var lastID int
var siswa []kelas

func main() {
	argumen := os.Args
	if len(argumen) > 2 {
		deac, _ := strconv.Atoi(argumen[2])
		if (deac) != 0 {
			excelToKelas()
		}
	} else {
		excelToKelas()
	}

	tambahAnggota("Heru", "Jatihandap", "FGA", "Mempelajari GoLang")
	tambahAnggota("Irman", "Cigending", "FGA", "Mempelajari Bahasa")
	fmt.Println("Path program:", argumen[0])
	fmt.Println("Untuk akses berdasar absen, gunakan argumen pertama contoh go run biodata.go 1 untuk menampilkan absen 1, untuk menonaktifkan excel gunakan 0 pada argumen ke 2.")
	fmt.Println("Argumen:", argumen[1:])
	absen, _ := strconv.Atoi(argumen[1])

	if absen >= 0 && absen < len(siswa)+1 {
		fmt.Printf("Kamu memilih absen nomor %v \n Dengan identitas sebagai berikut \n %v", absen, siswa[absen-1])
	} else {
		fmt.Printf("Kamu memilih absen nomor %v \n Akan tetapi absen yang kamu cari tidak ditemukan", absen)

	}
}

func excelToKelas() {
	excel, err := excelize.OpenFile("./kelas.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	rows, err := excel.GetRows("Table 1")
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, row := range rows {
		if i == 0 {

			continue
		}
		if len(row) < 6 {
			continue
		}
		if lastID == 0 {
			lastID = i
		} else {
			lastID++
		}
		anggota := kelas{
			ID:        lastID,
			Nama:      row[2],
			Alamat:    row[3],
			Pekerjaan: row[4],
			Alasan:    row[5],
		}
		siswa = append(siswa, anggota)
		fmt.Printf("%v\n", siswa[lastID-1])

	}
}

func tambahAnggota(Nama string, Alamat string, Pekerjaan string, Alasan string) {
	lastID++
	anggota := kelas{
		ID:        lastID,
		Nama:      Nama,
		Alamat:    Alamat,
		Pekerjaan: Pekerjaan,
		Alasan:    Alasan,
	}
	siswa = append(siswa, anggota)
	fmt.Printf("%v\n", siswa[lastID-1])
	// fmt.Printf("%v\n", siswa)

}

// Dibawah ini untuk program langsung akses excel.

// func main() {
// 	filePath := flag.String("file", "", "path to Excel file")
// 	sheetName := flag.String("sheet", "Sheet1", "Sheet name, default is Sheet1")
// 	noAbsen := flag.String("absen", "", "Row number to retrieve, default is +1 from the last read row")
// 	// cariHeader := flag.String("cari", "Nama", "Silakan pilih header, sesuai dengan header yang akan dicari.")

// 	flag.Parse()
// 	if *filePath == "" {
// 		fmt.Println("Please provide a path to the Excel file using the -file flag.")
// 		return
// 	}

// 	excel, err := excelize.OpenFile(*filePath)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	if *noAbsen == "" {
// 		rows, err := excel.GetRows(*sheetName)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}

// 		for _, row := range rows {
// 			for _, colCell := range row {
// 				fmt.Print(colCell, "\t")
// 			}
// 			fmt.Println()
// 		}
// 	} else {
// 		rowNum, err := strconv.Atoi(*noAbsen)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}

// 		rowNumber := "C" + strconv.Itoa(rowNum+1)

// 		absen, err := excel.GetCellValue(*sheetName, rowNumber)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		fmt.Printf("Absen %s: %s\n", *noAbsen, absen)
// 	}
// }
