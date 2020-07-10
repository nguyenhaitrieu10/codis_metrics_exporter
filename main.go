package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var DOMAIN = os.Getenv("CODIS_HOST")
var CODIS_API = DOMAIN + "/proxy/stats"
var INTERVAL = 5 * time.Second

var (
	ops_qps          = promauto.NewGauge(prometheus.GaugeOpts{Name: "codis_ops_qps"})
	sessions_alive   = promauto.NewGauge(prometheus.GaugeOpts{Name: "codis_sessions_alive"})
	rusage_cpu       = promauto.NewGauge(prometheus.GaugeOpts{Name: "codis_rusage_cpu"})
	rusage_mem       = promauto.NewGauge(prometheus.GaugeOpts{Name: "codis_rusage_mem"})
	ops_fails        = promauto.NewGauge(prometheus.GaugeOpts{Name: "codis_ops_fails"})
	ops_redis_errors = promauto.NewGauge(prometheus.GaugeOpts{Name: "codis_ops_redis_errors"})
	runtime_gc_num   = promauto.NewGauge(prometheus.GaugeOpts{Name: "codis_runtime_gc_num"})
)

type MetricsCodis struct {
	Ops struct {
		Failt uint `json:"fails,omitempty"`
		Redis struct {
			Errors uint `json:"errors,omitempty"`
		} `json:"redis,omitempty"`
		Qps uint `json:"qps,omitempty"`
	} `json:"ops,omitempty"`

	Sessions struct {
		Total uint `json:"total,omitempty"`
		Alive uint `json:"alive,omitempty"`
	} `json:"sessions,omitempty"`

	Rusage struct {
		Cpu float64 `json:"cpu,omitempty"`
		Mem uint64  `json:"mem,omitempty"`
	} `json:"rusage,omitempty"`

	Runtime struct {
		Gc struct {
			Num uint `json:"num,omitempty"`
		} `json:"gc,omitempty"`
	} `json:"runtime,omitempty"`
}

func recordMetrics() {
	go func() {
		for {
			resp, err := http.Get(CODIS_API)
			if err != nil {
				fmt.Println("Request Error")
			} else if resp.StatusCode == 200 {
				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Read Body Response Error")
					time.Sleep(INTERVAL)
					continue
				}

				var metrics MetricsCodis
				err = json.Unmarshal(data, &metrics)
				if err != nil {
					fmt.Println("Unmarshal Body Response Error")
					time.Sleep(INTERVAL)
					continue
				}

				ops_qps.Set(float64(metrics.Ops.Qps))
				sessions_alive.Set(float64(metrics.Sessions.Alive))
				rusage_cpu.Set(float64(metrics.Rusage.Cpu))
				rusage_mem.Set(float64(metrics.Rusage.Mem))
				ops_fails.Set(float64(metrics.Ops.Failt))
				ops_redis_errors.Set(float64(metrics.Ops.Redis.Errors))
				runtime_gc_num.Set(float64(metrics.Runtime.Gc.Num))
			} else {
				fmt.Println("Status Code ", resp.StatusCode)
			}

			time.Sleep(INTERVAL)
		}
	}()
}

func main() {
	if DOMAIN == "" {
		fmt.Println("CODIS_HOST is not set, using default: http://localhost:11080")
		CODIS_API = "http://localhost:11080/proxy/stats"
	}

	recordMetrics()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
