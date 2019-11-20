package redis

import (
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gvar"
	"github.com/gomodule/redigo/redis"
	"time"
)

const (
	DEFAULT_POOL_IDLE_TIMEOUT = 60 * time.Second
	DEFAULT_POOL_CONN_TIMEOUT = 10 * time.Second
	DEFAULT_POOL_MAX_LIFETIME = 60 * time.Second
)

type Redis struct {
	pool *redis.Pool
	group string
	config Config
}


type Conn struct {
	redis.Conn
}

type Config struct {
	Host string
	Port int
	Db int
	Passwd string
	MaxIdle int
	MaxActive int
	IdleTimeout time.Duration
	MaxConnLifetime time.Duration
	ConnectTimeout time.Duration
}

type PoolStats struct {
	redis.PoolStats
}


var (
	pools = gmap.NewStrAnyMap(true)
)


func New(config Config) *Redis {
	if config.IdleTimeout == 0 {
		config.IdleTimeout = DEFAULT_POOL_IDLE_TIMEOUT
	}
	if config.ConnectTimeout == 0 {
		config.ConnectTimeout = DEFAULT_POOL_CONN_TIMEOUT
	}
	if config.MaxConnLifetime == 0 {
		config.MaxConnLifetime = DEFAULT_POOL_MAX_LIFETIME
	}

	return &Redis{
		config: config,
		pool: pools.GetOrSetFuncLock(fmt.Sprintf("%v", config), func() interface{} {
			return &redis.Pool{
				IdleTimeout: config.IdleTimeout,
				MaxActive: config.MaxActive,
				MaxIdle: config.MaxIdle,
				MaxConnLifetime: config.MaxConnLifetime,
				Dial: func() (conn redis.Conn, err error) {
					conn, err = redis.Dial(
						"tcp",
						fmt.Sprintf("%s:%d", config.Host, config.Port),
						redis.DialConnectTimeout(config.ConnectTimeout),
						)
					if err != nil {
						return
					}
					// auth
					if len(config.Passwd) > 0 {
						if _, err = conn.Do("AUTH", config.Passwd); err != nil {
							return
						}
					}
					// db
					if _, err = conn.Do("SELECT", config.Db); err != nil {
						return
					}
					return
				},

				TestOnBorrow: func(c redis.Conn, t time.Time) error {
					_, err := c.Do("PING")
					return err
				},
			}
		}).(*redis.Pool),
	}
}

func (r *Redis) Close() error {
	if r.group != "" {
		instances.Remove(r.group)
	}
	pools.Remove(fmt.Sprintf("%v", r.config))
	return r.pool.Close()
}

func (r *Redis) Conn() *Conn {
	return &Conn{r.pool.Get()}
}

func (r *Redis) GetConn() *Conn {
	return r.Conn()
}

func (r *Redis) SetMaxIdle(value int) {
	r.pool.MaxIdle = value
}

func (r *Redis) SetMaxActive(value int) {
	r.pool.MaxActive = value
}

func (r *Redis) SetIdleTimeout(value time.Duration) {
	r.pool.IdleTimeout = value
}


func (r *Redis) SetMaxConnLifetime(value time.Duration) {
	r.pool.MaxConnLifetime = value
}

func (r *Redis) Stats() *PoolStats {
	return &PoolStats{r.pool.Stats()}
}

func (r *Redis) Do(command string, args ...interface{}) (interface{}, error) {
	conn := &Conn{r.pool.Get()}
	defer conn.Close()
	return conn.Do(command, args...)
}


func (r *Redis) DoVar(command string, args ...interface{}) (*gvar.Var, error) {
	v, err := r.Do(command, args...)
	return gvar.New(v), err
}
