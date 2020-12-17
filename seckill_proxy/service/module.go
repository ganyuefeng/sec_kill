package service

import (
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

var (
	Seckillconf = &SeckillConf{}
)

type RedisConf struct {
	Redis_addr        string
	Redis_maxIdle     int
	Redis_maxActive   int
	Redis_idleTimeout int
}
type EtcdConf struct {
	EtcdAddr          string
	Timeout           int
	EtcdSecKeyPrefix  string
	EtcdSecProductKey string
}

type SecInfoConf struct {
	ProductId int
	StartTime int64
	EndTime   int64
	Status    int
	Total     int
	Left      int
}

type SeckillConf struct {
	RedisBlackConf       RedisConf
	RedisProxy2LayerConf RedisConf
	RedisLayer2ProxyConf RedisConf

	EtcdConf          EtcdConf
	LogPath           string
	LogLevel          string
	SecProductInfoMap map[int]*SecInfoConf
	RWSecProductLock  sync.RWMutex
	CookieSecretKey   string

	ReferWhiteList []string

	IpBlackMap map[string]bool
	IdBlackMap map[int]bool

	AccessLimitConf      AccessLimitConf
	blackRedisPool       *redis.Pool
	proxy2LayerRedisPool *redis.Pool
	layer2ProxyRedisPool *redis.Pool

	SecLimitMgr *SecLimitMgr

	RWBlackLock                  sync.RWMutex
	WriteProxy2LayerGoroutineNum int
	ReadProxy2LayerGoroutineNum  int

	SecReqChan     chan *SecRequest
	SecReqChanSize int

	UserConnMap     map[string]chan *SecResult
	UserConnMapLock sync.Mutex
}

const (
	ProductStatusNormal       = 0
	ProductStatusSaleOut      = 1
	ProductStatusForceSaleOut = 2
)

type SecResult struct {
	ProductId int
	UserId    int
	Code      int
	Token     string
}

type SecRequest struct {
	ProductId     int
	Source        string
	AuthCode      string
	SecTime       string
	Nance         string
	UserId        int
	UserAuthSign  string
	AccessTime    time.Time
	ClientAddr    string
	ClientRefence string
	CloseNotify   <-chan bool `json:"-"`

	ResultChan chan *SecResult `json:"-"`
}

type AccessLimitConf struct {
	IPSecAccessLimit   int
	UserSecAccessLimit int
	IPMinAccessLimit   int
	UserMinAccessLimit int
}
