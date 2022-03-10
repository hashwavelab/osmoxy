package proxy

import (
	"github.com/bitly/go-simplejson"
	"github.com/hashwavelab/osmoxy/proxy/execproxy"
)

func (_p *Proxy) GetLatestBlock() (*simplejson.Json, error) {
	c := execproxy.NewPlaintextGrpcurlCommand()
	bytes, err := c.MaxTime(QueryTimeOut).Address(_p.address).Service("cosmos.base.tendermint.v1beta1.Service/GetLatestBlock").Execute()
	if err != nil {
		return nil, err
	}
	j, err := simplejson.NewJson(bytes)
	if err != nil {
		return nil, err
	}
	return j, err
}

func (_p *Proxy) GetPools() (*simplejson.Json, error) {
	c := execproxy.NewPlaintextGrpcurlCommand()
	bytes, err := c.MaxTime(QueryTimeOut).Address(_p.address).Service("osmosis.gamm.v1beta1.Query/Pools").Data(`{"pagination":{"limit":1000}}`).Execute()
	if err != nil {
		return nil, err
	}
	j, err := simplejson.NewJson(bytes)
	if err != nil {
		return nil, err
	}
	return j, err
}
