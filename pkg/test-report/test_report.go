package test_report

import (
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
	"strings"
	"sync"
)

type reportStyles struct {
	failure   int
	success   int
	notTested int
	headers   int
}

type Report struct {
	Path         string
	file         *excelize.File
	nextRow      int
	sheetName    string
	mutex        *sync.Mutex
	styleNumbers reportStyles
	logger       *log.Logger
}

const (
	top    string = "top"
	left   string = "left"
	right  string = "right"
	bottom string = "bottom"
)

const continuousBorderStyle = 1
const blackRGB = "000000"
const redPastelRGB = "FF6961"
const yellowPastelRGB = "FDFD96"
const greenPastelRGB = "77DD77"
const lightBluePastelRGB = "CEE5ED"
const fullFillType = "pattern"
const fullFillPattern = 1
const centerVerticalAlignment = "center"

func New(path string, sheetName string, logger *log.Logger) Report {
	f := excelize.NewFile()
	rp := Report{Path: path, file: f, nextRow: 1, sheetName: sheetName, mutex: new(sync.Mutex), logger: logger}
	rp.createStyles()
	return rp
}

func (rp *Report) createStyles() {
	borderStyle := continuousBorderStyle
	border := []excelize.Border{{Color: blackRGB, Style: borderStyle, Type: left},
		{Color: blackRGB, Style: borderStyle, Type: right},
		{Color: blackRGB, Style: borderStyle, Type: bottom},
		{Color: blackRGB, Style: borderStyle, Type: top}}
	redStyle, err := rp.file.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: fullFillType, Color: []string{redPastelRGB}, Pattern: fullFillPattern},
		Alignment: &excelize.Alignment{Vertical: centerVerticalAlignment, WrapText: true},
		Border:    border,
	})
	if err != nil {
		rp.logger.Println("Error: Failed To Create Red Style %v", err)
	}
	greenStyle, err := rp.file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: fullFillType, Color: []string{greenPastelRGB}, Pattern: fullFillPattern}, Border: border,
	})
	if err != nil {
		rp.logger.Println("Error: Failed To Create Green Style %v", err)
	}
	yellowStyle, err := rp.file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: fullFillType, Color: []string{yellowPastelRGB}, Pattern: fullFillPattern}, Border: border,
	})
	if err != nil {
		rp.logger.Println("Error: Failed To Create Yellow Style %v", err)
	}
	headerStyle, err := rp.file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: fullFillType, Color: []string{lightBluePastelRGB}, Pattern: fullFillPattern}, Border: border,
	})
	if err != nil {
		rp.logger.Println("Error: Failed To Create Header Style %v", err)
	}
	rp.styleNumbers = reportStyles{success: greenStyle, notTested: yellowStyle, failure: redStyle, headers: headerStyle}
}

type column struct {
	name  string
	width float64
}

func (rp *Report) AddHeaders() {
	rp.mutex.Lock()
	defer rp.mutex.Unlock()
	char := 'A'
	fields := []column{
		{"Repository", 20.0},
		{"Package Name", 60.0},
		{"Was Tested", 15},
		{"Was Skipped", 15},
		{"Time", 15},
		{"Number Of Runs", 15},
		{"Number of Failures", 18},
		{"Failure Message", 160},
	}
	for _, field := range fields {
		cell := string(char) + "1"
		err := rp.file.SetCellValue(rp.sheetName, cell, field.name)
		if err != nil {
			rp.logger.Println("Error: Failed To Set Cell %s in Header %v", field.name, err)
		}
		err = rp.file.SetColWidth(rp.sheetName, string(char), string(char), field.width)
		if err != nil {
			rp.logger.Println("Error: Failed To Set Cell Width Of %s in Header %v", field.name, err)
		}
		char = char + 1
	}
	err := rp.file.SetCellStyle(rp.sheetName, "A"+strconv.Itoa(rp.nextRow), string('A'+len(fields))+strconv.Itoa(rp.nextRow), rp.styleNumbers.headers)
	if err != nil {
		rp.logger.Println("Error: Failed To Set Style For Headers %v", err)
	}
	rp.nextRow += 1
}

const maxHeightCell = 409

func (rp *Report) AddTest(ts ITestSummary) {
	rp.mutex.Lock()
	defer rp.mutex.Unlock()
	var style int = rp.styleNumbers.success
	row := strconv.Itoa(rp.nextRow)
	values := []string{ts.GetRepoName(), ts.GetPackageName(), strconv.FormatBool(ts.WasTested()), strconv.FormatBool(ts.WasSkipped()),
		ts.GetTestTime(), strconv.Itoa(ts.GetNumberOfRuns()), strconv.Itoa(ts.GetNumberOfFailures()), ts.GetFailureMessage()}
	for i, value := range values {
		err := rp.file.SetCellValue(rp.sheetName, string('A'+i)+row, value)
		if err != nil {
			rp.logger.Println("Error: Failed To Set Cell Value For Row - (Pacakge Name: %s , Row Number: %+d  , Value: %s) %v", ts.GetPackageName(), rp.nextRow, value, err)
		}
	}

	if ts.HasFailed() {
		style = rp.styleNumbers.failure
		height := 12 * float64(len(strings.Split(ts.GetFailureMessage(), "\n")))
		if height > maxHeightCell {
			height = maxHeightCell
		}
		err := rp.file.SetRowHeight(rp.sheetName, rp.nextRow, height)
		if err != nil {
			rp.logger.Println("Error: Failed To Set Height For Row - (Pacakge Name: %s , Row Number: %+d ) %v", ts.GetPackageName(), rp.nextRow, err)
		}
	}

	if ts.WasSkipped() || !ts.WasTested() {
		style = rp.styleNumbers.notTested
	}

	err := rp.file.SetCellStyle(rp.sheetName, "A"+row, string('A'+len(values))+row, style)
	if err != nil {
		rp.logger.Println("Error: Failed To Set Style For Row - (Pacakge Name: %s , Row Number: %+d  %v", ts.GetPackageName(), rp.nextRow, err)
	}
	rp.nextRow++

}

func (rp *Report) Save() {
	rp.mutex.Lock()
	defer rp.mutex.Unlock()
	err := rp.file.SaveAs(rp.Path)
	if err != nil {
		rp.logger.Println("Error: Failed To Save Report: %v", err)
	}
}
