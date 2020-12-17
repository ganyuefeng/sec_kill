package service

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"
)

func NewSecRequest() (secRequest *SecRequest) {
	secRequest = &SecRequest{
		ResultChan: make(chan *SecResult, 1),
	}

	return
}

func SecInfoList() (data []map[string]interface{}, code int, err error) {

	Seckillconf.RWSecProductLock.RLock()
	defer Seckillconf.RWSecProductLock.RUnlock()

	for _, v := range Seckillconf.SecProductInfoMap {

		item, _, err := SecInfoById(v.ProductId)
		if err != nil {
			code = ErrNotFoundProductId
			logs.Error("get product_id[%d] failed, err:%v", v.ProductId, err)
			continue
		}

		logs.Debug("get product[%d]， result[%v], all[%v] v[%v]", v.ProductId, item, Seckillconf.SecProductInfoMap, v)
		data = append(data, item)
	}

	return
}

func SecInfo(productId int) (data []map[string]interface{}, code int, err error) {

	Seckillconf.RWSecProductLock.RLock()
	defer Seckillconf.RWSecProductLock.RUnlock()

	item, code, err := SecInfoById(productId)
	if err != nil {
		return
	}

	data = append(data, item)
	return
}

func SecInfoById(productId int) (data map[string]interface{}, code int, err error) {

	Seckillconf.RWSecProductLock.RLock()
	defer Seckillconf.RWSecProductLock.RUnlock()

	v, ok := Seckillconf.SecProductInfoMap[productId]
	if !ok {
		code = ErrNotFoundProductId
		err = fmt.Errorf("not found product_id:%d", productId)
		return
	}

	start := false
	end := false
	status := "success"

	now := time.Now().Unix()
	if now-v.StartTime < 0 {
		start = false
		end = false
		status = "sec kill is not start"
		code = ErrActiveNotStart
	}

	if now-v.StartTime > 0 {
		start = true
	}

	if now-v.EndTime > 0 {
		start = false
		end = true
		status = "sec kill is already end"
		code = ErrActiveAlreadyEnd
	}

	if v.Status == ProductStatusForceSaleOut || v.Status == ProductStatusSaleOut {
		start = false
		end = true
		status = "product is sale out"
		code = ErrActiveSaleOut
	}

	data = make(map[string]interface{})
	data["product_id"] = productId
	data["start"] = start
	data["end"] = end
	data["status"] = status

	return
}

func userCheck(req *SecRequest) (err error) {
	found := false
	for _, refer := range Seckillconf.ReferWhiteList {
		if refer == req.ClientRefence {
			found = true
			break
		}
	}

	if !found {
		err = fmt.Errorf("invalid request")
		logs.Warn("user[%d] is reject by refer, req[%v]", req.UserId, req)
		return
	}

	authData := fmt.Sprintf("%d:%s", req.UserId, Seckillconf.CookieSecretKey)
	authSign := fmt.Sprintf("%x", md5.Sum([]byte(authData)))

	if authSign != req.UserAuthSign {
		err = fmt.Errorf("invalid user cookie auth")
		return
	}
	return
}

func SecKill(req *SecRequest) (data map[string]interface{}, code int, err error) {

	Seckillconf.RWSecProductLock.RLock()
	defer Seckillconf.RWSecProductLock.RUnlock()
	/*
		err = userCheck(req)
		if err != nil {
			code = ErrUserCheckAuthFailed
			logs.Warn("userId[%d] invalid, check failed, req[%v]", req.UserId, req)
			return
		}
	*/
	logs.Debug("this is service SecKill")
	/*
		err = userCheck(req)
		if err != nil {
			code = ErrUserServiceBusy
			logs.Warn("userId[%d] invalid, check failed, req[%v]", req.UserId, req)
			return
		}

	*/

	err = antiSpam(req)
	if err != nil {
		code = ErrUserServiceBusy
		logs.Warn("userId[%d] invalid, check failed, req[%v]", req.UserId, req)
		return
	}
	logs.Debug("start SecInfoById")
	data, code, err = SecInfoById(req.ProductId)
	if err != nil {
		logs.Warn("userId[%d] secInfoBy Id failed, req[%v]", req.UserId, req)
		return
	}

	if code != 0 {
		logs.Warn("userId[%d] secInfoByid failed, code[%d] req[%v]", req.UserId, code, req)
		return
	}

	userKey := fmt.Sprintf("%s_%s", req.UserId, req.ProductId)
	Seckillconf.UserConnMap[userKey] = req.ResultChan
	logs.Debug("2222")
	Seckillconf.SecReqChan <- req

	ticker := time.NewTicker(time.Second * 10)

	defer func() {
		ticker.Stop()
		Seckillconf.UserConnMapLock.Lock()
		delete(Seckillconf.UserConnMap, userKey)
		Seckillconf.UserConnMapLock.Unlock()
	}()

	select {
	case <-ticker.C:
		code = ErrProcessTimeout
		err = fmt.Errorf("request timeout")

		return
	case <-req.CloseNotify:
		code = ErrClientClosed
		err = fmt.Errorf("client already closed")
		return
	case result := <-req.ResultChan:
		code = result.Code
		data["product_id"] = result.ProductId
		data["token"] = result.Token
		data["user_id"] = result.UserId

		return
	}

	return
}

func SyncIdBlackList() {
	for {
		conn := Seckillconf.blackRedisPool.Get()
		defer conn.Close()
		reply, err := conn.Do("BLPOP", "blackidlist", time.Second)
		id, err := redis.Int(reply, err)
		if err != nil {
			continue
		}

		Seckillconf.RWBlackLock.Lock()
		Seckillconf.IdBlackMap[id] = true
		Seckillconf.RWBlackLock.Unlock()

		logs.Info("sync id list from redis succ, ip[%v]", id)
	}
}

func SyncIpBlackList() {
	var ipList []string
	lastTime := time.Now().Unix()
	for {
		conn := Seckillconf.blackRedisPool.Get()
		defer conn.Close()
		reply, err := conn.Do("BLPOP", "blackiplist", time.Second)
		ip, err := redis.String(reply, err)
		if err != nil {
			continue
		}

		curTime := time.Now().Unix()
		ipList = append(ipList, ip)

		if len(ipList) > 100 || curTime-lastTime > 5 {
			Seckillconf.RWBlackLock.Lock()
			for _, v := range ipList {
				Seckillconf.IpBlackMap[v] = true
			}
			Seckillconf.RWBlackLock.Unlock()

			lastTime = curTime
			logs.Info("sync ip list from redis succ, ip[%v]", ipList)
		}
	}
}
