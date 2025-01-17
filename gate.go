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
