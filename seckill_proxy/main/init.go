package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/garyburd/redigo/redis"
	etcd_client "go.etcd.io/etcd/clientv3"
	"seckill_proxy/service"
	"strings"
	"time"
)

var (
	redisPool  *redis.Pool
	etcdClient *etcd_client.Client
)

func InitConfig() (err error) {
	redis_black_addr := beego.AppConfig.String("redis_black_addr")
	etcd_addr := beego.AppConfig.String("etcd_addr")
	redis_black_MaxIdle, err := beego.AppConfig.Int("redis_black_maxIdle")
	if err != nil {
		logs.Debug("read redisMaxIdle failed")
		return err
	}
	redis_black_IdleTimeout, err := beego.AppConfig.Int("redis_black_idleTimeout")
	if err != nil {
		logs.Debug("read redis_idleTimeout failed")
		return err
	}
	redis_black_MaxActive, err := beego.AppConfig.Int("redis_black_maxActive")
	if err != nil {
		logs.Debug("read redis_maxActive failed")
		return err
	}

	logs.Debug("redis addr is %v", redis_black_addr)
	logs.Debug("etcd_addr addr is %v", etcd_addr)
	logs.Debug("redis_IdleTimeout  is %v", redis_black_IdleTimeout)
	logs.Debug("redis_MaxActive  is %v", redis_black_MaxActive)
	logs.Debug("redisMaxIdle  is %v", redis_black_MaxIdle)

	service.Seckillconf.EtcdConf.EtcdAddr = etcd_addr
	service.Seckillconf.RedisBlackConf.Redis_addr = redis_black_addr
	service.Seckillconf.RedisBlackConf.Redis_idleTimeout = redis_black_IdleTimeout
	service.Seckillconf.RedisBlackConf.Redis_maxActive = redis_black_MaxActive
	service.Seckillconf.RedisBlackConf.Redis_maxIdle = redis_black_MaxIdle

	if len(redis_black_addr) == 0 || len(etcd_addr) == 0 {
		err = fmt.Errorf("init failed")
		return
	}

	etcdTimeout, err := beego.AppConfig.Int("etcd_timeout")
	if err != nil {
		logs.Debug("etcdTimeout get failed")
		return
	}

	service.Seckillconf.EtcdConf.Timeout = etcdTimeout

	service.Seckillconf.LogLevel = beego.AppConfig.String("log_level")
	if err != nil {
		logs.Debug("log_path get failed")
		return
	}
	service.Seckillconf.LogPath = beego.AppConfig.String("log_path")
	if err != nil {
		logs.Error("log_level get failed")
		return
	}

	service.Seckillconf.EtcdConf.EtcdSecKeyPrefix = beego.AppConfig.String("etcd_sec_key_prefix")
	if err != nil {
		logs.Error("etcd_sec_key_prefix get failed")
		return
	}
	etcdSecProductKey := beego.AppConfig.String("etcd_product_key")
	if err != nil {
		logs.Error("etcd_product_key get failed")
		return
	}
	service.Seckillconf.EtcdConf.EtcdSecProductKey = fmt.Sprintf("%s/%s", service.Seckillconf.EtcdConf.EtcdSecKeyPrefix, etcdSecProductKey)
	//limit conf get
	service.Seckillconf.CookieSecretKey = beego.AppConfig.String("cookie_secretkey")
	secLimit, err := beego.AppConfig.Int("user_sec_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read user_sec_access_limit error:%v", err)
		return
	}

	service.Seckillconf.AccessLimitConf.UserSecAccessLimit = secLimit
	referList := beego.AppConfig.String("refer_whitelist")
	if len(referList) > 0 {
		service.Seckillconf.ReferWhiteList = strings.Split(referList, ",")
	}

	ipLimit, err := beego.AppConfig.Int("ip_sec_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read ip_sec_access_limit error:%v", err)
		return
	}

	service.Seckillconf.AccessLimitConf.IPSecAccessLimit = ipLimit

	minIdLimit, err := beego.AppConfig.Int("user_min_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read user_min_access_limit error:%v", err)
		return
	}

	service.Seckillconf.AccessLimitConf.UserMinAccessLimit = minIdLimit
	minIpLimit, err := beego.AppConfig.Int("ip_min_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read ip_min_access_limit error:%v", err)
		return
	}

	service.Seckillconf.AccessLimitConf.IPMinAccessLimit = minIpLimit

	/*redisProxy2LayerAddr := beego.AppConfig.String("redis_proxy2layer_addr")
	logs.Debug("read config succ, redis addr:%v", redisProxy2LayerAddr)

	Seckillconf.RedisProxy2LayerConf.redis_addr = redisProxy2LayerAddr

	if len(redisProxy2LayerAddr) == 0 {
		err = fmt.Errorf("init config failed, redis[%s]  config is null", redisProxy2LayerAddr)
		return
	}


	*/
	redisProxy2LayerAddr := beego.AppConfig.String("redis_proxy2layer_addr")
	logs.Debug("read config succ, redis addr:%v", redisProxy2LayerAddr)

	service.Seckillconf.RedisProxy2LayerConf.Redis_addr = redisProxy2LayerAddr

	if len(redisProxy2LayerAddr) == 0 {
		err = fmt.Errorf("init config failed, redis[%s]  config is null", redisProxy2LayerAddr)
		return
	}

	redisMaxIdle, err := beego.AppConfig.Int("redis_proxy2layer_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_proxy2layer_idle error:%v", err)
		return
	}

	redisMaxActive, err := beego.AppConfig.Int("redis_proxy2layer_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_proxy2layer_active error:%v", err)
		return
	}

	redisIdleTimeout, err := beego.AppConfig.Int("redis_proxy2layer_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_proxy2layer_idle_timeout error:%v", err)
		return
	}

	service.Seckillconf.RedisProxy2LayerConf.Redis_maxIdle = redisMaxIdle
	service.Seckillconf.RedisProxy2LayerConf.Redis_maxActive = redisMaxActive
	service.Seckillconf.RedisProxy2LayerConf.Redis_idleTimeout = redisIdleTimeout

	writeGoNums, err := beego.AppConfig.Int("write_proxy2layer_goroutine_num")
	if err != nil {
		err = fmt.Errorf("init config failed, read write_proxy2layer_goroutine_num error:%v", err)
		return
	}

	service.Seckillconf.WriteProxy2LayerGoroutineNum = writeGoNums

	readGoNums, err := beego.AppConfig.Int("read_layer2proxy_goroutine_num")
	if err != nil {
		err = fmt.Errorf("init config failed, read read_layer2proxy_goroutine_num error:%v", err)
		return
	}

	service.Seckillconf.ReadProxy2LayerGoroutineNum = readGoNums

	//璇诲彇涓氬姟閫昏緫灞傚埌proxy鐨剅edis閰嶇疆
	redisLayer2ProxyAddr := beego.AppConfig.String("redis_layer2proxy_addr")
	logs.Debug("read config succ, redis addr:%v", redisLayer2ProxyAddr)

	service.Seckillconf.RedisProxy2LayerConf.Redis_addr = redisLayer2ProxyAddr

	if len(redisLayer2ProxyAddr) == 0 {
		err = fmt.Errorf("init config failed, redis[%s]  config is null", redisProxy2LayerAddr)
		return
	}

	redisMaxIdle, err = beego.AppConfig.Int("redis_layer2proxy_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_layer2proxy_idle error:%v", err)
		return
	}

	redisMaxActive, err = beego.AppConfig.Int("redis_layer2proxy_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_layer2proxy_active error:%v", err)
		return
	}

	redisIdleTimeout, err = beego.AppConfig.Int("redis_layer2proxy_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_layer2proxy_idle_timeout error:%v", err)
		return
	}

	service.Seckillconf.RedisLayer2ProxyConf.Redis_maxIdle = redisMaxIdle
	service.Seckillconf.RedisLayer2ProxyConf.Redis_maxActive = redisMaxActive
	service.Seckillconf.RedisLayer2ProxyConf.Redis_idleTimeout = redisIdleTimeout
	return

	logs.Debug("InitConfig ok!")
	return
}

func initRedis() (err error) {
	logs.Debug("initredis start")
	redisPool = &redis.Pool{
		MaxIdle:     service.Seckillconf.RedisBlackConf.Redis_maxIdle,
		MaxActive:   service.Seckillconf.RedisBlackConf.Redis_maxActive,
		IdleTimeout: time.Duration(service.Seckillconf.RedisBlackConf.Redis_idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", service.Seckillconf.RedisBlackConf.Redis_addr)
		},
	}

	conn := redisPool.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed ", err)
		return
	}

	logs.Debug("initRedis success!!")
	return
}

func initEtcd() (err error) {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{service.Seckillconf.EtcdConf.EtcdAddr},
		DialTimeout: time.Duration(service.Seckillconf.EtcdConf.Timeout) * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err:", err)
		return
	}

	etcdClient = cli
	logs.Debug("initEtcd success")
	return
}

func convertLogLevel(level string) int {

	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}

	return logs.LevelDebug
}

func initLogs() (err error) {
	config := make(map[string]interface{})
	config["filename"] = service.Seckillconf.LogPath
	config["level"] = convertLogLevel(service.Seckillconf.LogLevel)

	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("marshal failed, err:", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))
	return

}

func initSecProductWatcher() {
	go watchSecProductKey(service.Seckillconf.EtcdConf.EtcdSecProductKey)
}

func watchSecProductKey(key string) {

	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err:", err)
		return
	}

	logs.Debug("begin watch key:%s", key)
	for {
		rch := cli.Watch(context.Background(), key)
		var secProductInfo []service.SecInfoConf
		var getConfSucc = true

		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s] 's config deleted", key)
					continue
				}

				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err = json.Unmarshal(ev.Kv.Value, &secProductInfo)
					if err != nil {
						logs.Error("key [%s], Unmarshal[%s], err:%v ", err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("get config from etcd, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}

			if getConfSucc {
				logs.Debug("get config from etcd succ, %v", secProductInfo)
				updateSecProductInfo(secProductInfo)
			}
		}

	}
}

func loadSecConf() (err error) {
	logs.Debug("key is %v", service.Seckillconf.EtcdConf.EtcdSecProductKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := etcdClient.Get(ctx, service.Seckillconf.EtcdConf.EtcdSecProductKey)
	defer cancel()
	if err != nil {
		logs.Error("111get [%s] from etcd failed, err:%v", service.Seckillconf.EtcdConf.EtcdSecProductKey, err)
		//return
	}
	logs.Debug("key is %v", service.Seckillconf.EtcdConf.EtcdSecProductKey)

	//logs.Debug("111key is %v", Seckillconf.etcdConf.EtcdSecProductKey)
	var secProductInfo []service.SecInfoConf
	for k, v := range resp.Kvs {
		logs.Debug("key[%v] valud[%v]", k, v)
		err = json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			logs.Error("Unmarshal sec product info failed, err:%v", err)
			return
		}

		logs.Debug("ganyuefeng info conf is [%v]", secProductInfo)
	}

	logs.Debug("loadSecConf success!")
	updateSecProductInfo(secProductInfo)
	return
}

func InitSec() (err error) {
	err = initLogs()
	if err != nil {
		//panic(err)
		logs.Error("initLogs failed err is %v", err)
		return
	}

	err = initRedis()
	if err != nil {
		//panic(err)
		logs.Error("initRedis failed err is %v", err)
		return
	}

	err = initEtcd()
	if err != nil {
		//panic(err)
		logs.Error("initEtcd failed err is %v", err)
		return
	}

	loadSecConf()
	if err != nil {
		logs.Error("loadSecConf failed")
		return
	}

	service.InitService()
	initSecProductWatcher()
	logs.Info("init sec success!")
	return
}

func updateSecProductInfo(secProductInfo []service.SecInfoConf) {

	var tmp map[int]*service.SecInfoConf = make(map[int]*service.SecInfoConf, 1024)
	for _, v := range secProductInfo {
		produtInfo := v
		tmp[v.ProductId] = &produtInfo
	}

	service.Seckillconf.RWSecProductLock.Lock()
	service.Seckillconf.SecProductInfoMap = tmp
	service.Seckillconf.RWSecProductLock.Unlock()
}
