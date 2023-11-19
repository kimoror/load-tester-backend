package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() {
	log.Println("Init db connection...")
	host := os.Getenv("DATABASE.HOST")
	port, _ := strconv.Atoi(os.Getenv("DATABASE.PORT"))
	dbname := os.Getenv("DATABASE.NAME")
	user := os.Getenv("DATABASE.USER")
	password := os.Getenv("DATABASE.PASSWORD")
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, dbname, user, password)

	var errConnection error
	db, errConnection = sql.Open("postgres", connStr)

	if errConnection != nil {
		log.Fatalf("Initilizing DB failed: %s", errConnection)
	}

	err := db.Ping()

	checkDbError(err)

	log.Println("Db init completed")
}

func SaveResults(latenciesP99 time.Duration, latenciesP95 time.Duration,
	latenciesMin time.Duration, latenciesMax time.Duration, rate float64, duration time.Duration,
	throughput float64, requestsTotalCount uint64, requestsSuccessPercentage float64) {
	insertResult := `INSERT INTO load_tester.load_tests_results(latencies_p99, latencies_p95, latencies_min, latencies_max, rate, duration, throughput, requests_total_count, requests_success_percentage) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	testId := -1
	err := db.QueryRow(insertResult, latenciesP99.String(), latenciesP95.String(), latenciesMin.String(), latenciesMax.String(), rate, duration.String(), throughput, requestsTotalCount, requestsSuccessPercentage).Scan(&testId)
	checkDbError(err)
}

func checkDbError(err error) {
	if err != nil {
		log.Panicln("Db error", err)
	}
}
