package main

import (
	"flag"
	"os"
	"runtime"
	"sync"
)

var (
	mutex     sync.Mutex
	waitgroup sync.WaitGroup
)

func main() {
	var (
		configurationFilePath string // path to configuration JSON file

		configurations = make(map[string]interface{}) // configuration map

		periods = make(map[string][]string)       // map[period][staff1, staff2...]
		staffs  = make(map[string][]string)       // map[staff][period1, period2...]
		output  = make(map[string]map[int]string) // map[period][0, 1, 2...][staff1, staff2...]

		timetableTemplatePositions          = make(map[string][]string) // positions to render in timeable template
		availableTimetableTemplatePositions = make(map[string][]string) // positions to render in available timeable template
	)

	runtime.GOMAXPROCS(runtime.NumCPU())

	// secretary

	// 	-c string
	// 	    	configuration JSON file (default "settings.json")
	// 	-h    to print help message
	// 	-r    to remove all caches and generate sheets
	// 	-s    to generate an example settings.json

	flag.StringVar(&configurationFilePath, "c", "settings.json", "configuration JSON file")

	needHelp := flag.Bool("h", false, "to print help message")
	needReset := flag.Bool("r", false, "to remove all caches and generate sheets")
	needExampleSettings := flag.Bool("s", false, "to generate an example settings.json")

	flag.Parse()

	tail := flag.Args()

	if len(tail) == 0 {
		tail = append(tail, configurationFilePath)
	}

	if *needHelp {
		flag.PrintDefaults()
		return
	}

	if *needExampleSettings {
		loadDefaultConfiguration(configurations)
		return
	}

	for _, configurationFile := range tail {

		// To load configuration JSON file.

		loadConfiguration(
			configurationFile,
			configurations,
		)

		cachePath := configurations["__cache_path"].(string)

		if *needReset {
			os.RemoveAll(cachePath)
		}

		os.Mkdir(cachePath, os.ModePerm)

		// To load available timeable template.

		loadTemplate(
			"available_timetable",
			availableTimetableTemplatePositions,
			configurations,
		)

		// To load timeable template.

		loadTemplate(
			"timetable",
			timetableTemplatePositions,
			configurations,
		)

		waitgroup.Wait()

		// To load all staffs'timetable.

		loadStaffsTimeable(
			staffs,
			periods,
			availableTimetableTemplatePositions,
			configurations,
		)

		// To print the available timetable.

		printavailableTimetable(
			staffs,
			periods,
			availableTimetableTemplatePositions,
			configurations,
		)

		// To schedule shifts for all staffs.

		scheduler(
			staffs,
			periods,
			output,
			configurations,
		)

		// To print the available timetable.

		printTimetable(
			output,
			timetableTemplatePositions,
			configurations,
		)

		waitgroup.Wait()
	}
}
