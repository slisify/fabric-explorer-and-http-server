package prometheus

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/syncx"
	"github.com/tal-tech/go-zero/core/threading"
)

var (
	once    sync.Once
	enabled syncx.AtomicBool
)

// Enabled returns if prometheus is enabled.
func Enabled() bool {
	return enabled.True()
}

// StartAgent starts a prometheus agent.
func StartAgent(c Config) {
	once.Do(func() {
		if len(c.Host) == 0 {
			return
		}

		enabled.Set(true)
		threading.GoSafe(func() {
			http.Handle(c.Path, promhttp.Handler())
			addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
			logx.Infof("Starting prometheus agent at %s", addr)
			if err := http.ListenAndServe(addr, nil); err != nil {
				logx.Error(err)
			}
		})
	})
}
