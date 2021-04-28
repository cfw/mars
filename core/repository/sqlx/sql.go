package sqlx

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"time"

	"xorm.io/xorm"
)

const (
	healthCheckInterval = 10
	pingTimeout         = 10
	maxRetry            = 3
	pingCommand         = "select 1"
)

func NewSql(c *Config) *xorm.Engine {
	var err error
	log.Infof("Init db %s", c.Database)
	engine, err := xorm.NewEngine(c.Driver, c.Url())
	if err != nil {
		panic(err)
	}

	engine.SetLogger(NewLogCtx(log.StandardLogger()))
	if c.Debug {
		engine.ShowSQL(true)
	}

	if err = engine.Ping(); err != nil {
		panic(err)
	}

	go func() {
		ticker := time.NewTicker(time.Duration(healthCheckInterval) * time.Second)
		for range ticker.C {
			for i := 0; i < maxRetry; i++ {
				if err := ping(engine); err != nil {
					if i == 2 {
						log.Errorf("Db health check failed after retry %d times", maxRetry)
						break
					}
					log.Errorf("Db health check retry")
					continue
				}
				break
			}
		}
	}()
	log.Infof("Init db connection %s success", c.Database)
	return engine
}

func ping(engine *xorm.Engine) error {
	errCh := make(chan error)
	go func() {
		if _, err := engine.Exec(pingCommand); err != nil {
			log.Errorf("Engine ping error %v", err.Error())
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-time.After(time.Duration(pingTimeout) * time.Second):
		log.Error("Db ping connection timeout")
		return errors.New("db ping connection timeout")
	}
	return nil
}
