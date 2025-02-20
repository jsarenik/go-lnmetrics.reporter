package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/vincenzopalazzo/glightning/jrpc2"
)

type MetricOneRpcMethod struct {
	StartPeriod string `json:"start"`
	EndPeriod   string `json:"end"`

	// Metric Reference
	plugin *MetricsPlugin `json:"-"`
}

func (instance *MetricOneRpcMethod) Name() string {
	return "metric_one"
}

func NewMetricPlugin(plugin *MetricsPlugin) *MetricOneRpcMethod {
	return &MetricOneRpcMethod{
		StartPeriod: "",
		EndPeriod:   "",
		plugin:      plugin,
	}
}

func (instance *MetricOneRpcMethod) New() interface{} {
	return NewMetricPlugin(instance.plugin)
}

func (instance *MetricOneRpcMethod) Call() (jrpc2.Result, error) {
	metricOne, found := instance.plugin.Metrics[1]

	if !found {
		return nil, fmt.Errorf("Metric with id %d not found", 1)
	}

	if instance.StartPeriod == "" &&
		instance.EndPeriod == "" {
		return nil, fmt.Errorf("Missing at list the start parameter in the rpc method")
	}

	if instance.StartPeriod == "now" {
		return metricOne, nil
	}

	if instance.StartPeriod == "last" {
		jsonValue, err := instance.plugin.Storage.LoadLastMetricOne()
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(*jsonValue), &metricOne); err != nil {
			return nil, err
		}
		return metricOne, nil
	}

	return nil, fmt.Errorf("We don't support the filter operation right now")
}

// ForceUpdateRPC enable the force update command
type ForceUpdateRPC struct {
	// the instance of the plugin
	plugin *MetricsPlugin
}

func NewForceUpdateRPC(plugin *MetricsPlugin) *ForceUpdateRPC {
	return &ForceUpdateRPC{plugin}
}

func (instance *ForceUpdateRPC) New() interface{} {
	return instance
}

func (instance *ForceUpdateRPC) Name() string {
	return "lnmetrics-force-update"
}

func (instance *ForceUpdateRPC) Call() (jrpc2.Result, error) {
	for _, metric := range instance.plugin.Metrics {
		msg := Msg{
			cmd:    "plugin_rpc_method",
			params: map[string]interface{}{"event": "on_force_update"},
		}
		instance.plugin.callUpdateOnMetric(metric, &msg)
	}
	response := struct {
		result string
	}{result: "force call update on all the metrics succeeded"}
	return response, nil
}
