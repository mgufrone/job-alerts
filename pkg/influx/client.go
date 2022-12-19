package influx

import (
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
	"os"
)

func NewClient() (api.WriteAPI, api.QueryAPI) {
	host := fmt.Sprintf("http://%s:%s", os.Getenv("INFLUXDB_HOST"), os.Getenv("INFXLUDB_PORT"))
	cli := influxdb2.NewClient(host, os.Getenv("INFLUXDB_TOKEN"))
	org := os.Getenv("INFLUXDB_ORG")
	bucket := os.Getenv("INFLUXDB_BUCKET")
	return cli.WriteAPI(org, bucket), cli.QueryAPI(org)
}
