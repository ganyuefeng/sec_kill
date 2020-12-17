package service

import (
	"github.com/astaxie/beego/logs"
	etcd_client "github.com/coreos/etcd/clientv3"
	"time"
)

func initEtcd(conf *SecLayerConf) (err error) {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{conf.EtcdConfig.EtcdAddr},
		DialTimeout: time.Duration(conf.EtcdConfig.Timeout) * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err:", err)
		return
	}

	secLayerContext.etcdClient = cli
	logs.Debug("init etcd succ")
	return
}

func InitSecLayer(conf *SecLayerConf) (err error) {

	err = initRedis(conf)
	if err != nil {
		logs.Error("initRedis failed!", err)
	}

	err = initEtcd(conf)
	if err != nil {
		logs.Error("initEtcd failed!", err)
	}

	logs.Debug("init etcd succ")
	err = loadProductFromEtcd(conf)
	if err != nil {
		logs.Error("load product from etcd failed, err:%v", err)
		return
	}

	secLayerContext.secLayerConf = conf
	secLayerContext.Read2HandleChan = make(chan *SecRequest,
		secLayerContext.secLayerConf.Read2handleChanSize)

	secLayerContext.Handle2WriteChan = make(chan *SecResponse,
		secLayerContext.secLayerConf.Handle2WriteChanSize)

	secLayerContext.HistoryMap = make(map[int]*UserBuyHistory, 1000000)

	secLayerContext.productCountMgr = NewProductCountMgr()
	return
}
