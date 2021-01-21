package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// To load default configurations.
//
// 	configurations map[string]interface{} // configuration map
//
func loadDefaultConfiguration(
	configurations map[string]interface{},

) {
	// Name of sheet

	configurations["__sheet_name"] = "Sheet1"

	// Date format for timetable filename and {{__date}} in template
	//
	//   Month: 1, 01, Jan, January
	//   Date: 2, 02, _2
	//   Hour: 3, 03, 15, PM, pm, AM, am
	//   Minute: 4, 04
	//   Second: 5, 05
	//   Year: 06, 2006
	//   Timezone: -07, -0700, Z0700, Z07:00, -07:00, MST
	//   Day: Mon, Monday
	//
	// e.g. as `2021-01-01 08:00:00`,
	//      for "2006-01-02" => `2021-01-01`
	//      for "Jan 02, 2006" => `Jan 01 2021`
	//      for "3:02 02/01/2006" => `8:00 01/01/2021`

	configurations["__date_format"] = "2006-01-02"

	// Separator between staff names in available timetable
	// e.g. If `, ` is set as __join_separator,
	//      `Ethan, John` will refer to Ethan and John.

	configurations["__join_separator"] = ", "

	// Template pattern for parsing cell strings
	// e.g. If `\\{\\{(.*?)\\}\\}` is set as ____template_pattern,
	//      `{{Monday:1}}`will be parsed to `Monday:1`.

	configurations["__template_pattern"] = "\\{\\{(.*?)\\}\\}"

	// When content in some cell in certain staff's available timetable is equal to  __available_tag,
	// then the time period which the cell refers to is available to this staff.

	configurations["__available_tag"] = "free"

	// Path to cache files such as __staffs.json

	configurations["__cache_path"] = ".cache"

	// Path to template files

	configurations["__template_path"] = "template"

	// Path to all staffs'available timetable

	configurations["__available_timetable_path"] = "available-timetable"

	// Path to the timetable

	configurations["__timetable_template_path"] = "template\\timetable-template.xlsx"

	// Path to the available timetable

	configurations["__available_timetable_template_path"] = "template\\available-timetable-template.xlsx"

	// Title of the timetable and {{__title}} in templates

	configurations["__timetable_title"] = "timetable"

	// Title of the available timetable and {{__title}} in templates

	configurations["__available_timetable_title"] = "available-timetable"

	// Try to assign staffs to different positions or not.
	// e.g. If Ethan has been assigned to the first position of some assigned shift,
	//      then Ethan will be secondary consideration to the first positions of the other shifts.

	configurations["__try_to_assign_staffs_to_different_positions"] = true

	// The number of shifts for one staff by default

	configurations["__shifts_per_staff"] = 4

	// The number of staffs for one shift by default

	configurations["__staffs_per_shift"] = 5

	// The number of shifts for the specified staff
	// e.g. If **Ethan** has to be on call in 10 shifts
	//      rather than 4 shifts by default as above,
	//      then the below should be set.

	// configurations["Ethan__shifts_per_staff"] = 10

	// The number of staffs for the specified shift
	// e.g. If there has to be 7 staffs on call for **the first shift on Monday**
	//      rather than 5 staffs by default as above,
	//      then the below should be set.

	// configurations["Mon:1__staffs_per_shift"] = 7

	configurations["Mon:1__staffs_per_shift"] = 4
	configurations["Tue:1__staffs_per_shift"] = 4
	configurations["Wed:1__staffs_per_shift"] = 4
	configurations["Thu:1__staffs_per_shift"] = 4
	configurations["Fri:1__staffs_per_shift"] = 4
	configurations["Sat:1__staffs_per_shift"] = 4
	configurations["Sun:1__staffs_per_shift"] = 4

	generateExampleSettings(configurations)
}

// To load configuration JSON file.
//
//  configurationFilePath string // path to configuration JSON file
//  configurations map[string]interface{} // configuration map
//
func loadConfiguration(
	configurationFilePath string,
	configurations map[string]interface{},

) {
	configurationFile, err := os.Open(configurationFilePath)

	// if anything goes wrong, use default configurations.

	if err != nil {
		loadDefaultConfiguration(configurations)
		return
	}
	defer configurationFile.Close()

	bytes, err := ioutil.ReadAll(configurationFile)

	// if anything goes wrong, use default configurations.

	if err != nil {
		loadDefaultConfiguration(configurations)
		return
	}

	json.Unmarshal([]byte(bytes), &configurations)
}

// To load all staffs'timetable.
//  staffs map[string][]string // map[staff][period1, period2...]
//  periods map[string][]string // map[period][staff1, staff2...]
//  availableTimetableTemplatePositions map[string]string // positions to render in available timetable template
//  configurations map[string]interface{} // configuration map
//
func loadStaffstimetable(
	staffs map[string][]string,
	periods map[string][]string,
	availableTimetableTemplatePositions map[string][]string,
	configurations map[string]interface{},

) {
	cachePath := configurations["__cache_path"].(string)
	templatesPath := configurations["__template_path"].(string)
	availableTimetablePath := configurations["__available_timetable_path"].(string)

	// If any files have been modified in path to available timetables,
	// compared the modified time of path to available timetables with the modified time of path to cache files.

	if getFileModTime(availableTimetablePath) > getFileModTime(cachePath) ||

		// If any files have been modified in path to template files,
		// compared the modified time of path to template files with the modified time of path to cache files.

		getFileModTime(templatesPath) > getFileModTime(cachePath) ||

		// If __staffs.json doesn't exist...

		getFileModTime(filepath.Join(cachePath, "__staffs.json")) == 0 ||

		// If __periods.json doesn't exist...

		getFileModTime(filepath.Join(cachePath, "__periods.json")) == 0 {

		fmt.Println("The staffs'timetables have been updated.")

		os.Remove(filepath.Join(cachePath, "__staffs.json"))
		os.Remove(filepath.Join(cachePath, "__periods.json"))

		for _, name := range getStaffNames(readPath(availableTimetablePath)) {
			fmt.Printf("Reading timetable for %v...\n", name)

			waitgroup.Add(1)
			go getStaffTimetable(
				name,
				staffs,
				periods,
				availableTimetableTemplatePositions,
				configurations,
			)
		}
		waitgroup.Wait()

		waitgroup.Add(1)
		go writeCache("__staffs.json", staffs, configurations)
		waitgroup.Add(1)
		go writeCache("__periods.json", periods, configurations)
		waitgroup.Wait()

	} else {
		waitgroup.Add(1)
		go getCache("__staffs.json", staffs, configurations)
		waitgroup.Add(1)
		go getCache("__periods.json", periods, configurations)
		waitgroup.Wait()
	}
}

// To load template for timetables.
//
//  templateName string // name of the template ("available_timetable" or "timetable")
//  templatePositions map[string]string // positions to render in available timetable template
//  configurations map[string]interface{} // configuration map
//
func loadTemplate(
	templateName string,
	templatePositions map[string][]string,
	configurations map[string]interface{},

) {
	cachePath := configurations["__cache_path"].(string)
	templateCachePath := filepath.Join(cachePath, "__"+templateName+".json")
	templatePath := configurations["__"+templateName+"_template_path"].(string)

	// If any files have been modified in path to certain template file,
	// compared the modified time of path to certain template file with the modified time of path to certain template cache file.

	if getFileModTime(templatePath) > getFileModTime(templateCachePath) ||

		// If the template cache doesn't exist...

		getFileModTime(templateCachePath) == 0 {

		os.Remove(filepath.Join(cachePath, "__"+templateName+".json"))

		workbook, err := excelize.OpenFile(templatePath)
		if err != nil {
			panic(err)
		}
		rows, err := workbook.GetRows(configurations["__sheet_name"].(string))
		if err != nil {
			panic(err)
		}
		for rowIndex, row := range rows {
			for colIndex, cellString := range row {

				// To skin the cellString with specified pattern.

				skinnedCellString := skinCellString(cellString, configurations)

				if skinnedCellString != "" {

					// 1 -> A, 2 -> B...
					// (1, 1) => A1, (2, 1) => B1...

					colName, err := excelize.ColumnNumberToName(colIndex + 1)
					if err != nil {
						panic(err)
					}

					templatePositions[skinnedCellString] = append(
						templatePositions[skinnedCellString],
						colName+fmt.Sprint(rowIndex+1),
					)
				}
			}
		}

		waitgroup.Add(1)
		go writeCache("__"+templateName+".json", templatePositions, configurations)

	} else {
		waitgroup.Add(1)
		go getCache("__"+templateName+".json", templatePositions, configurations)
		waitgroup.Wait()
	}
}
