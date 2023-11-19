package loader

import (
	"encoding/json"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"load-tester-backend/modules/db"
	"load-tester-backend/modules/model"
	"log"
	"net/http"
	"time"
)

func StartLoadTesting(requestBody model.StartLoadRequest) {
	for _, val := range requestBody.StepConfig {
		startStep(val.Rate, val.DurationSeconds, requestBody.Url, requestBody.TestBody)
	}
}

func getBody(TestBody model.TestBody) []byte {
	byteBody, err := json.Marshal(TestBody)

	if err != nil {
		log.Println("Error while convert body to byte array")
		return nil
	}
	return byteBody
}

func newCustomTargeter(url string, testBody model.TestBody) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "POST"
		headers := make(http.Header)
		headers.Set("Content-Type", "application/json")
		tgt.Header = headers
		tgt.URL = url

		tgt.Body = getBody(testBody)
		//добавить логирование

		return nil
	}
}

func startStep(loadRate int, timeSeconds int, url string, testBody model.TestBody) {
	rate := vegeta.Rate{Freq: loadRate, Per: time.Second}
	duration := time.Duration(timeSeconds) * time.Second

	targeter := newCustomTargeter(url, testBody)

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics

	for res := range attacker.Attack(targeter, rate, duration, "Big Big Bang!") {
		log.Printf("Making request. Code %d, Error: %s, Latencies: %s", res.Code, res.Error, res.Latency)
		metrics.Add(res)
	}
	metrics.Close()
	log.Printf("\n\n==========\nTest ended. Latencies 99 percentiles: %s\n"+
		"Latencies 95 percentiles: %s\n"+
		"Latencies minimum: %s\n"+
		"Latencies maximum: %s\n"+
		"Rate: %f\n"+
		"Duration: %s\n"+
		"Throughput: %f\n"+
		"Requests total number: %d\n"+
		"Success count percentage: %f",
		metrics.Latencies.P99, metrics.Latencies.P95, metrics.Latencies.Min, metrics.Latencies.Max,
		metrics.Rate, metrics.Duration, metrics.Throughput, metrics.Requests, metrics.Success)
	//сохранение в бд
	db.SaveResults(metrics.Latencies.P99, metrics.Latencies.P95, metrics.Latencies.Min, metrics.Latencies.Max,
		metrics.Rate, metrics.Duration, metrics.Throughput, metrics.Requests, metrics.Success)
}
