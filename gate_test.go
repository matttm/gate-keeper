package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ShouldCreateQueryString(t *testing.T) {
	type QueryStringTest struct {
		config        *GateConfig
		year          int
		pastGate      string
		magnitude     int
		expectedQuery string
	}
	table := []QueryStringTest{
		{
			config:        {},
			year:          2026,
			pastGate:      "",
			magnitude:     0,
			expectedQuery: "",
		},
	}
	for _, v := range table {
		q := createQueryString(v.config, v.year, v.pastGate, v.magnitude)
		assert.Equal(t, q, v.expectedQuery)
	}
}
