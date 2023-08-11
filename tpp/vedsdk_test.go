package tpp

import (
	"testing"
)

func TestQueryAll(t *testing.T) {
	var v VenafiTpp
	var config Config
	config.GetConf("/home/porteb/dev/tool/venafi-exporter/conf/venafi.yml")
	v.Connect(&config)
	v.QueryAll()

	if len(v.Certificates) < 1 {
		t.Fatalf(`len(v.QueryAll.Certificates) < 0. Must be > 0.`)
	}

}
