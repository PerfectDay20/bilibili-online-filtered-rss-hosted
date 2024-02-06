package main

import (
	"fmt"
	"testing"
	"time"
)

func TestDynamoGet(t *testing.T) {
	table := initTable()
	//table.SetRecord("test unix time")
	record, _ := table.GetRecord()
	fmt.Println(record)
	parsedTime := time.Unix(record.UpdateTimestamp, 0)
	fmt.Println(parsedTime)
	duration := time.Now().Sub(parsedTime)
	fmt.Println(duration)
	if duration > 5*time.Minute {
		fmt.Println("over 5 minutes")
	}
	if duration < 10*time.Minute {
		fmt.Println("less 10 minutes")
	}
}
