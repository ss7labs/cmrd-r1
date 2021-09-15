package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"sync/atomic"
	"time"
)

var tsdb *sql.DB
var db *sql.DB

type CallRecord struct {
	EventTime time.Time
	NumberA   string
	NumberB   string
	Duration  uint16
	TrunkIn   string
	TrunkOut  string
}

func (e *Env) trunkTotalByDestination() {
	var err error
	query := "SELECT event_time,numa,numb,duration,t_in_name,t_out_name FROM cmrd.amts WHERE toYYYYMM(event_date)='202108' AND t_in_name='BTMSL2'"
	rows, err := tsdb.Query(query)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		var record CallRecord
		err = rows.Scan(&record.EventTime, &record.NumberA, &record.NumberB, &record.Duration, &record.TrunkIn, &record.TrunkOut)
		if err != nil {
			panic(err.Error())
		}
		go e.rateRecord(record)
	}
}

func (e *Env) rateRecord(record CallRecord) {

	if match, _ := regexp.MatchString(`^(8101|8107)`, record.NumberB); match {
		count <- 1
		return
	}

	//var prefix string
	for _, v := range e.PrefixOrdered {

		patt := "^" + v + ".+"
		if match, _ := regexp.MatchString(patt, record.NumberB); match {
			//prefix = v
			break
		}
	}
	//shortName := e.PrefixShortname[prefix]
	//fmt.Println(record.NumberB,prefix,shortName)
	count <- 1
	//return shortName
}

const (
	flushInterval = time.Duration(1) * time.Second
)

var ops uint64 = 0
var total uint64 = 0
var count chan uint64
var flushTicker *time.Ticker

func initCounter() {
	count = make(chan uint64)
	flushTicker = time.NewTicker(flushInterval)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			atomic.AddUint64(&total, ops)
			fmt.Printf("Total ops %d\n", total)
			os.Exit(0)
		}
	}()

	go func(count chan uint64) {
		for {
			_ = <-count
			atomic.AddUint64(&ops, 1)
		}
	}(count)

	go func() {
		for range flushTicker.C {
			fmt.Printf("Ops/s %f\n", float64(ops)/flushInterval.Seconds())
			atomic.AddUint64(&total, ops)
			atomic.StoreUint64(&ops, 0)
		}
	}()

}
func (e *Env) initPrefix() {
	query := "SELECT prefix,IFNULL(shortname,0) FROM route_prices ORDER BY CHAR_LENGTH(prefix) DESC"
	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var prefix, shortName string
	var prefixShortname map[string]string = make(map[string]string)
	var prefixOrdered []string

	for rows.Next() {
		err = rows.Scan(&prefix, &shortName)
		if err != nil {
			panic(err.Error())
		}
		prefixOrdered = append(prefixOrdered, prefix)
		prefixShortname[prefix] = shortName
	}
	e.PrefixOrdered = prefixOrdered
	e.PrefixShortname = prefixShortname

}
