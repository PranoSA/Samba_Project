package sambaservermanagement

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
			fileSpace.Set(500)
			fileSpace2.Add(55)
		}
	}()
}

func ExposeFSMetrics() {
	go func() {
		for {
			fileSpace.Add(1)
			fileSpace2.Set(2)
			time.Sleep(2 * time.Second)
		}
	}()
}

var DiskLabel []map[string]string = make([]map[string]string, 0)

var DiskGauges []prometheus.Gauge

func InitFromDiskLabels(fs []FileSystem) {
	for i, fs := range fs {
		mount_path := fs.MouthPath
		if fs.MouthPath == "" {
			mount_path = "/mount/samba_server/" + fs.Fsid
		}
		DiskLabel = append(DiskLabel, make(map[string]string))
		DiskLabel[i]["Device"] = fs.Dev
		DiskLabel[i]["Fsid"] = fs.Fsid
		DiskLabel[i]["space"] = strconv.FormatInt(fs.RoomLeft, 10)
		DiskLabel[i]["mount_path"] = mount_path

		DiskGauges = append(DiskGauges, promauto.NewGauge(prometheus.GaugeOpts{
			Namespace:   "Disks",
			Name:        "Samba_File_Systems",
			Help:        "Room,Mount Path, and fsid of allocated file systems",
			ConstLabels: DiskLabel[i],
		}))

		DiskGauges[i].Set(float64(fs.RoomLeft))
	}

}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})

	labels2   = map[string]string{"disk": "/dev/sda2"}
	fileSpace = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace:   "poopoo",
		Name:        "Gauge_Opts",
		Help:        "Seee",
		ConstLabels: labels2,
	})

	labels = map[string]string{"disk": "/dev/sda1"}

	fileSpace2 = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace:   "poopoo",
		Name:        "Gauge_Opts",
		Help:        "Seee",
		ConstLabels: labels,
	})
)

func Initialize() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
		wg.Done()
	}()

	wg.Wait()
}
