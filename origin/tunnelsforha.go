package origin

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// tunnelsForHA maps this cloudflared instance's HA connections to the tunnel IDs they serve.
type tunnelsForHA struct {
	sync.Mutex
	metrics *prometheus.GaugeVec
	entries map[uint8]string
}

// NewTunnelsForHA initializes the Prometheus metrics etc for a tunnelsForHA.
func NewTunnelsForHA() tunnelsForHA {
	metrics := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tunnel_ids",
			Help: "The ID of all tunnels (and their corresponding HA connection ID) running in this instance of cloudflared.",
		},
		[]string{"tunnel_id", "ha_conn_id"},
	)
	prometheus.MustRegister(metrics)

	return tunnelsForHA{
		metrics: metrics,
		entries: make(map[uint8]string),
	}
}

// Track a new tunnel ID, removing the disconnected tunnel (if any) and update metrics.
func (t *tunnelsForHA) AddTunnelID(haConn uint8, tunnelID string) {
	t.Lock()
	defer t.Unlock()
	if oldTunnelID, ok := t.entries[haConn]; ok {
		t.metrics.WithLabelValues(oldTunnelID).Dec()
	}
	t.entries[haConn] = tunnelID
	t.metrics.WithLabelValues(tunnelID, fmt.Sprintf("%v", haConn)).Inc()
}

func (t *tunnelsForHA) String() string {
	t.Lock()
	defer t.Unlock()
	return fmt.Sprintf("%v", t.entries)
}