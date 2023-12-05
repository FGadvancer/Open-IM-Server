package prommetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	OnlineUserGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "online_user_num",
		Help: "The number of online user num",
	})
	GateWaySendMsgTotalCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gateway_send_msg_total",
		Help: "The number of gateway send msg total",
	})
)
