package plugin

import (
	"time"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth"
	"net/http"
	"github.com/loveczp/fqimg/lib"
	"github.com/pkg/errors"
	"strconv"
)

func Plugin_throttle_total(h http.HandlerFunc) http.HandlerFunc {

	var Total_limit *limiter.Limiter
	Total_limit = tollbooth.NewLimiter(float64(lib.Conf.UploadThrottleTotal), &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	Total_limit.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}).SetOnLimitReached(func(w http.ResponseWriter, r *http.Request) {
		lib.WriteErr(w, http.StatusBadRequest, errors.New("upload rate is tow high, exceed the total limit  "+strconv.Itoa(lib.Conf.UploadThrottleTotal)+" upload / second+IP"))
	})
	return func(writer http.ResponseWriter, request *http.Request) {
		if lib.Conf.UploadThrottleTotal != 0 {
			tollbooth.LimitFuncHandler(Total_limit, h.ServeHTTP)
		} else {
			h.ServeHTTP(writer, request)
		}
	}
}

func Plugin_throttle_ip(h http.HandlerFunc) http.HandlerFunc {
	var IP_limit *limiter.Limiter
	IP_limit = tollbooth.NewLimiter(float64(lib.Conf.UploadThrottlePerIp), &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	IP_limit.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}).
		SetOnLimitReached(func(w http.ResponseWriter, r *http.Request) {
		lib.WriteErr(w, http.StatusBadRequest, errors.New("upload rate is tow high, exceed the ip limit  "+strconv.Itoa(lib.Conf.UploadThrottlePerIp)+" upload / second+IP"))
	})

	return func(writer http.ResponseWriter, request *http.Request) {
		if lib.Conf.UploadThrottlePerIp != 0 {
			tollbooth.LimitFuncHandler(IP_limit, h.ServeHTTP)
		} else {
			h.ServeHTTP(writer, request)
		}
	}
}
