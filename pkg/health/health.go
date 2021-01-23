package health

import (
	"encoding/json"
	"github.com/douyu/jupiter/pkg/server/governor"
	"net/http"
)

const (
	StatusName = "status"
	StatusUp   = "up"
	StatusDown = "down"
)

func Init() {
	governor.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		rets := make(map[string]interface{})
		status := StatusUp

		for name := range RegisterMap {
			ok, ret := RegisterMap[name]()
			subStatus := StatusUp
			if !ok {
				subStatus = StatusDown
				status = StatusDown
			}
			if ret == nil {
				ret = make(map[string]interface{})
			}
			ret[StatusName] = subStatus
			rets[name] = ret
		}

		rets[StatusName] = status

		body, _ := json.Marshal(&rets)
		if status == StatusDown {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		w.Write(body)

	})
}
