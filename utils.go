package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// To get the modified time of the file.
//
//  path string // path to the file
//
func getFileModTime(
	path string,

) int64 {
	file, err := os.Open(path)

	// If the file doesn't exist, return 0.

	if err != nil {
		return 0
	}
	defer file.Close()

	fileStat, err := file.Stat()

	// If anything goes wrong, return 0.

	if err != nil {
		return 0
	}
	return fileStat.ModTime().Unix()
}

// To get the cache (JSON files).
//
//  cacheFilename string // name of the cache JSON file
//  targetMap map[string][]string // the target map
//  configurations map[string]interface{} // configuration map
//
func getCache(
	cacheFilename string,
	targetMap map[string][]string,
	configurations map[string]interface{},

) {
	cachePath := configurations["__cache_path"].(string)

	cacheFile, err := os.Open(filepath.Join(cachePath, cacheFilename))
	if err != nil {
		panic(err)
	}
	defer cacheFile.Close()

	bytes, err := ioutil.ReadAll(cacheFile)
	if err != nil {
		panic(err)
	}

	json.Unmarshal([]byte(bytes), &targetMap)
	waitgroup.Done()
}

// To write map data into cache (JSON files).
//
//  cacheFilename string // name of the cache JSON file
//  targetMap map[string][]string // the target map
//  configurations map[string]interface{} // configuration map
//
func writeCache(
	cacheFilename string,
	targetMap map[string][]string,
	configurations map[string]interface{},

) {
	cachePath := configurations["__cache_path"].(string)

	bytes, _ := json.MarshalIndent(targetMap, "", "  ")
	err := ioutil.WriteFile(
		filepath.Join(cachePath, cacheFilename),
		bytes,
		os.ModeAppend,
	)
	if err != nil {
		panic(err)
	}
	waitgroup.Done()
}

// To generate an example settings.json.
//
//  configurations map[string]interface{} // configurations map
//
func generateExampleSettings(
	configurations map[string]interface{},

) {
	bytes, _ := json.MarshalIndent(configurations, "", "  ")
	err := ioutil.WriteFile(
		"settings.json",
		bytes,
		os.ModeAppend,
	)
	if err != nil {
		panic(err)
	}
}

// To unmerge all merged-cells in the specified sheet from the specified workbook,
// with merged-cells as []excelize.MergeCell and errors returned.
//
//  workbook *excelize.File // file struct of the specified workbook
//  sheetName string // name of the specified sheet
//
func unmergeAllCells(
	workbook *excelize.File,
	sheetName string,

) ([]excelize.MergeCell, error) {
	mergeCells, err := workbook.GetMergeCells(sheetName)
	if err != nil {
		return mergeCells, err
	}

	var mergeCellsSplit []string
	for _, mergeCell := range mergeCells {
		mergeCellsSplit = strings.Split(mergeCell[0], ":")
		err = workbook.UnmergeCell(
			sheetName,
			mergeCellsSplit[0],
			mergeCellsSplit[1],
		)
		if err != nil {
			return mergeCells, err
		}
	}

	return mergeCells, err
}

// To merge merged-cells given in the specified sheet from the specified workbook,
// with errors returned.
//
//  mergeCells []excelize.MergeCell // the merged-cells given
//  workbook *excelize.File // file struct of the specified workbook
//  sheetName string // name of the specified sheet
//
func mergeAllCells(
	mergeCells []excelize.MergeCell,
	workbook *excelize.File,
	sheetName string,

) error {
	for _, mergeCell := range mergeCells {
		mergeCellsSplit := strings.Split(mergeCell[0], ":")
		err := workbook.MergeCell(
			sheetName,
			mergeCellsSplit[0],
			mergeCellsSplit[1],
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// To skin the cellString with specified pattern.
// If the string matches te pattern, return the wrapped content.
// Or else return empty string.
//
//  cellString string // string to be skinned
//
func skinCellString(
	cellString string,
	configurations map[string]interface{},

) string {
	reg := regexp.MustCompile(configurations["__template_pattern"].(string))
	result := reg.FindAllStringSubmatch(cellString, -1)
	if len(result) > 0 {
		return result[0][1]
	}
	return ""
}
