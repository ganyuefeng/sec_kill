package service

import (
	"fmt"

	"sync"
	//"seckill_proxy/conf"
	"github.com/astaxie/beego/logs"
)

type SecLimitMgr struct {
	UserLimitMap map[int]*Limit
	IpLimitMap   map[string]*Limit
	lock         sync.Mutex
}

func antiSpam(req *SecRequest) (err error) {
	logs.Debug("this is antiSpam")
	_, ok := Seckillconf.IdBlackMap[req.UserId]
	if ok {
		err = fmt.Errorf("invalid request")
		logs.Error("useId[%v] is block by id black", req.UserId)
		return
	}
	/*
		_, ok = Seckillconf.IpBlackMap[req.ClientAddr]
		if ok {
			err = fmt.Errorf("invalid request")
			logs.Error("useId[%v] ip[%v] is block by ip black", req.UserId, req.ClientAddr)
			return
		}

	*/

	Seckillconf.SecLimitMgr.lock.Lock()
	//uid 频率控制
	limit, ok := Seckillconf.SecLimitMgr.UserLimitMap[req.UserId]
	logs.Debug("22this is antiSpam")
	if !ok {
		limit = &Limit{
			secLimit: &SecLimit{},
			minLimit: &MinLimit{},
		}
		Seckillconf.SecLimitMgr.UserLimitMap[req.UserId] = limit
	}
	logs.Debug("this is antiSpam")
	secIdCount := limit.secLimit.Count(req.AccessTime.Unix())
	minIdCount := limit.minLimit.Count(req.AccessTime.Unix())

	//ip 频率控制
	limit, ok = Seckillconf.SecLimitMgr.IpLimitMap[req.ClientAddr]
	if !ok {
		limit = &Limit{
			secLimit: &SecLimit{},
			minLimit: &MinLimit{},
		}
		Seckillconf.SecLimitMgr.IpLimitMap[req.ClientAddr] = limit
	}

	secIpCount := limit.secLimit.Count(req.AccessTime.Unix())
	minIpCount := limit.minLimit.Count(req.AccessTime.Unix())

	Seckillconf.SecLimitMgr.lock.Unlock()

	if secIpCount > Seckillconf.AccessLimitConf.IPSecAccessLimit {
		err = fmt.Errorf("One second IP invalid request")
		return
	}

	if minIpCount > Seckillconf.AccessLimitConf.IPMinAccessLimit {
		err = fmt.Errorf("One minite IP invalid request")
		return
	}

	if secIdCount > Seckillconf.AccessLimitConf.UserSecAccessLimit {
		err = fmt.Errorf("One seconed user invalid request")
		return
	}

	if minIdCount > Seckillconf.AccessLimitConf.UserMinAccessLimit {
		err = fmt.Errorf("One minite user invalid request")
		return
	}

	logs.Debug("this is antiSpam out")
	return
}
