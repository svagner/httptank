package httptank

import (
	"html/template"
	"io"
	"sync"
	"time"

	"github.com/golang/glog"
)

type graphPoints [2]float64

type Graph struct {
	Title                                        string
	Settings                                     tankSettings
	Count, Error, MaxTime, MinTime, AvgTime      []graphPoints
	Error50x, Error40x, ErrorOther, ErrorTimeout []graphPoints
	Tmpl                                         *template.Template `json:"-"`
	mu                                           sync.RWMutex       `json:"-"`
}

var StartTime = time.Now()

func NewGraph(title, tmpl string) Graph {
	g := Graph{
		Title:        title,
		Count:        []graphPoints{},
		Error:        []graphPoints{},
		Error40x:     []graphPoints{},
		Error50x:     []graphPoints{},
		ErrorTimeout: []graphPoints{},
		ErrorOther:   []graphPoints{},
		MaxTime:      []graphPoints{},
		MinTime:      []graphPoints{},
		AvgTime:      []graphPoints{},
	}
	g.setTmpl(tmpl)

	return g
}

func (g *Graph) setSettings(settings tankSettings) {
	g.Settings = settings
}

func (g *Graph) setTmpl(tmplStr string) {
	g.Tmpl = template.Must(template.New("vis").Parse(tmplStr))
}

func (g *Graph) Write(w io.Writer) error {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.Tmpl.Execute(w, g)
}

func (g *Graph) CleanStat() {
	g.Count = []graphPoints{}
	g.Error = []graphPoints{}
	g.MaxTime = []graphPoints{}
	g.MinTime = []graphPoints{}
	g.AvgTime = []graphPoints{}
	g.Error40x = []graphPoints{}
	g.Error50x = []graphPoints{}
	g.ErrorTimeout = []graphPoints{}
	g.ErrorOther = []graphPoints{}
}

func (g *Graph) AddTankGraphPoint(data *tankTrace) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var elapsedTime float64
	glog.V(1).Infoln("Data for draw:", *data)
	if data.ElapsedTime == 0 {
		elapsedTime = time.Now().Sub(StartTime).Seconds()
	} else {
		elapsedTime = data.ElapsedTime
	}
	g.Error = append(g.Error, graphPoints{elapsedTime, float64(data.Error)})
	g.Error50x = append(g.Error50x, graphPoints{elapsedTime, float64(data.Errors.E50x)})
	g.Error40x = append(g.Error40x, graphPoints{elapsedTime, float64(data.Errors.E40x)})
	g.ErrorTimeout = append(g.ErrorTimeout, graphPoints{elapsedTime, float64(data.Errors.ETimeout)})
	g.ErrorOther = append(g.ErrorOther, graphPoints{elapsedTime, float64(data.Errors.EOther)})
	g.Count = append(g.Count, graphPoints{elapsedTime, float64(data.Count)})
	var (
		avgtime int64
	)

	avgtime = (data.MinTime + data.MaxTime) / 2

	g.AvgTime = append(g.AvgTime, graphPoints{elapsedTime, float64(avgtime / 1000000)})
	g.MinTime = append(g.MinTime, graphPoints{elapsedTime, float64(data.MinTime / 1000000)})
	g.MaxTime = append(g.MaxTime, graphPoints{elapsedTime, float64(data.MaxTime / 1000000)})

}
