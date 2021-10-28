package speedtest

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gnum"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var dlSizes = [...]int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
var ulSizes = [...]int{100, 300, 500, 800, 1000, 1500, 2500, 3000, 3500, 4000} //kB
var client = http.Client{}

func (s *TestServer) testDownload() (float64, error) {
	if len(s.URL) == 0 {
		return 0, gerrors.New("Empty test server URL")
	}
	latency, err := s.testLatency()
	if err != nil {
		return 0, err
	}

	dlURL := strings.Split(s.URL, "/upload")[0]
	wg := new(sync.WaitGroup)

	// Warming up
	sTime := time.Now()
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go dlWarmUp(wg, dlURL)
	}
	wg.Wait()
	fTime := time.Now()
	// 1.125MB for each request (750 * 750 * 2)
	wuSpeed := 1.125 * 8 * 2 / fTime.Sub(sTime.Add(latency)).Seconds()

	// Decide workload by warm up speed
	workload := 0
	weight := 0
	skip := false
	if 10.0 < wuSpeed {
		workload = 16
		weight = 4
	} else if 4.0 < wuSpeed {
		workload = 8
		weight = 4
	} else if 2.5 < wuSpeed {
		workload = 4
		weight = 4
	} else {
		skip = true
	}

	// Main speedtest
	dlSpeed := wuSpeed
	if skip == false {
		sTime = time.Now()
		for i := 0; i < workload; i++ {
			wg.Add(1)
			go downloadRequest(wg, dlURL, weight)
		}
		wg.Wait()
		fTime = time.Now()

		reqMB := dlSizes[weight] * dlSizes[weight] * 2 / 1000 / 1000
		dlSpeed = float64(reqMB) * 8 * float64(workload) / fTime.Sub(sTime).Seconds()
	}
	return dlSpeed, nil
}

func (s *TestServer) testUpload() (float64, error) {
	if len(s.URL) == 0 {
		return 0, gerrors.New("Empty test server URL")
	}
	latency, err := s.testLatency()
	if err != nil {
		return 0, err
	}

	wg := new(sync.WaitGroup)

	// Warm up
	sTime := time.Now()
	wg = new(sync.WaitGroup)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go ulWarmUp(wg, s.URL)
	}
	wg.Wait()
	fTime := time.Now()
	// 1.0 MB for each request
	wuSpeed := 1.0 * 8 * 2 / fTime.Sub(sTime.Add(latency)).Seconds()

	// Decide workload by warm up speed
	workload := 0
	weight := 0
	skip := false
	if 10.0 < wuSpeed {
		workload = 16
		weight = 9
	} else if 4.0 < wuSpeed {
		workload = 8
		weight = 9
	} else if 2.5 < wuSpeed {
		workload = 4
		weight = 5
	} else {
		skip = true
	}

	// Main speedtest
	ulSpeed := wuSpeed
	if skip == false {
		sTime = time.Now()
		for i := 0; i < workload; i++ {
			wg.Add(1)
			go uploadRequest(wg, s.URL, weight)
		}
		wg.Wait()
		fTime = time.Now()

		reqMB := float64(ulSizes[weight]) / 1000
		ulSpeed = reqMB * 8 * float64(workload) / fTime.Sub(sTime).Seconds()
	}
	return ulSpeed, nil
}

func dlWarmUp(wg *sync.WaitGroup, dlURL string) {
	size := dlSizes[2]
	url := dlURL + "/random" + strconv.Itoa(size) + "x" + strconv.Itoa(size) + ".jpg"

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)

	wg.Done()
}

func ulWarmUp(wg *sync.WaitGroup, ulURL string) {
	size := ulSizes[4]
	v := url.Values{}
	v.Add("content", strings.Repeat("0123456789", size*100-51))

	resp, err := client.PostForm(ulURL, v)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)

	wg.Done()
}

func downloadRequest(wg *sync.WaitGroup, dlURL string, w int) {
	size := dlSizes[w]
	url := dlURL + "/random" + strconv.Itoa(size) + "x" + strconv.Itoa(size) + ".jpg"

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)

	wg.Done()
}

func uploadRequest(wg *sync.WaitGroup, ulURL string, w int) {
	size := ulSizes[9]
	v := url.Values{}
	v.Add("content", strings.Repeat("0123456789", size*100-51))

	resp, err := client.PostForm(ulURL, v)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)

	wg.Done()
}

func (s *TestServer) testLatency() (time.Duration, error) {
	if len(s.URL) == 0 {
		return 0, gerrors.New("Empty test server URL")
	}

	// I guess latency.txt is almost an empty file, it's used for test latency only
	testURL := strings.Split(s.URL, "/upload")[0] + "/latency.txt"

	l := time.Hour
	for i := 0; i < 3; i++ {
		startTime := time.Now()
		resp, err := http.Get(testURL)
		endTime := time.Now()
		defer func() {
			{
				if resp != nil {
					resp.Body.Close()
				}
			}
		}()
		if err != nil {
			return 0, err
		}
		l = time.Duration(gnum.MinInt64(int64(endTime.Sub(startTime)), int64(l))) // Get min latency in three loops
	}
	l = l / 2 // 往返折算成单向
	return l, nil
}
