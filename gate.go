package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Gate struct {
	Id       uint64 // Unique identifier for the gate
	GateName string // Name of the gate
	GateYear string // Year the gate is associated with
	Start    string // Start datetime (string format)
	End      string // End datetime (string format)
}

const createdFormat = "2006-01-02 15:04:05" // Time format for DB
const DAYS_PER_WINDOW = 3                   // Number of days per gate window

// Returns all years with applicable gates from the database
func selectAllYears(c *GateConfig) []string {
	db := GetDatabase()
	query := fmt.Sprintf("SELECT DISTINCT g.%s FROM %s.%s g WHERE g.%s = 'Y';", c.GateYearKey, c.Dbname, c.TableName, c.GateIsApplicableFlag)
	log.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err.Error())
	}
	years := []string{}
	defer rows.Close()
	for rows.Next() {
		var year string
		if err := rows.Scan(&year); err != nil {
			log.Fatal(err.Error())
		}
		years = append(years, year)
	}
	return years
}

// Returns all gates for a given year, ordered by their configured order
func selectAllGates(c *GateConfig, year int) []*Gate {
	db := GetDatabase()
	query := fmt.Sprintf("SELECT g.%s, g.%s, g.%s FROM %s.%s g WHERE g.%s = %d AND g.%s = 'Y' ORDER BY g.%s ASC;", c.GateNameKey, c.StartKey, c.EndKey, c.Dbname, c.TableName, c.GateYearKey, year, c.GateIsApplicableFlag, c.GateOrderKey)
	fmt.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
	}
	gates := []*Gate{}
	defer rows.Close()
	for rows.Next() {
		var gate = &Gate{}
		if err := rows.Scan(&gate.GateName, &gate.Start, &gate.End); err != nil {
			log.Fatal(err.Error())
		}
		gates = append(gates, gate)
	}
	return gates
}

// Sets all gates' start/end times relative to a selected gate and position
func setGatesRelativeTo(c *GateConfig, year int, gate string, pos RelativePosition) {
	var wg sync.WaitGroup
	gates := selectAllGates(c, year) // Get all gates for the year
	queries := _setGatesRelativeTo(c, gates, time.Now(), year, gate, pos)
	for _, q := range queries {
		wg.Add(1)
		log.Println(q)
		go func() {
			defer wg.Done()
			ExecSql(q)
		}()
	}
	wg.Wait()
}

// Helper: builds update queries for all gates relative to the selected gate
func _setGatesRelativeTo(c *GateConfig, gates []*Gate, now time.Time, year int, gate string, pos RelativePosition) []string {
	index := -1 // Find index of the selected gate
	for i, g := range gates {
		if g.GateName == gate {
			index = i
			break
		}
	}
	if index < 0 {
		log.Fatal("Specified gate not found")
	}
	queryStrings := []string{}
	for i := 0; i < len(gates); i++ {
		pastGate := gates[i]
		s := createQueryString(c, now, year, pastGate.GateName, i-index, pos)
		queryStrings = append(queryStrings, s)
	}
	return queryStrings
}

// Builds a single SQL update query for a gate based on its relative position
func createQueryString(c *GateConfig, now time.Time, year int, pastGate string, magnitude int, pos RelativePosition) string {
	return _createQueryString(c, now, year, pastGate, magnitude, pos)
}

// Internal: builds the SQL update string for a gate
func _createQueryString(c *GateConfig, now time.Time, year int, pastGate string, magnitude int, pos RelativePosition) string {
	halfday := time.Hour * time.Duration(12)
	var date time.Time
	if magnitude != 0 {
		date = now.Add(time.Hour * 24 * time.Duration(magnitude) * DAYS_PER_WINDOW)
	} else {
		date = now
	}
	start := date.Add(-halfday).Format(createdFormat)
	end := date.Add(halfday).Format(createdFormat)
	if magnitude == 0 && pos != INSIDE {
		shift := time.Hour * 6
		if pos == BEFORE {
			start = date.Add(shift).Format(createdFormat)
		}
		if pos == AFTER {
			end = date.Add(-shift).Format(createdFormat)
		}
	}
	return fmt.Sprintf("UPDATE %s.%s SET %s = '%s', %s = '%s' WHERE %s = '%s' AND %s = %d;", c.Dbname, c.TableName, c.StartKey, start, c.EndKey, end, c.GateNameKey, pastGate, c.GateYearKey, year)
}

// Checks if all gates are in a strictly linear timeline (no overlaps)
func isTimelineLinear(gates []*Gate) bool {
	for i := 1; i < len(gates); i++ {
		currentStart, _ := time.Parse(createdFormat, gates[i].Start)
		previousEnd, _ := time.Parse(createdFormat, gates[i-1].End)
		if !currentStart.After(previousEnd) {
			return false
		}
	}
	return true
}

// Returns true if the gate is currently open (now is between start and end)
func isGateOpen(g *Gate) bool {
	now := time.Now()
	s, _ := time.Parse(createdFormat, g.Start)
	e, _ := time.Parse(createdFormat, g.End)
	return s.Before(now) && e.After(now)
}

// Returns -1 if gate is in the past, 0 if open, 1 if in the future
func getGatePosition(g *Gate) int {
	now := time.Now()
	s, _ := time.Parse(createdFormat, g.Start)
	e, _ := time.Parse(createdFormat, g.End)
	if s.Before(now) && e.After(now) {
		return 0
	}
	if e.Before(now) {
		return -1
	}
	if s.After(now) {
		return 1
	}
	panic("Unexpected gate status")
}
