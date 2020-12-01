package resources

import (
	"WiFiAuth/internal/config"
	"WiFiAuth/internal/model"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"time"
)

type R struct {
}

func New(logger *zap.SugaredLogger) (*R, error) {

	postgresCon, err := sqlx.Connect("postgres", config.Config.POSTGRESURL)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	config.Config.POSTGRESDB = postgresCon

	config.Config.REDISPOOL = newPool()
	conn := config.Config.REDISPOOL.Get()

	err = ping(conn)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	go func() {
		for {
			model.UpdateStateCallNumber()
			time.Sleep(time.Minute * 5)
		}
	}()

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
			c, err := redis.Dial("tcp", config.Config.REDISURL)
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
