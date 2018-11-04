package main

import (
	"flag"

	"fmt"
	"runtime"

	"github.com/golang/glog"
	"github.com/pkg/browser"
	"httptank/random_data"
)

var version = "dev"

func main() {
	version := flag.Bool("version", false, "Show version")
	host := flag.String("host", "127.0.0.1", "listen host")
	port := flag.String("port", "10000", "listen port")
	file := flag.String("file", "", "file with parameters list")
	htpasswd := flag.String("htpasswd", "", ".htpasswd path location for http-basic auth user")
	mode := flag.Bool("client", false, "mode")
	flag.Parse()

	if *version {
		printVersion()
		return
	}
	runtime.GOMAXPROCS(runtime.NumCPU())

	random_data.Init(*file)

	graph := NewGraph("Tank test", TANK_TMPL)
	server := NewHttpServer(*host, *port, *htpasswd, &graph)
	tank := NewTank()
	go tank.Run()
	go server.Start(tank.start, tank.stop)

	if *mode {
		err := browser.OpenURL(server.Url())
		if err != nil {
			glog.Warning("Open browser error: " + err.Error() + ". Please open link in your browser: http://" + *host + ":" + *port)
		}
	} else {
		glog.Info("Listen: http://" + *host + ":" + *port)
	}

	for {
		select {
		case trace := <-tank.DataChan:
			graph.AddTankGraphPoint(trace)
		case _ = <-tank.CleanChan:
			graph.CleanStat()
		}
	}
}

func printVersion() {
	fmt.Println("HTTPTank version " + version)
}
