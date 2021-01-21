package main

import (
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// To print the available timetable.
//
//  staffs map[string][]string // map[staff][period1, period2...]
//  periods map[string][]string // map[period][staff1, staff2...]
//  availableTimetableTemplatePositions map[string]string // positions to render in available timetable template
//  configurations map[string]interface{} // configuration map
//
func printavailableTimetable(
	staffs map[string][]string,
	periods map[string][]string,
	availableTimetableTemplatePositions map[string][]string,
	configurations map[string]interface{},

) {
	sheetName := configurations["__sheet_name"].(string)
	dateFormat := configurations["__date_format"].(string)
	joinSeparator := configurations["__join_separator"].(string)
	title := configurations["__available_timetable_title"].(string)
	availableTimetableTitle := configurations["__available_timetable_title"].(string)
	availableTimetableTemplatePath := configurations["__available_timetable_template_path"].(string)

	workbook, err := excelize.OpenFile(availableTimetableTemplatePath)
	if err != nil {
		panic(err)
	}

	// To unmerge all merged-cells at first.

	mergeCells, err := unmergeAllCells(
		workbook,
		sheetName,
	)
	if err != nil {
		panic(err)
	}

	for period, positions := range availableTimetableTemplatePositions {
		for _, position := range positions {
			if period == "__title" {
				workbook.SetCellValue(
					sheetName,
					position,
					title,
				)
			} else if period == "__date" {
				workbook.SetCellValue(
					sheetName,
					position,
					time.Now().Format(dateFormat),
				)
			} else {
				workbook.SetCellValue(
					sheetName,
					position,
					strings.Join(periods[period], joinSeparator),
				)
			}
		}
	}

	// To remerge all unmerged merged-cells just now.

	err = mergeAllCells(
		mergeCells,
		workbook,
		sheetName,
	)
	if err != nil {
		panic(err)
	}

	// The output filename will be xxx-2006-01-02.xlsx.
	// p.s. xxx refers to __availableTimetableTitle in configurations map
	//      2006-01-02 as date format refers to __date_format in configurations map

	err = workbook.SaveAs(
		availableTimetableTitle + "-" + time.Now().Format(dateFormat) + ".xlsx",
	)
	if err != nil {
		panic(err)
	}
}

// To print the timetable.
//  output map[string]map[int]string // map[period][0, 1, 2...][staff1, staff2...]
//  timetableTemplatePositions map[string]string // positions to render in timetable template
//  configurations map[string]interface{} // configuration map
//
func printTimetable(
	output map[string]map[int]string,
	timetableTemplatePositions map[string][]string,
	configurations map[string]interface{},

) {
	sheetName := configurations["__sheet_name"].(string)
	title := configurations["__timetable_title"].(string)
	dateFormat := configurations["__date_format"].(string)
	timetableTitle := configurations["__timetable_title"].(string)
	timetableTemplatePath := configurations["__timetable_template_path"].(string)

	workbook, err := excelize.OpenFile(timetableTemplatePath)
	if err != nil {
		panic(err)
	}

	// To unmerge all merged-cells at first.

	mergeCells, err := unmergeAllCells(
		workbook,
		sheetName,
	)
	if err != nil {
		panic(err)
	}

	for period, positions := range timetableTemplatePositions {
		for index, position := range positions {
			if index < len(output[period]) {

				// When the cells for this period are enough,
				// fill staffs'names in.

				workbook.SetCellValue(
					sheetName,
					position,
					output[period][index+1],
				)
			} else if period == "__title" {
				workbook.SetCellValue(
					sheetName,
					position,
					title,
				)
			} else if period == "__date" {
				workbook.SetCellValue(
					sheetName,
					position,
					time.Now().Format(dateFormat),
				)
			} else {

				// The blank cells will be filled with empty string.

				workbook.SetCellValue(
					sheetName,
					position,
					"",
				)
			}
		}
	}

	// To remerge all unmerged merged-cells just now.

	err = mergeAllCells(
		mergeCells,
		workbook,
		sheetName,
	)
	if err != nil {
		panic(err)
	}

	// The output filename will be xxx-2006-01-02.xlsx.
	// p.s. xxx refers to __availableTimetableTitle in configurations map
	//      2006-01-02 as date format refers to __date_format in configurations map

	err = workbook.SaveAs(
		timetableTitle + "-" + time.Now().Format(dateFormat) + ".xlsx",
	)
	if err != nil {
		panic(err)
	}
}
