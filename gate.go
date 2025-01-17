package main

import (
	"fmt"
	"log"
)

func selectAllGates(c *GateConfig, year int) []string {
	db := GetDatabase()
	query := fmt.Sprintf("SELECT g.%s FROM %s.%s g WHERE g.%s = %d AND g.APLCTN_CYC = 'Y';", c.GateNameKey, c.Dbname, c.TableName, c.GateYearKey, year)
	fmt.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
	}
	gates := []string{}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err.Error())
		}
		gates = append(gates, name)
	}
	fmt.Println(gates)
	return gates
}
