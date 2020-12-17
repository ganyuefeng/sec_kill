package service

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

func InitService() (err error) {
	//secKillConf = serviceConf

	err = loadBlackList()
	if err != nil {
		logs.Error("load black list err:%v", err)
		return
	}

	err = initProxy2LayerRedis()
	if err != nil {
		logs.Error("load proxy2layer redis pool failed, err:%v", err)
		return
	}

	logs.Debug("init service succ, config:%v", Seckillconf)

	Seckillconf.SecLimitMgr = &SecLimitMgr{
		UserLimitMap: make(map[int]*Limit, 10000),
		IpLimitMap:   make(map[string]*Limit, 10000),
	}

	Seckillconf.SecReqChan = make(chan *SecRequest, Seckillconf.SecReqChanSize)
	Seckillconf.UserConnMap = make(map[string]chan *SecResult, 10000)

	initRedisProcessFunc()

	return
}

func initRedisProcessFunc() {
	logs.Info("111111:", Seckillconf.WriteProxy2LayerGoroutineNum)
	for i := 0; i < Seckillconf.WriteProxy2LayerGoroutineNum; i++ {
		logs.Info("2222")
		go WriteHandle()
	}

	for i := 0; i < Seckillconf.ReadProxy2LayerGoroutineNum; i++ {
		go ReadHandle()
	}
}

func initProxy2LayerRedis() (err error) {
	Seckillconf.proxy2LayerRedisPool = &redis.Pool{
		MaxIdle:     Seckillconf.RedisProxy2LayerConf.Redis_maxIdle,
		MaxActive:   Seckillconf.RedisProxy2LayerConf.Redis_maxActive,
		IdleTimeout: time.Duration(Seckillconf.RedisProxy2LayerConf.Redis_idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", Seckillconf.RedisProxy2LayerConf.Redis_addr)
		},
	}

	conn := Seckillconf.proxy2LayerRedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err:%v", err)
		return
	}

	return
}

func initBlackRedis() (err error) {
	logs.Debug("initredis start")
	Seckillconf.blackRedisPool = &redis.Pool{
		MaxIdle:     Seckillconf.RedisBlackConf.Redis_maxIdle,
		MaxActive:   Seckillconf.RedisBlackConf.Redis_maxActive,
		IdleTimeout: time.Duration(Seckillconf.RedisBlackConf.Redis_idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", Seckillconf.RedisBlackConf.Redis_addr)
		},
	}

	conn := Seckillconf.blackRedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed ", err)
		return
	}

	logs.Debug("initRedis success!!")
	return
}

func loadBlackList() (err error) {
	err = initBlackRedis()
	if err != nil {
		logs.Error("init reids failed")
	}
	conn := Seckillconf.blackRedisPool.Get()
	defer conn.Close()

	reply, err := conn.Do("hgetall", "idblacklist")
	if err != nil {
		logs.Error("hget all failed")
	}

	idlist, err := redis.Strings(reply, err)
	if err != nil {
		return
	}

	for _, v := range idlist {
		id, err := strconv.Atoi(v)
		if err != nil {
			logs.Warn("invalid user id %v", id)
			return err
		}
		Seckillconf.IdBlackMap[id] = true
	}

	reply, err = conn.Do("hgetall", "ipblacklist")
	if err != nil {
		logs.Error("hget all failed")
	}
	for _, v := range idlist {
		Seckillconf.IpBlackMap[v] = true
	}
	return
}
