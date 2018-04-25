package main

import (

	"bytes"
	"net/http"
	"fmt"
	"time"

	"github.com/gin-gonic/gin/json"
	"math/rand"

	"github.com/tcnksm/go-httpstat"
	"runtime"
	"sync"
)

const (
	MAXConCurrency = 10
	HeaderContentType = "application/json"
	TimeFormat = "2006-01-02 15:04:05"
)

var(

	seq=0
	over = make(chan bool)
  	sem = make(chan bool, MAXConCurrency) //控制并发任务数
	maxCount = 100
)
func checkErr(err error) {
	if err != nil {
		fmt.Errorf("程序执行异常: %v", err)
		//panic(err)
	}
}

//定义KPI结构体
type Kpi struct {
	NodeKey string   `json:"nodeKey"`
	Severity int `json:"severity"`
	Message	string `json:"message,omitempty"`
	ArisingTime string `json:"arisingTime"`
	AgentId string `json:agentId`
	Title string `json:title,omitempty`
	AlertKey string `json:alertKey,omitempty`
}

//定义性能指标结构体
type Performance struct {
	Status 		string `json:Status`
	RespCode	string `json:RespCode`
	KeepAlive	string `json:keepAlive,omitempty`

	Duration 	string `json:duration`
	Dns 		string `json:dns,omitempty`
	Connect		string `json:connect,omitempty`
	Request 	string `json:request,omitempty`
	Response	string `json:response,omitempty`
	Latency 	string `json:latency,omitempty`
	Speed		string `json:speed,omitempty`
	RespSpeed	string `json:ResponseSpeed,omitempty`

}

func generateTimeString() string {

	var secs = time.Now().Unix() - int64(seq)
	t := time.Unix(secs,0)
	timeStr := t.Format(TimeFormat)
	//fmt.Printf("DEBUG:%v \r\n",timeStr)
	return timeStr
}

func parsePerformance(header http.Header) Performance {
	var perf Performance
	perf.Status 	= header.Get("Status")
	perf.RespCode 	= header.Get("Response Code")
	perf.KeepAlive	= header.Get("Kept Alive")
	perf.Duration	= header.Get("Duration")
	perf.Dns		= header.Get("DNS")
	perf.Connect	= header.Get("Connect")
	perf.Request	= header.Get("Request")
	perf.Response	= header.Get("Response")
	perf.Latency	= header.Get("Latency")

	//s,err := json.Marshal(perf)
	//checkErr(err)
	//ss := string(s)

	return perf
}

func createKpis() []Kpi{
	var kpis []Kpi


	seq++
	var kpi Kpi
	//构造一个告警①
	kpi.NodeKey="bjhr"
	kpi.Severity=rand.Intn(4)+1
	kpi.Message=fmt.Sprintf("golang告警推送,seq=%v", seq)
	kpi.ArisingTime= generateTimeString()
	kpi.AgentId = "Inspection"
	kpi.Title = fmt.Sprintf("golang告警标题,seq=%v", seq)
	kpi.AlertKey = "ApplicationSystem_SubHealth_Alert"
	kpis = append(kpis, kpi)

	//构造一个告警②
	kpi.NodeKey="xmhr"
	kpi.Severity=rand.Intn(4)+1
	kpi.Message=fmt.Sprintf("golang告警推送,seq=%v", seq)
	kpi.ArisingTime= generateTimeString()
	kpi.AgentId = "Inspection"
	kpi.Title = fmt.Sprintf("golang告警标题,seq=%v", seq)
	kpi.AlertKey = "ApplicationSystem_SubHealth_Alert"

	kpis = append(kpis, kpi)

	return kpis
}

func postKpi() string{
	//if time.Now().Unix()

	sem <- true
	defer func() {
		<-sem
	}()


	kpis := createKpis()
	s, _ := json.Marshal(kpis)
	body := bytes.NewBuffer(s)

	fmt.Printf("Request Seq=%v\t Response:\t\r\n",seq)
	fmt.Printf("DEBUG:\t %v \r\n",string(s))
	req, err := http.NewRequest("POST", "http://192.168.20.171:48080/restcenter/innerKpiApi/sendNewKpi/InspectionKpiData", body)

	req.Header.Set("Content-Type", HeaderContentType)



	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)


	resp, err := http.DefaultClient.Do(req)
	checkErr(err)

	result.End(time.Now())
	fmt.Printf("%v \r\n %+v \r\n",resp.Status, result)

	result.End(time.Now())

	perf := parsePerformance(resp.Header)
	perf_js,err := json.MarshalIndent(&perf,"","\t")
	perf_json := string(perf_js)

	defer resp.Body.Close()


	//fmt.Printf("Response Seq=%v\t Response:\t\r\n",seq)
	//fmt.Print("\t%v",perf_json)
	return perf_json
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup

	for i:=0; i<maxCount;i++  {
		wg.Add(1)
		fmt.Printf("#######For num:%v\n", i)
		go postKpi()
		time.Sleep(1* time.Second)
	}
	wg.Wait() //等待所有goroutine退出
	//fmt.Println(s)
}