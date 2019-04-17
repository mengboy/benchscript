package script

import (
	"github.com/montanaflynn/stats"
	"log"
	"sync"
	"time"
)

func BenchScript(totalRequest int, clientNum int, f func() error) {
	if totalRequest%clientNum != 0 {
		panic("total must be a mutiple of n")
	}

	time.Sleep(time.Second)

	every := totalRequest / clientNum
	var wg sync.WaitGroup

	consume := make([][]float64, clientNum)
	success := make([]int, clientNum)
	for i := range consume {
		consume[i] = make([]float64, 0, every)
	}
	start := time.Now()
	log.Println("start send requests at", start)
	//rest, _ := pingClient.QueryMethod(context.Background(), &servicePb.MethodRequest{Service: "ofo.ping.v1.PingService", Method: "ping", Body: "{\"token\":\"ping\"}"})
	//log.Infoln(rest)
	//time.Sleep(1000 * time.Second)
	for i := 0; i < clientNum; i++ {
		wg.Add(1)
		go func(i int) {
			for j := 0; j < every; j++ {
				start := time.Now()
				if err := f(); err != nil {

				} else {
					//log.Infoln("success")
					success[i]++
					end := time.Now()
					//consume[i][j] = float64(end.Sub(start).Nanoseconds())
					consume[i] = append(consume[i], float64(end.Sub(start).Nanoseconds()))
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	end := time.Now()
	log.Println("finish all requests at", end)
	t := end.Sub(start).Nanoseconds()

	successCount := 0
	for i := range success {
		successCount += success[i]
	}
	// merge result
	result := make([]float64, 0, totalRequest)
	for i := range consume {
		for _, t := range consume[i] {
			result = append(result, t)
		}
	}
	qps := int64(successCount) * int64(time.Second) / t
	mean, err := stats.Mean(result)
	panicIfError(err)
	median, err := stats.Median(result)
	panicIfError(err)
	max, err := stats.Max(result)
	panicIfError(err)
	min, err := stats.Min(result)
	panicIfError(err)

	ms := float64(time.Millisecond)
	us := float64(time.Microsecond)
	log.Printf("total send request: %d, success request: %d\n", totalRequest, successCount)
	log.Printf("client number: %d\n", clientNum)
	log.Printf("qps: %d\n", qps)
	log.Printf("mean: %f us, median: %f us, max: %f us, min: %f us\n", mean/us, median/us, max/us, min/us)
	log.Printf("mean: %f ms, median: %f ms, max: %f ms, min: %f ms\n", mean/ms, median/ms, max/ms, min/ms)
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
