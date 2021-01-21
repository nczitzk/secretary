package main

import (
	"math/rand"
	"sort"
	"time"
)

type periodStruct struct {
	period string
	number int
}

type periodSlice []periodStruct

func (slice periodSlice) Len() int {
	return len(slice)
}

func (slice periodSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice periodSlice) Less(i, j int) bool {
	return slice[i].number < slice[j].number
}

// To get limitation on the number.
//  key string // key to the number
//  target string // some specified target
//  configurations map[string]interface{} // configuration map
//
func getLimitaionOnNumber(
	key string,
	target string,
	configurations map[string]interface{},

) int {
	limitFloat, ok := configurations[target+key].(float64)
	if !ok {
		limitFloat = configurations[key].(float64)
	}
	return int(limitFloat)
}

// To confirm if the staff is assignable.
//  name string // staff's name
//  staffShifts map[string]map[int]int // all staffs'shifts
//  configurations map[string]interface{} // configuration map
//
func isStaffAssignable(
	name string,
	staffShifts map[string]map[int]int,
	configurations map[string]interface{},

) bool {
	shiftsPerStaff := getLimitaionOnNumber("__shifts_per_staff", name, configurations)

	shiftNumber, ok := staffShifts[name][0]
	if !ok {
		shiftNumber = 0
		staffShifts[name] = make(map[int]int)
		staffShifts[name][0] = 0
	}

	// According to labor law...

	if shiftNumber < shiftsPerStaff {
		return true
	}
	return false
}

// To schedule shifts for all staffs.
//  staffs map[string][]string // map[staff][period1, period2...]
//  periods map[string][]string // map[period][staff1, staff2...]
//  configurations map[string]interface{} // configuration map
//
func scheduler(
	staffs map[string][]string,
	periods map[string][]string,
	output map[string]map[int]string,
	configurations map[string]interface{},

) {
	// staffShifts[staff] = [
	//   0: 3  // number of shifts for the staff
	//   1: 1  // the 1st position has been assigned
	//   2: 1  // the 2nd position has been assigned
	//   3: 0  // the 3rd position has never been assigned
	//   4: 1  // the 4th position has been assigned
	//   5: 0  // the 5th position has never been assigned
	//   ...
	// ]
	staffShifts := make(map[string]map[int]int)

	tryToAssignStaffsToDifferentPositions :=
		configurations["__try_to_assign_staffs_to_different_positions"].(bool)

	var thePeriodSlice periodSlice
	for period, names := range periods {
		thePeriodSlice = append(thePeriodSlice, periodStruct{
			period: period,
			number: len(names),
		})
	}

	sort.Stable(thePeriodSlice)

	for _, thisPeriod := range thePeriodSlice {

		period := thisPeriod.period
		names := periods[period]

		output[period] = make(map[int]string)

		// Shuffle the names.

		rand.Seed(time.Now().Unix())
		rand.Shuffle(len(names), func(i, j int) {
			names[i], names[j] = names[j], names[i]
		})
		staffsPerShift := getLimitaionOnNumber("__staffs_per_shift", period, configurations)

		index := 1

		if tryToAssignStaffsToDifferentPositions {

			// Make staffs who has never been assigned to No.index position step to the front.

			for i, name := range names {

				// If the staff has been assigned to No.index position,
				// pop his/her name and add it to the end.

				if staffShifts[name][index] == 1 {
					names = append(names[:i], names[i+1:]...)
					names = append(names, name)
				}
			}
		}

		for _, name := range names {

			// If the shift is full already,
			// there is no need to deal with the shift.

			if index > staffsPerShift {
				break
			}

			if isStaffAssignable(
				name,
				staffShifts,
				configurations,
			) {
				// No.index position has been assigned to staff.

				staffShifts[name][index] = 1

				// Add 1 to the number of shifts for the staff.

				staffShifts[name][0]++

				output[period][index] = name
				index++
			}
		}
	}
}
