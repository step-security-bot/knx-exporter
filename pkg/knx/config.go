package knx

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Config defines the structure of the configuration file which defines which
// KNX Group Addresses were mapped into prometheus metrics.
type Config struct {
	Connection Connection `json:",omitempty"`
	// MetricsPrefix is a short prefix which will be added in front of the actual metric name.
	MetricsPrefix  string
	AddressConfigs map[GroupAddress]GroupAddressConfig
}

type Connection struct {
	Type     ConnectionType
	Endpoint string
}

type ConnectionType string

const Tunnel = ConnectionType("Tunnel")
const Router = ConnectionType("Router")

func (t ConnectionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(t))
}

func (t *ConnectionType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	switch strings.ToLower(str) {
	case "tunnel":
		*t = Tunnel
	case "router":
		*t = Router
	default:
		return fmt.Errorf("invalid connection type given: \"%s\"", str)
	}
	return nil
}

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	duration, err := time.ParseDuration(str)
	if err != nil {
		return err
	}
	*d = Duration(duration)
	return nil
}

// GroupAddressConfig defines all information to map a KNX group address to a prometheus metric.
type GroupAddressConfig struct {
	// Name defines the prometheus metric name without the MetricsPrefix.
	Name string
	// Comment to identify the group address.
	Comment string `json:",omitempty"`
	// DPT defines the DPT at the knx bus. This is required to parse the values correctly.
	DPT string
	// MetricType is the type that prometheus uses when exporting it. i.e. gauge or counter
	MetricType string
	// Export the metric to prometheus
	Export bool
	// ReadActive allows the exporter to actively send `GroupValueRead` telegrams to actively poll the value instead waiting for it.
	ReadActive bool `json:",omitempty"`
	// MaxAge of a value until it will actively send a `GroupValueRead` telegram to read the value if ReadActive is set to true.
	MaxAge Duration `json:",omitempty"`
}
