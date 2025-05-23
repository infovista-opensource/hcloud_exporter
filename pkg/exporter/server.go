package exporter

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/promhippie/hcloud_exporter/pkg/config"
)

// ServerCollector collects metrics about the servers.
type ServerCollector struct {
	client   *hcloud.Client
	logger   *slog.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	config   config.Target

	Running         *prometheus.Desc
	Created         *prometheus.Desc
	IncludedTraffic *prometheus.Desc
	OutgoingTraffic *prometheus.Desc
	IngoingTraffic  *prometheus.Desc
	Cores           *prometheus.Desc
	Memory          *prometheus.Desc
	Disk            *prometheus.Desc
	Backup          *prometheus.Desc
	PriceHourly     *prometheus.Desc
	PriceMonthly    *prometheus.Desc
}

// NewServerCollector returns a new ServerCollector.
func NewServerCollector(logger *slog.Logger, client *hcloud.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *ServerCollector {
	if failures != nil {
		failures.WithLabelValues("server").Add(0)
	}

	labels := []string{"id", "name", "datacenter", "network"}
	pricingLabels := append(labels, "vat")
	return &ServerCollector{
		client:   client,
		logger:   logger.With("collector", "server"),
		failures: failures,
		duration: duration,
		config:   cfg,

		Running: prometheus.NewDesc(
			"hcloud_server_running",
			"If 1 the server is running, 0 otherwise",
			labels,
			nil,
		),
		Created: prometheus.NewDesc(
			"hcloud_server_created_timestamp",
			"Timestamp when the server have been created",
			labels,
			nil,
		),
		IncludedTraffic: prometheus.NewDesc(
			"hcloud_server_included_traffic_bytes",
			"Included traffic for the server in bytes",
			labels,
			nil,
		),
		OutgoingTraffic: prometheus.NewDesc(
			"hcloud_server_outgoing_traffic_bytes",
			"Outgoing traffic from the server in bytes",
			labels,
			nil,
		),
		IngoingTraffic: prometheus.NewDesc(
			"hcloud_server_incoming_traffic_bytes",
			"Ingoing traffic to the server in bytes",
			labels,
			nil,
		),
		Cores: prometheus.NewDesc(
			"hcloud_server_cores",
			"Server number of cores",
			labels,
			nil,
		),
		Memory: prometheus.NewDesc(
			"hcloud_server_memory_bytes",
			"Server memory in bytes",
			labels,
			nil,
		),
		Disk: prometheus.NewDesc(
			"hcloud_server_disk_bytes",
			"Server disk in bytes",
			labels,
			nil,
		),
		Backup: prometheus.NewDesc(
			"hcloud_server_backup",
			"If 1 server backups are enabled, 0 otherwise",
			labels,
			nil,
		),
		PriceHourly: prometheus.NewDesc(
			"hcloud_server_price_hourly",
			"Price of the server billed hourly in €",
			pricingLabels,
			nil,
		),
		PriceMonthly: prometheus.NewDesc(
			"hcloud_server_price_monthly",
			"Price of the server billed monthly in €",
			pricingLabels,
			nil,
		),
	}
}

// Metrics simply returns the list metric descriptors for generating a documentation.
func (c *ServerCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.Running,
		c.Created,
		c.IncludedTraffic,
		c.OutgoingTraffic,
		c.IngoingTraffic,
		c.Cores,
		c.Memory,
		c.Disk,
		c.Backup,
		c.PriceHourly,
		c.PriceMonthly,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *ServerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Running
	ch <- c.Created
	ch <- c.IncludedTraffic
	ch <- c.OutgoingTraffic
	ch <- c.IngoingTraffic
	ch <- c.Cores
	ch <- c.Memory
	ch <- c.Disk
	ch <- c.Backup
	ch <- c.PriceHourly
	ch <- c.PriceMonthly
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *ServerCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	now := time.Now()
	servers, err := c.client.Server.All(ctx)

	if err != nil {
		c.logger.Error("Failed to fetch servers",
			"err", err,
		)

		c.failures.WithLabelValues("server").Inc()
		return
	}

	c.logger.Debug("Fetched servers",
		"count", len(servers),
	)

	for _, server := range servers {
		var (
			running float64
			network string
		)

		// Set network label based on whether server is in a private network
		// If it has private networks, use the ID of the first one, otherwise use "public"
		if len(server.PrivateNet) > 0 {
			network = strconv.FormatInt(server.PrivateNet[0].Network.ID, 10)
			c.logger.Info("Server is in private network",
				"server", server.Name,
				"network ID", network,
			)
		} else {
			network = "public"
		}

		labels := []string{
			strconv.FormatInt(server.ID, 10),
			server.Name,
			server.Datacenter.Name,
			network,
		}

		if server.Status == "running" {
			running = 1.0
		}

		ch <- prometheus.MustNewConstMetric(
			c.Running,
			prometheus.GaugeValue,
			running,
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Created,
			prometheus.GaugeValue,
			float64(server.Created.Unix()),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.IncludedTraffic,
			prometheus.GaugeValue,
			float64(server.IncludedTraffic),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.OutgoingTraffic,
			prometheus.GaugeValue,
			float64(server.OutgoingTraffic),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.IngoingTraffic,
			prometheus.GaugeValue,
			float64(server.IngoingTraffic),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Cores,
			prometheus.GaugeValue,
			float64(server.ServerType.Cores),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Memory,
			prometheus.GaugeValue,
			float64(server.ServerType.Memory*1024*1024),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Disk,
			prometheus.GaugeValue,
			float64(server.ServerType.Disk*1024*1024),
			labels...,
		)

		backupEnabled := 0.0
		if server.BackupWindow != "" {
			backupEnabled = 1.0
		}

		ch <- prometheus.MustNewConstMetric(
			c.Backup,
			prometheus.GaugeValue,
			backupEnabled,
			labels...,
		)

		// Create correct pricing labels that include both network and vat
		pricingLabelsNet := []string{
			strconv.FormatInt(server.ID, 10),
			server.Name,
			server.Datacenter.Name,
			network,
			"net",
		}
		pricingLabelsGross := []string{
			strconv.FormatInt(server.ID, 10),
			server.Name,
			server.Datacenter.Name,
			network,
			"gross",
		}

		for _, pricing := range server.ServerType.Pricings {
			if server.Datacenter.Location.Name == pricing.Location.Name {
				if net, err := strconv.ParseFloat(pricing.Hourly.Net, 64); err != nil {
					c.logger.Error("Failed to parse hourly server type net costs",
						"name", server.Name,
						"err", err,
					)

					c.failures.WithLabelValues("server").Inc()
				} else {
					ch <- prometheus.MustNewConstMetric(
						c.PriceHourly,
						prometheus.GaugeValue,
						net,
						pricingLabelsNet...,
					)
				}

				if gross, err := strconv.ParseFloat(pricing.Hourly.Gross, 64); err != nil {
					c.logger.Error("Failed to parse hourly server type gross costs",
						"name", server.Name,
						"err", err,
					)

					c.failures.WithLabelValues("server").Inc()
				} else {
					ch <- prometheus.MustNewConstMetric(
						c.PriceHourly,
						prometheus.GaugeValue,
						gross,
						pricingLabelsGross...,
					)
				}

				if net, err := strconv.ParseFloat(pricing.Monthly.Net, 64); err != nil {
					c.logger.Error("Failed to parse monthly server type net costs",
						"name", server.Name,
						"err", err,
					)

					c.failures.WithLabelValues("server").Inc()
				} else {
					ch <- prometheus.MustNewConstMetric(
						c.PriceMonthly,
						prometheus.GaugeValue,
						net,
						pricingLabelsNet...,
					)
				}

				if gross, err := strconv.ParseFloat(pricing.Monthly.Gross, 64); err != nil {
					c.logger.Error("Failed to parse monthly server type gross costs",
						"name", server.Name,
						"err", err,
					)

					c.failures.WithLabelValues("server").Inc()
				} else {
					ch <- prometheus.MustNewConstMetric(
						c.PriceMonthly,
						prometheus.GaugeValue,
						gross,
						pricingLabelsGross...,
					)
				}
			}
		}
	}

	c.logger.Debug("Processed server collector",
		"duration", time.Since(now),
	)

	c.duration.WithLabelValues("server").Observe(time.Since(now).Seconds())
}
