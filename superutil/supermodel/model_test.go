package supermodel

import (
	"os"
	"testing"

	"github.com/ironzhang/superlib/fileutil"
)

func TestWriteServiceModel(t *testing.T) {
	m := ServiceModel{
		Domain:             "http.myapp",
		DefaultDestination: "dev.default.k8s",
		Clusters: []Cluster{
			{
				Name: "dev.default.k8s",
				Labels: map[string]string{
					ZoneKey: "dev",
					LaneKey: "default",
					KindKey: "k8s",
				},
				Endpoints: []Endpoint{
					{
						Addr:   "192.168.1.1:8000",
						State:  Enabled,
						Weight: 100,
					},
					{
						Addr:   "192.168.1.2:8000",
						State:  Enabled,
						Weight: 100,
					},
				},
			},
			{
				Name: "dev.sim00.k8s",
				Labels: map[string]string{
					ZoneKey: "dev",
					LaneKey: "sim00",
					KindKey: "k8s",
				},
				Endpoints: []Endpoint{
					{
						Addr:   "192.168.2.1:8000",
						State:  Enabled,
						Weight: 100,
					},
					{
						Addr:   "192.168.2.2:8000",
						State:  Enabled,
						Weight: 100,
					},
				},
			},
		},
	}
	os.MkdirAll("./testdata/services", os.ModePerm)
	fileutil.WriteJSON("./testdata/services/http.myapp.json", m)
}

func TestWriteRouteModel(t *testing.T) {
	m := RouteModel{
		Domain: "http.myapp",
		Policy: RoutePolicy{
			EnableScript: false,
		},
	}
	os.MkdirAll("./testdata/routes", os.ModePerm)
	fileutil.WriteJSON("./testdata/routes/http.myapp.json", m)
}
