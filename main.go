package main

import (
	"database/sql"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	excelize "github.com/xuri/excelize/v2"
	"os"
)

type Env struct {
	PrefixShortname map[string]string
	PrefixOrdered   []string
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <YYYYMM> db.conf\n", os.Args[0])
		return
	}
	dbConf := os.Args[2]

	err := godotenv.Load(os.ExpandEnv(dbConf))
	if err != nil {
		panic(err.Error())
	}
	dsn := os.Getenv("TS_DB")
	tsdb, err = sql.Open("clickhouse", dsn)
	if err != nil {
		panic(err.Error())
	}

	dsn = os.Getenv("DB_DSN")
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	e := &Env{}
	e.initPrefix()
	initCounter()
	date := os.Args[1]
	fn := "AASYR" + date + ".xlsx"

	fmt.Println("Doing job for AltynAsyr, month ", date)
	f := excelize.NewFile()
	activeSheet := "Sheet1"
	nextRaw := createHeader(f, activeSheet, date)

	e.createInternationalCalls(nextRaw)

	if err := f.SaveAs(fn); err != nil {
		fmt.Println(err)
	}
}

func (e *Env) createInternationalCalls(nextRaw int) {
	e.trunkTotalByDestination()
}

func createHeader(f *excelize.File, sheet string, date string) int {
	if err := f.SetColWidth(sheet, "A", "E", 20.00); err != nil {
		fmt.Println(err)
		return 0
	}

	f.SetCellValue(sheet, "A1", "ГКЭ \"Туркментелеком\"")
	f.SetCellValue(sheet, "A2", "Расчетный счет для манатов")
	f.SetCellValue(sheet, "A3", "МФО:")
	f.SetCellValue(sheet, "A4", "НК")
	f.SetCellValue(sheet, "A5", "Станция \"Алтын асыр\"")

	f.SetCellValue(sheet, "C1", "хранить  3 года")
	f.SetCellValue(sheet, "C2", "638501")
	f.SetCellValue(sheet, "C3", "390101201")
	f.SetCellValue(sheet, "C4", "01161000158")
	f.SetCellValue(sheet, "C5", date)
	var boldFont int
	var err error
	if boldFont, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	}); err != nil {
		panic(err.Error())
	}

	if err := f.SetCellStyle(sheet, "A7", "E7", boldFont); err != nil {
		fmt.Println(err)
		return 0
	}

	f.SetCellValue(sheet, "A7", "СТРАНА")
	f.SetCellValue(sheet, "B7", "КОЛ-ВО ЗАКАЗОВ")
	f.SetCellValue(sheet, "C7", "КОЛ-ВО МИН")
	f.SetCellValue(sheet, "D7", "ТАРИФ В МАНАТАХ")
	f.SetCellValue(sheet, "E7", "ОПЛАТА В МАНАТАХ")

	return 8
}
