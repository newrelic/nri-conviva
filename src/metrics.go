package main

import (
	"time"

	sdk_metric "github.com/newrelic/infra-integrations-sdk/v4/data/metric"
	"github.com/newrelic/infra-integrations-sdk/v4/integration"
	sdk_log "github.com/newrelic/infra-integrations-sdk/v4/log"
	"github.com/newrelic/nri-conviva/src/api"
)

const (
	METRIC_PREFIX = "conviva."
	PERCENTAGE_SUFFIX = ".percentage"
)

type GetCountMetricFunc func (m *api.Metrics) *api.Count
type GetGaugeMetricFunc func (m *api.Metrics) *api.Gauge
type GetCountPercentageFunc func (m *api.Metrics) *api.CountPercentage

type AddCountFunc func (metricName string, c int64) error
type AddGaugeFunc func (metricName string, g float64) error
type CreateMetricsFunc func (*api.Metrics, AddCountFunc, AddGaugeFunc) error

type AddMetricFunc func(
	entity        *integration.Entity,
	timestamp     time.Time,
	dimension     *api.Dimension,
	Metrics       *api.Metrics,
) error

var (
	metricAdders = []AddMetricFunc{}
)

func PercentageToGauge(p *api.Percentage) *api.Gauge {
	if p != nil {
		return &api.Gauge{Value: p.Value}
	}
	return nil
}

func BitrateToGauge(b *api.Bitrate) *api.Gauge {
	if b != nil {
		return &api.Gauge{Value: b.Bps}
	}
	return nil
}

func FramerateToGauge(f *api.Framerate) *api.Gauge {
	if f != nil {
		return &api.Gauge{Value: f.Fps}
	}
	return nil
}

func RatioToGauge(r *api.Ratio) *api.Gauge {
	if r != nil {
		return &api.Gauge{Value: r.Ratio}
	}
	return nil
}

var addAdEndedPlays CreateMetricsFunc = func (
	m *api.Metrics,
	addCount AddCountFunc,
	addGauge AddGaugeFunc,
) error {
	r := m.AdEndedPlays
	if r != nil {
		err := addCount("ad_ended_plays", r.Count.Value)
		if err != nil {
			return err
		}
		return addGauge("ad_ended_plays.per_unique_device", r.PerUniqueDevice)
	}
	return nil
}

var addAdMinutesPlayed CreateMetricsFunc = func (
	m *api.Metrics,
	addCount AddCountFunc,
	addGauge AddGaugeFunc,
) error {
	r := m.AdMinutesPlayed
	if r != nil {
		err := addCount("ad_minutes_played", r.Count.Value)
		if err != nil {
			return err
		}
		return addGauge("ad_minutes_played.per_ended_play", r.PerEndedPlay)
	}
	return nil
}

var addEndedPlays CreateMetricsFunc = func (
	m *api.Metrics,
	addCount AddCountFunc,
	addGauge AddGaugeFunc,
) error {
	r := m.EndedPlays
	if r != nil {
		err := addCount("ended_plays", r.Count.Value)
		if err != nil {
			return err
		}
		return addGauge("ended_plays.per_unique_device", r.PerUniqueDevice)
	}
	return nil
}

var addMinutesPlayed CreateMetricsFunc = func (
	m *api.Metrics,
	addCount AddCountFunc,
	addGauge AddGaugeFunc,
) error {
	r := m.MinutesPlayed
	if r != nil {
		err := addCount("minutes_played", r.Count.Value)
		if err != nil {
			return err
		}
		return addGauge("minutes_played.per_ended_play", r.PerEndedPlay)
	}
	return nil
}

func newMetric(
	entity *integration.Entity,
	dimension *api.Dimension,
	metric sdk_metric.Metric,
) error {
	if dimension != nil {
		err := metric.AddDimension(dimension.Key, dimension.Value)
		if err != nil {
			return err
		}
	}

	entity.AddMetric(metric)

	return nil
}

func newCountMetric(
	entity        *integration.Entity,
	timestamp     time.Time,
	metricName    string,
	count         int64,
	dimension     *api.Dimension,
) error {
	metric, err := sdk_metric.NewCount(
		timestamp,
		METRIC_PREFIX + metricName,
		float64(count),
	)
	if err != nil {
		return err
	}

	return newMetric(entity, dimension, metric)
}

func newGaugeMetric(
	entity        *integration.Entity,
	timestamp     time.Time,
	metricName    string,
	value         float64,
	dimension     *api.Dimension,
) error {
	metric, err := sdk_metric.NewGauge(
		timestamp,
		METRIC_PREFIX + metricName,
		value,
	)
	if err != nil {
		return err
	}

	return newMetric(entity, dimension, metric)
}

func createCountFunc(metricName string, fn GetCountMetricFunc) AddMetricFunc {
	return func (
		e *integration.Entity,
		t time.Time,
		d *api.Dimension,
		m *api.Metrics,
	) error {
		c := fn(m)
		if c != nil {
			return newCountMetric(
				e,
				t,
				metricName,
				c.Value,
				d,
			)
		}

		return nil
	}
}

func createGaugeFunc(metricName string, fn GetGaugeMetricFunc) AddMetricFunc {
	return func (
		e *integration.Entity,
		t time.Time,
		d *api.Dimension,
		m *api.Metrics,
	) error {
		g := fn(m)
		if g != nil {
			return newGaugeMetric(
				e,
				t,
				metricName,
				g.Value,
				d,
			)
		}

		return nil
	}
}

func createMetricsFunc(fn CreateMetricsFunc) AddMetricFunc {
	return func (
		e *integration.Entity,
		t time.Time,
		d *api.Dimension,
		m *api.Metrics,
	) error {
		return fn(
			m,
			func (metricName string, c int64) error {
				return newCountMetric(
					e,
					t,
					metricName,
					c,
					d,
				)
			},
			func (metricName string, g float64) error {
				return newGaugeMetric(
					e,
					t,
					metricName,
					g,
					d,
				)
			},
		)
	 }
}

func createCountPercentageFunc(
	metricName string,
	fn GetCountPercentageFunc,
) AddMetricFunc {
	return createMetricsFunc(func (
		m *api.Metrics,
		addCount AddCountFunc,
		addGauge AddGaugeFunc,
	) error {
		cp := fn(m)
		if (cp != nil) {
			err := addCount(metricName, cp.Count.Value)
			if err != nil {
				return err
			}

			return addGauge(metricName + PERCENTAGE_SUFFIX, cp.Percentage.Value)
		}
		return nil
	})
}


func initMetrics() {
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"abandonment",
		func (m *api.Metrics) *api.CountPercentage {
			return m.Abandonment
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"abandonment_with_pre_roll",
		func (m *api.Metrics) *api.Gauge {
			return PercentageToGauge(m.AbandonmentWithPreRoll)
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"abandonment_without_pre_roll",
		func (m *api.Metrics) *api.Gauge {
			return PercentageToGauge(m.AbandonmentWithoutPreRoll)
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"ad_actual_duration",
		func (m *api.Metrics) *api.Gauge {
			return m.AdActualDuration
		},
	))
	metricAdders = append(metricAdders, createCountFunc(
		"ad_attempts",
		func (m *api.Metrics) *api.Count {
			return m.AdAttempts
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"ad_bitrate",
		func (m *api.Metrics) *api.Gauge {
			return BitrateToGauge(m.AdBitrate)
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"ad_completed_creative_plays",
		func (m *api.Metrics) *api.Gauge {
			return m.AdCompletedCreativePlays
		},
	))
	metricAdders = append(metricAdders, createCountFunc(
		"ad_concurrent_plays",
		func (m *api.Metrics) *api.Count {
			return m.AdConcurrentPlays
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"ad_connection_induced_rebuffering_ratio",
		func (m *api.Metrics) *api.Gauge {
			return RatioToGauge(m.AdConnectionInducedRebufferingRatio)
		},
	))
	metricAdders = append(metricAdders, createMetricsFunc(addAdEndedPlays))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"ad_exit_before_video_starts",
		func (m *api.Metrics) *api.CountPercentage {
			return m.AdExitBeforeVideoStarts
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"ad_framerate",
		func (m *api.Metrics) *api.Gauge {
			return FramerateToGauge(m.AdFramerate)
		},
	))
	metricAdders = append(metricAdders, createMetricsFunc(addAdMinutesPlayed))
	metricAdders = append(metricAdders, createGaugeFunc(
		"ad_percentage_complete",
		func (m *api.Metrics) *api.Gauge {
			return PercentageToGauge(m.AdPercentageComplete)
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"ad_plays",
		func (m *api.Metrics) *api.CountPercentage {
			return m.AdPlays
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"ad_rebuffering_ratio",
		func (m *api.Metrics) *api.Gauge {
			return RatioToGauge(m.AdRebufferingRatio)
		},
	))
	metricAdders = append(metricAdders, createCountFunc(
		"ad_unique_devices",
		func (m *api.Metrics) *api.Count {
			return m.AdUniqueDevices
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"ad_video_playback_failures",
		func (m *api.Metrics) *api.CountPercentage {
			return m.AdVideoPlaybackFailures
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"ad_video_restart_time",
		func (m *api.Metrics) *api.Gauge {
			return m.AdVideoRestartTime
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"ad_video_start_failures",
		func (m *api.Metrics) *api.CountPercentage {
			return m.AdVideoStartFailures
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"ad_video_start_time",
		func (m *api.Metrics) *api.Gauge {
			return m.AdVideoStartTime
		},
	))
	metricAdders = append(metricAdders, createCountFunc(
		"attempts",
		func (m *api.Metrics) *api.Count {
			return m.Attempts
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"attempts_with_pre_roll",
		func (m *api.Metrics) *api.CountPercentage {
			return m.AttemptsWithPreRoll
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"attempts_without_pre_roll",
		func (m *api.Metrics) *api.CountPercentage {
			return m.AttemptsWithoutPreRoll
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"bad_session",
		func (m *api.Metrics) *api.CountPercentage {
			return m.BadSession
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"bad_session_average_life_playing_time_mins",
		func (m *api.Metrics) *api.Gauge {
			return m.BadSessionAverageLifePlayingTimeMins
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"bad_unique_devices",
		func (m *api.Metrics) *api.CountPercentage {
			return m.BadUniqueDevices
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"bad_unique_viewers",
		func (m *api.Metrics) *api.CountPercentage {
			return m.BadUniqueViewers
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"bitrate",
		func (m *api.Metrics) *api.Gauge {
			return BitrateToGauge(m.Bitrate)
		},
	))
	metricAdders = append(metricAdders, createCountFunc(
		"concurrent_plays",
		func (m *api.Metrics) *api.Count {
			return m.ConcurrentPlays
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"connection_induced_rebuffering_ratio",
		func (m *api.Metrics) *api.Gauge {
			return RatioToGauge(m.ConnectionInducedRebufferingRatio)
		},
	))
	metricAdders = append(metricAdders, createMetricsFunc(addEndedPlays))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"ended_plays_with_ads",
		func (m *api.Metrics) *api.CountPercentage {
			return m.EndedPlaysWithAds
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"ended_plays_without_ads",
		func (m *api.Metrics) *api.CountPercentage {
			return m.EndedPlaysWithoutAds
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"exit_before_video_starts",
		func (m *api.Metrics) *api.CountPercentage {
			return m.ExitBeforeVideoStarts
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"framerate",
		func (m *api.Metrics) *api.Gauge {
			return FramerateToGauge(m.Framerate)
		},
	))
	metricAdders = append(metricAdders, createCountFunc(
		"good_session",
		func (m *api.Metrics) *api.Count {
			return m.GoodSession
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"good_session_average_life_playing_time_mins",
		func (m *api.Metrics) *api.Gauge {
			return m.GoodSessionAverageLifePlayingTimeMins
		},
	))
	metricAdders = append(metricAdders, createCountFunc(
		"good_unique_devices",
		func (m *api.Metrics) *api.Count {
			return m.GoodUniqueDevices
		},
	))
	metricAdders = append(metricAdders, createCountFunc(
		"good_unique_viewers",
		func (m *api.Metrics) *api.Count {
			return m.GoodUniqueViewers
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"high_rebuffering",
		func (m *api.Metrics) *api.CountPercentage {
			return m.HighRebuffering
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"high_rebuffering_with_ads",
		func (m *api.Metrics) *api.Gauge {
			return PercentageToGauge(m.HighRebufferingWithAds)
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"high_rebuffering_without_ads",
		func (m *api.Metrics) *api.Gauge {
			return PercentageToGauge(m.HighRebufferingWithoutAds)
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"high_startup_time",
		func (m *api.Metrics) *api.CountPercentage {
			return m.HighStartupTime
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"high_startup_time_with_pre_roll",
		func (m *api.Metrics) *api.Gauge {
			return PercentageToGauge(m.HighStartupTimeWithPreRoll)
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"high_startup_time_without_pre_roll",
		func (m *api.Metrics) *api.Gauge {
			return PercentageToGauge(m.HighStartupTimeWithoutPreRoll)
		},
	))
	metricAdders = append(metricAdders, createCountFunc(
		"interval_minutes_played",
		func (m *api.Metrics) *api.Count {
			return m.IntervalMinutesPlayed
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"low_bitrate",
		func (m *api.Metrics) *api.CountPercentage {
			return m.LowBitrate
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"low_bitrate_with_ads",
		func (m *api.Metrics) *api.Gauge {
			return PercentageToGauge(m.LowBitrateWithAds)
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"low_bitrate_without_ads",
		func (m *api.Metrics) *api.Gauge {
			return PercentageToGauge(m.LowBitrateWithoutAds)
		},
	))
	metricAdders = append(metricAdders, createMetricsFunc(addMinutesPlayed))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"non_zero_cirr_ended_plays",
		func (m *api.Metrics) *api.CountPercentage {
			return m.NonZeroCirrEndedPlays
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"percentage_complete",
		func (m *api.Metrics) *api.Gauge {
			return PercentageToGauge(m.PercentageComplete)
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"plays",
		func (m *api.Metrics) *api.CountPercentage {
			return m.Plays
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"rebuffering_ratio",
		func (m *api.Metrics) *api.Gauge {
			return RatioToGauge(m.RebufferingRatio)
		},
	))
	metricAdders = append(metricAdders, createCountFunc(
		"spi_streams",
		func (m *api.Metrics) *api.Count {
			return m.SpiStreams
		},
	))
	metricAdders = append(metricAdders, createCountFunc(
		"spi_unique_devices",
		func (m *api.Metrics) *api.Count {
			return m.SpiUniqueDevices
		},
	))
	metricAdders = append(metricAdders, createCountFunc(
		"spi_unique_viewers",
		func (m *api.Metrics) *api.Count {
			return m.SpiUniqueViewers
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"streaming_performance_index",
		func (m *api.Metrics) *api.Gauge {
			return m.StreamingPerformanceIndex
		},
	))
	metricAdders = append(metricAdders, createCountFunc(
		"unique_devices",
		func (m *api.Metrics) *api.Count {
			return m.UniqueDevices
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"video_playback_failures",
		func (m *api.Metrics) *api.CountPercentage {
			return m.VideoPlaybackFailures
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"video_playback_failures_business",
		func (m *api.Metrics) *api.CountPercentage {
			return m.VideoPlaybackFailuresBusiness
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"video_playback_failures_tech",
		func (m *api.Metrics) *api.CountPercentage {
			return m.VideoPlaybackFailuresTech
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"video_playback_failures_tech_with_ads",
		func (m *api.Metrics) *api.Gauge {
			return PercentageToGauge(m.VideoPlaybackFailuresTechWithAds)
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"video_playback_failures_tech_without_ads",
		func (m *api.Metrics) *api.Gauge {
			return PercentageToGauge(m.VideoPlaybackFailuresTechWithoutAds)
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"video_restart_time",
		func (m *api.Metrics) *api.Gauge {
			return m.VideoRestartTime
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"video_start_failures",
		func (m *api.Metrics) *api.CountPercentage {
			return m.VideoStartFailures
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"video_start_failures_business",
		func (m *api.Metrics) *api.CountPercentage {
			return m.VideoStartFailuresBusiness
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"video_start_failures_tech",
		func (m *api.Metrics) *api.CountPercentage {
			return m.VideoStartFailuresTech
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"video_start_failures_tech_with_pre_roll",
		func (m *api.Metrics) *api.Gauge {
			return PercentageToGauge(m.VideoStartFailuresTechWithPreRoll)
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"video_start_failures_tech_without_pre_roll",
		func (m *api.Metrics) *api.Gauge {
			return PercentageToGauge(m.VideoStartFailuresTechWithoutPreRoll)
		},
	))
	metricAdders = append(metricAdders, createGaugeFunc(
		"video_start_time",
		func (m *api.Metrics) *api.Gauge {
			return m.VideoStartTime
		},
	))
	metricAdders = append(metricAdders, createCountPercentageFunc(
		"zero_cirr_ended_plays",
		func (m *api.Metrics) *api.CountPercentage {
			return m.ZeroCirrEndedPlays
		},
	))
}

func addMetrics(
	entity        *integration.Entity,
	metrics       *api.Metrics,
) {
	ts := time.UnixMilli(metrics.TimeStamp.EpochMs)
	for i := 0; i < len(metricAdders); i += 1 {
		err := metricAdders[i](
			entity,
			ts,
			nil,
			metrics,
		)
		fatalIfErr(err)
	}
}

func addDimensionalMetrics(
	entity        *integration.Entity,
	timestamp     time.Time,
	dimensionData *api.DimensionalData,
) {
	for i := 0; i < len(metricAdders); i += 1 {
		err := metricAdders[i](
			entity,
			timestamp,
			&dimensionData.Dimension,
			&dimensionData.Metrics,
		)
		fatalIfErr(err)
	}
}

func getMetricsData(
	entity *integration.Entity,
	log sdk_log.Logger,
	cfg *Config,
) error {
	log.Debugf("creating a new conviva collector.")

	c, err := api.NewConvivaCollector(
		cfg.ApiV3URL,
		args.ClientId,
		args.ClientSecret,
		cfg.StartOffset,
		cfg.EndOffset,
		cfg.Granularity,
		cfg.RealTime,
		log,
	)
	if err != nil {
		return err
	}

	for _, m := range cfg.Metrics {
		if len(m.Dimensions) == 0 {
			metricData, err := getMetricData(c, log, &m)
			if err != nil {
				return err
			} else if metricData != nil {
				for i := 0; i < len(metricData.TimeSeries); i += 1 {
					addMetrics(entity, &metricData.TimeSeries[i])
				}
			}
			continue
		}

		for _, d := range m.Dimensions {
			metricData, err := getMetricDataByDimension(c, log, &m, d)
			if err != nil {
				return err
			} else if metricData != nil {
				for i := 0; i < len(metricData.TimeSeries); i += 1 {
					dimensions := metricData.TimeSeries[i]
					ts := time.UnixMilli(dimensions.TimeStamp.EpochMs)

					for j := 0; j < len(dimensions.DimensionalData); j += 1 {
						addDimensionalMetrics(
							entity,
							ts,
							&dimensions.DimensionalData[j],
						)
					}
				}
			}
		}
	}

	return nil
}

func getMetricData(
	c *api.ConvivaCollector,
	log sdk_log.Logger,
	m *ConfigMetric,
) (*api.MetricData, error) {
	if m.MetricGroup != "" {
		log.Debugf(
			"collecting conviva metrics for metric group %s...",
			m.MetricGroup,
		)

		return c.CollectMetricGroup(
			m.MetricGroup,
			m.Filters,
			m.StartOffset,
			m.EndOffset,
			m.Granularity,
			m.RealTime,
		)
	} else if m.Metric != "" {
		log.Debugf(
			"collecting conviva metrics for metric %s...",
			m.Metric,
		)

		return c.CollectMetrics(
			[]string {m.Metric},
			m.Filters,
			m.StartOffset,
			m.EndOffset,
			m.Granularity,
			m.RealTime,
		)
	} else if len(m.Names) > 0 {
		log.Debugf(
			"collecting conviva metrics for metrics %v...",
			m.Names,
		)

		return c.CollectMetrics(
			m.Names,
			m.Filters,
			m.StartOffset,
			m.EndOffset,
			m.Granularity,
			m.RealTime,
		)
	}

	return nil, nil
}

func getMetricDataByDimension(
	c *api.ConvivaCollector,
	log sdk_log.Logger,
	m *ConfigMetric,
	d string,
) (*api.DimMetricData, error) {
	if m.MetricGroup != "" {
		log.Debugf(
			"collecting conviva metrics for metric group %s and dimension %s...",
			m.MetricGroup,
			d,
		)

		return c.CollectMetricGroupByDimension(
			m.MetricGroup,
			d,
			m.Filters,
			m.StartOffset,
			m.EndOffset,
			m.Granularity,
			m.RealTime,
		)
	} else if m.Metric != "" {
		log.Debugf(
			"collecting conviva metrics for metric %s and dimension %s...",
			m.Metric,
			d,
		)

		return c.CollectMetricsByDimension(
			[]string {m.Metric},
			d,
			m.Filters,
			m.StartOffset,
			m.EndOffset,
			m.Granularity,
			m.RealTime,
		)
	} else if len(m.Names) > 0 {
		log.Debugf(
			"collecting conviva metrics for metrics %v and dimension %s...",
			m.Names,
			d,
		)

		return c.CollectMetricsByDimension(
			m.Names,
			d,
			m.Filters,
			m.StartOffset,
			m.EndOffset,
			m.Granularity,
			m.RealTime,
		)
	}

	return nil, nil
}
