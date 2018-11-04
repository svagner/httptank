package main

type HttpErrors struct {
	E50x, E40x, ETimeout, EOther int64
}

type tankTrace struct {
	Count, Error, Time, MinTime, MaxTime int64
	Errors                               HttpErrors
	ElapsedTime                          float64
}

type tankSettings struct {
	Url                  string
	Timeout, Count, Time int64
	Username, Password   string
	Useragent, Cookie    string
}
