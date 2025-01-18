package main

import (
	"fmt"
	"log"
	"time"
)

type Gate struct {
	GateName string
	GateYear string
	Start    string
	End      string
}

const DAYS_PER_WINDOW = 3
const createdFormat = "2006-01-02 15:04:05" //"Jan 2, 2006 at 3:04pm (MST)"

func selectAllGates(c *GateConfig, year int) []*Gate {
	db := GetDatabase()
	query := fmt.Sprintf("SELECT g.%s FROM %s.%s g WHERE g.%s = %d AND g.APLCTN_CYC = 'Y' ORDER BY g.SRT_ORDR ASC;", c.GateNameKey, c.Dbname, c.TableName, c.GateYearKey, year)
	fmt.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
	}
	gates := []*Gate{}
	defer rows.Close()
	for rows.Next() {
		var gate = &Gate{}
		if err := rows.Scan(&gate.GateName); err != nil {
			log.Fatal(err.Error())
		}
		gates = append(gates, gate)
	}
	return gates
}

func updateGate(c *GateConfig, gate string, year int, start, end time.Time) {
	db := GetDatabase()
	query := fmt.Sprintf("UPDATE %s.%s SET %s = %s, %s = %s WHERE %s = %s AND %s = %d;", c.Dbname, c.TableName, c.StartKey, start, c.EndKey, end, c.GateNameKey, gate, c.GateYearKey, year)
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func setGatesRelativeTo(c *GateConfig, year int, gate string, pos int) []string {
	gates := selectAllGates(c, year) // these are ordered
	queries := _setGatesRelativeTo(c, gates, time.Now(), year, gate, pos)
	return queries
	// for _, s := range queries {
	// 	fmt.Println(s)
	// 	fmt.Println()
	// }
}
func _setGatesRelativeTo(c *GateConfig, gates []*Gate, now time.Time, year int, gate string, pos int) []string {
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
		s := createQueryString(c, now, year, pastGate.GateName, i-index)
		queryStrings = append(queryStrings, s)
	}
	return queryStrings
}
func createQueryString(c *GateConfig, now time.Time, year int, pastGate string, magnitude int) string {
	return _createQueryString(c, now, year, pastGate, magnitude)
}
func _createQueryString(c *GateConfig, now time.Time, year int, pastGate string, magnitude int) string {
	halfday := time.Hour * time.Duration(12)
	var date time.Time
	switch {
	case magnitude < 0:
		date = now.Add(time.Hour * 24 * time.Duration(magnitude) * DAYS_PER_WINDOW)
		break
	case magnitude > 0:
		date = now.Add(time.Hour * 24 * time.Duration(magnitude) * DAYS_PER_WINDOW)
		break
	default:
		date = now
	}
	start := date.Add(-halfday).Format(createdFormat)
	end := date.Add(halfday).Format(createdFormat)
	return fmt.Sprintf("UPDATE %s.%s SET %s = '%s', %s = '%s' WHERE %s = '%s' AND %s = %d;", c.Dbname, c.TableName, c.StartKey, start, c.EndKey, end, c.GateNameKey, pastGate, c.GateYearKey, year)
}
