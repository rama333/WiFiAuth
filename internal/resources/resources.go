package resources

import (
	"WiFiAuth/internal/config"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"log"
)

type R struct {
}

func New(logger *zap.SugaredLogger) (*R, error) {

	config.Config.REDISPOOL = newPool()
	conn := config.Config.REDISPOOL.Get()

	err := ping(conn)
	if err != nil {
		return nil, config.Config.REDISPOOLERR
	}

	return &R{}, nil
}

func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "192.168.1.3:6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func ping(c redis.Conn) error {
	// Send PING command to Redis
	pong, err := c.Do("PING")
	if err != nil {
		return err
	}

	// PING command returns a Redis "Simple String"
	// Use redis.String to convert the interface type to string
	s, err := redis.String(pong, err)
	if err != nil {
		return err
	}

	log.Printf("PING Response = %s\n", s)
	// Output: PONG

	return nil
}

func (r *R) Release() error {

	return config.Config.REDISPOOL.Close()
}
