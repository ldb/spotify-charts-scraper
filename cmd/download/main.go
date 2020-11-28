package main

import (
	"errors"
	"fmt"
	"github.com/ldb/spotify/pkg/spotify"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

const (
	baseURL        = "https://spotifycharts.com"
	dataPath       = "./data"
	maxConcurrency = 50
)

func main() {
	minDate := time.Date(2017, 01, 01, 0, 0, 0, 0, time.UTC)
	maxDate := time.Date(2020, 10, 30, 0, 0, 0, 0, time.UTC)

	var dates []string
	currentDate := minDate
	for {
		if currentDate.Equal(maxDate) {
			break
		}
		dates = append(dates, dateString(currentDate))
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	log.Printf("%d days, %d regions\n", len(dates), len(spotify.Regions))

	maxRegions := make(chan struct{}, 10)
	wg := sync.WaitGroup{}
	progress := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithRefreshRate(180*time.Millisecond))
	wg.Add(len(dates) * len(spotify.Regions))
	for regCode, regName := range spotify.Regions {
		regCode := regCode
		regName := regName
		maxRegions <- struct{}{}
		go func() {
			defer func() {
				<-maxRegions
			}()
			time.Sleep(time.Second * time.Duration(rand.Intn(10)))
			os.MkdirAll(path.Join(dataPath, regCode), 0700)
			nDates := int64(len(dates))
			bar := progress.AddBar(nDates,
				mpb.PrependDecorators(
					decor.Name(regName, decor.WC{W: 20, C: decor.DidentRight}),
					decor.OnComplete(decor.CountersNoUnit("%d / %d", decor.WCSyncWidth), "done!"),
				),
				mpb.BarFillerClearOnComplete(),
				mpb.AppendDecorators(
					decor.OnComplete(decor.AverageETA(decor.ET_STYLE_GO, decor.WCSyncSpace), ""),
					decor.Percentage(decor.WC{W: 6}),
				),
			)

			maxConcurrency := make(chan struct{}, maxConcurrency)
			for _, date := range dates {
				maxConcurrency <- struct{}{}
				go func(date string) {
					defer func() {
						wg.Done()
						<-maxConcurrency
					}()
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))
					err := retry(5, 3*time.Second, func() error {
						r, err := http.Get(fmt.Sprintf("%s/regional/%s/daily/%s/download", baseURL, regCode, date))
						if err != nil {
							return err
						}
						defer r.Body.Close()
						if r.StatusCode == http.StatusNotFound {
							// If the file is not there, there is no point in retrying.
							return stop{errors.New("not found")}
						}
						if h := r.Header.Get("Content-Type"); h != "text/csv;charset=UTF-8" {
							return errors.New("non CSV data")
						}
						if r.StatusCode != http.StatusOK {
							return errors.New("non 200 code")
						}
						p := path.Join(dataPath, fmt.Sprintf("%s/%s", regCode, date))
						f, err := os.Create(p)
						if err != nil {
							return err
						}
						defer f.Close()
						_, err = io.Copy(f, r.Body)
						if err != nil {
							return err
						}
						return nil
					})
					if err != nil {
					}
					bar.Increment()
				}(date)
			}
		}()
	}
	wg.Wait()
}

func dateString(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

// https://upgear.io/blog/simple-golang-retry-function/
func retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if s, ok := err.(stop); ok {
			return s.error
		}

		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return retry(attempts, 2*sleep, fn)
		}
		return err
	}
	return nil
}

type stop struct {
	error
}
