package collector

import (
	"sync"
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
        "io/ioutil"
        "os/exec"

        "fmt"
	"strings"
)

type Metrics struct {
	metrics map[string]*prometheus.Desc
	mutex   sync.Mutex
}

func newGlobalMetric(namespace string, metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName, docString, labels, nil)
}


func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		metrics: map[string]*prometheus.Desc{
			"gpclient_connections_metric": newGlobalMetric(namespace, "gpclient_connections_metric","gpclient connections", []string{"host"}),
		},
	}
}

func (c *Metrics) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

func (c *Metrics) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()  // 加锁
	defer c.mutex.Unlock()


        for k, v := range c.GetConnectionsData() {
                fmt.Println("hosts:", k, "value:", v)
		ch <-prometheus.MustNewConstMetric(c.metrics["gpclient_connections_metric"], prometheus.GaugeValue, v, k)
        }


}

func (c *Metrics) GetConnectionsData() map[string]float64 {
	result := ExecCommand("netstat -nat | grep ESTABLISHED")
        mapConns := make(map[string]float64)

        for _, v := range strings.Split(result, "\n") {
                if strings.Contains(v, "ESTABLISHED") {
                        lines := strings.Fields(v)
                        if "10.247.32.84:5432" == lines[3] {
                                delIndex := strings.Index(lines[4], ":")
                                ip := lines[4][:delIndex]
                                mapConns[ip] = mapConns[ip] + 1
                        }
                }
        }

	return mapConns
}

 func (c *Metrics) GenerateMockData() (mockCounterMetricData map[string]int, mockGaugeMetricData map[string]int) {
 	mockCounterMetricData = map[string]int{
		"yahoo.com": int(rand.Int31n(1000)),
		"google.com": int(rand.Int31n(1000)),
	}
	mockGaugeMetricData = map[string]int{
		"yahoo.com": int(rand.Int31n(10)),
		"google.com": int(rand.Int31n(10)),
	}
	return
 }


func ExecCommand(strCommand string) string {

        cmd := exec.Command("/bin/bash", "-c", strCommand)
        stdout, _ := cmd.StdoutPipe()
        if err := cmd.Start(); err != nil {
                fmt.Println("Execute failed when Start:" + err.Error())
                return ""
        }

        out_bytes, _ := ioutil.ReadAll(stdout)
        stdout.Close()

        if err := cmd.Wait(); err != nil {
                fmt.Println("Execute failed when Wait:" + err.Error())
                return ""
        }
        return string(out_bytes)
}

