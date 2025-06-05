package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Gate struct {
	Id       uint64
	GateName string
	GateYear string
	Start    string
	End      string
}

const createdFormat = "2006-01-02 15:04:05" //"Jan 2, 2006 at 3:04pm (MST)"
const DAYS_PER_WINDOW = 3

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

func setGatesRelativeTo(c *GateConfig, year int, gate string, pos RelativePosition) {
	var wg sync.WaitGroup
	gates := selectAllGates(c, year) // these are ordered
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
func _setGatesRelativeTo(c *GateConfig, gates []*Gate, now time.Time, year int, gate string, pos RelativePosition) []string {
	index := -1 // so lets find index of gate were working relative to
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
	for i := 0; i < len(gates); i++ { // set gates before selected gate
		pastGate := gates[i]
		s := createQueryString(c, now, year, pastGate.GateName, i-index, pos)
		queryStrings = append(queryStrings, s)
	}
	return queryStrings
}
func createQueryString(c *GateConfig, now time.Time, year int, pastGate string, magnitude int, pos RelativePosition) string {
	return _createQueryString(c, now, year, pastGate, magnitude, pos)
}
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
func isGateOpen(g *Gate) bool {
	now := time.Now()
	s, _ := time.Parse(createdFormat, g.Start)
	e, _ := time.Parse(createdFormat, g.End)
	return s.Before(now) && e.After(now)
}
