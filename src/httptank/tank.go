package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"text/template"
	"time"

	"httptank/random_data"

	"github.com/golang/glog"
)

type Tank struct {
	DataChan    chan *tankTrace
	CleanChan   chan bool
	NoMatchChan chan string
	done        chan bool
	stop        chan bool
	start       chan tankSettings
	runningJobs int64
	jobResult   chan tankTrace

	Err error
}

func NewTank() *Tank {
	return &Tank{
		DataChan:    make(chan *tankTrace, 1),
		CleanChan:   make(chan bool),
		NoMatchChan: make(chan string, 1),
		done:        make(chan bool),
		stop:        make(chan bool),
		start:       make(chan tankSettings),
	}
}

func (t *Tank) Run() {
	var (
		settings tankSettings
	)
	stopTimer := time.NewTimer(0)
	stat := new(tankTrace)
	res := make(chan tankTrace)
	stopJobs := make([]chan bool, 0)
	tick := time.NewTicker(1 * time.Second)
	var StartTime = time.Now()
	for {
		select {
		case _ = <-t.stop:
			for i := 0; i < len(stopJobs); i++ {
				close(stopJobs[i])
			}
			stopJobs = make([]chan bool, 0)

		case settings = <-t.start:
			if settings.Url == "" {
				settings = tankSettings{}
				continue
			}
			if settings.Count == 0 {
				settings.Count = 1
			}
			if settings.Time != 0 {
				stopTimer = time.NewTimer(time.Duration(settings.Time) * time.Second)
			}
			t.CleanChan <- true
			StartTime = time.Now()
			for i := int64(0); i < settings.Count; i++ {
				stopChan := make(chan bool)
				stopJobs = append(stopJobs, stopChan)
				go t.HttpClient(settings, res, stopChan)
			}
		case dt := <-res:
			stat.Count += dt.Count
			stat.Error += dt.Error
			stat.Errors.E50x += dt.Errors.E50x
			stat.Errors.E40x += dt.Errors.E40x
			stat.Errors.ETimeout += dt.Errors.ETimeout
			stat.Errors.EOther += dt.Errors.EOther
			if stat.MinTime > dt.Time || stat.MinTime == 0 {
				stat.MinTime = dt.Time
			}
			if stat.MaxTime < dt.Time || stat.MaxTime == 0 {
				stat.MaxTime = dt.Time
			}

		case _ = <-tick.C:
			if stat.Count == 0 {
				continue
			}
			stat.ElapsedTime = time.Now().Sub(StartTime).Seconds()
			t.DataChan <- stat
			stat = new(tankTrace)
		case _ = <-stopTimer.C:
			for i := 0; i < len(stopJobs); i++ {
				close(stopJobs[i])
			}
			stopJobs = make([]chan bool, 0)
		}
	}
}

func (t *Tank) HttpClient(settings tankSettings, res chan tankTrace, stop chan bool) {
	for {
		select {
		case _ = <-stop:
			return
		default:
			var (
				stat    tankTrace
				cookies []*http.Cookie
			)
			jar, _ := cookiejar.New(nil)
			cookiesText := strings.Split(settings.Cookie, ",")
			if len(cookiesText) != 0 {
				for _, cook := range cookiesText {
					cookiesValues := strings.Split(cook, ":")
					if len(cookiesValues) < 4 {
						continue
					}
					c := &http.Cookie{
						Name:   cookiesValues[0],
						Value:  cookiesValues[1],
						Path:   cookiesValues[2],
						Domain: cookiesValues[3],
					}
					cookies = append(cookies, c)
				}
			}
			uniqParams := random_data.GetRandArg()
			Data := struct {
				Param string
			}{
				uniqParams,
			}
			buf := new(bytes.Buffer)
			err := template.Must(template.New("put").Parse(settings.Url)).Execute(buf, Data)
			if err != nil {
				glog.Warningln("Error while render template for put data: ", err)
				continue
			}
			cookieURL, _ := url.Parse(buf.String())
			jar.SetCookies(cookieURL, cookies)
			timeout := time.Duration(time.Duration(settings.Timeout) * time.Millisecond)
			httpTransport := &http.Transport{
				DisableKeepAlives: true,
			}
			client := http.Client{
				Transport: httpTransport,
				Timeout:   timeout,
				Jar:       jar,
			}
			stat.Count++
			start := time.Now()
			req, err := http.NewRequest("GET", buf.String(), nil)
			if err != nil {
				glog.Warningln("Generate request", buf.String(), "failed:", err)
				continue
			}
			if settings.Useragent != "" {
				req.Header.Set("User-Agent", settings.Useragent)
			}
			if settings.Username != "" || settings.Password != "" {
				req.SetBasicAuth(settings.Username, settings.Password)
			}
			resp, err := client.Do(req)
			if err != nil {
				stat.Error++
				switch err.(type) {
				case *url.Error:
					stat.Errors.ETimeout++
				case net.Error:
					stat.Errors.ETimeout++
				default:
					stat.Errors.EOther++
					glog.V(2).Infoln("Get", buf.String(), "error", err)
				}
			} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				glog.V(2).Infoln("Get", buf.String(), resp.StatusCode)
				stat.Error++
				stat.Errors.E40x++
			} else if resp.StatusCode >= 500 {
				glog.V(2).Infoln("Get", buf.String(), resp.StatusCode)
				stat.Error++
				stat.Errors.E50x++
			}
			stat.Time = int64(time.Now().Sub(start))
			if err == nil {
				io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
			}
			glog.V(1).Infoln("Tank sent data:", stat)
			res <- stat
		}
	}
}
