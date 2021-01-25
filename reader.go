package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// To read the specified path with filenames matching "*.xlsx" returned.
//
//  path string // path to check out
//
func readPath(
	path string,

) []string {
	fileNames, err := filepath.Glob(filepath.Join(path, "*.xls*"))
	if err != nil {
		panic(err)
	} else if len(fileNames) == 0 {
		panic(fmt.Errorf("XLSX files not found in %v", path))
	}

	return fileNames
}

// To get all staff names from filenames.
//
//  fileNames []string // filenames to check out
//
func getStaffNames(
	fileNames []string,

) []string {
	var staffNames []string
	var fileNameSplit []string

	for _, fileName := range fileNames {
		fileNameSplit = strings.Split(string(fileName), string(os.PathSeparator))
		staffNameSplit := strings.Split(fileNameSplit[len(fileNameSplit)-1], ".")
		staffName := strings.Join(staffNameSplit[:len(staffNameSplit)-1], "")
		staffNames = append(staffNames, staffName)
	}

	return staffNames
}

// To get the timetable of the specified staff.
//
//  name string // name of the specified staff
//  staffs map[string][]string // map[staff][period1, period2...]
//  periods map[string][]string // map[period][staff1, staff2...]
//  availableTimetableTemplatePositions map[string]string // positions to render in available timetable template
//  configurations map[string]interface{} // configuration map
//
func getStaffTimetable(
	name string,
	staffs map[string][]string,
	periods map[string][]string,
	availableTimetableTemplatePositions map[string][]string,
	configurations map[string]interface{},

) {
	sheetName := configurations["__sheet_name"].(string)
	availableTag := configurations["__available_tag"].(string)
	availableTimetablePath := configurations["__available_timetable_path"].(string)

	workbook, err := excelize.OpenFile(filepath.Join(availableTimetablePath, name+".xlsx"))
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

	for period, position := range availableTimetableTemplatePositions {
		tag, err := workbook.GetCellValue(sheetName, position[0])
		if err != nil {
			panic(err)
		}
		if tag == availableTag {
			mutex.Lock()
			staffs[name] = append(staffs[name], period)
			periods[period] = append(periods[period], name)
			mutex.Unlock()
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

	waitgroup.Done()

	fmt.Printf("Done with timetable for %v\n", name)
}
