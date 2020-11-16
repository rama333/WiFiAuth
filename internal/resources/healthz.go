package resources

import (
	"WiFiAuth/internal/diagnostics/healthz"
	"time"
)

func (r *R) Healthz() []healthz.Resource {
	var err error

	dbStatus := healthz.Ok
	dbMsg := "It works!"
	for i := 0; i < 5; i++ {
		//_, err = config.Config.DB.Query("SELECT 1")
		//if err == nil {
		//	break
		//}
		time.Sleep(time.Second)
	}
	if err != nil {
		dbStatus = healthz.Fatal
		dbMsg = err.Error()
	}

	return []healthz.Resource{
		{
			Name:    "reformDB",
			Status:  dbStatus,
			Message: dbMsg,
		},
	}
}
