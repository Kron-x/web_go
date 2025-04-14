package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
)

var (
    PageVisits = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "page_visits_total",
            Help: "Total number of page visits",
        },
        []string{"path"},
    )

    ActiveUsers = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_users_current",
            Help: "Current number of active users",
        },
    )
)

func Init() {
    prometheus.MustRegister(PageVisits)
    prometheus.MustRegister(ActiveUsers)
}