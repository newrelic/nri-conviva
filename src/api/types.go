package api

/* Metrics */

type Count struct {
	Value int64 `json:"count"`
}

type Gauge struct {
	Value float64 `json:"value"`
}

type Percentage struct {
	Value float64 `json:"percentage"`
}

type CountPercentage struct {
	Count
	Percentage
}

type Ratio struct {
	Ratio float64 `json:"ratio"`
}

type Bitrate struct {
	Bps float64 `json:"bps"`
}

type EndedPlays struct {
	Count
	PerUniqueDevice float64 `json:"per_unique_device"`
}

type MinutesPlayed struct {
	EndedPlays
	PerEndedPlay float64 `json:"per_ended_play"`
}

type Framerate struct {
	Fps float64 `json:"fps"`
}

/* Metrics Map */

type Metrics struct {
	TimeStamp       TimeStamp         `json:"timestamp"`
	Abandonment							    *CountPercentage  `json:"abandonment"`
	AbandonmentWithPreRoll                  *Percentage       `json:"abandonment_with_pre_roll"`
	AbandonmentWithoutPreRoll               *Percentage       `json:"abandonment_without_pre_roll"`
	AdActualDuration			            *Gauge            `json:"ad_actual_duration"`
	AdAttempts                              *Count            `json:"ad_attempts"`
	AdBitrate                               *Bitrate          `json:"ad_bitrate"`
    AdCompletedCreativePlays                *Gauge            `json:"ad_completed_creative_plays"`
	AdConcurrentPlays                       *Count            `json:"ad_concurrent_plays"`
	AdConnectionInducedRebufferingRatio     *Ratio            `json:"ad_connection_induced_rebuffering_ratio"`
	AdEndedPlays                            *EndedPlays       `json:"ad_ended_plays"`
	AdExitBeforeVideoStarts                 *CountPercentage  `json:"exits_before_ad_start"`
	AdFramerate                             *Framerate        `json:"ad_framerate"`
    AdMinutesPlayed                         *MinutesPlayed    `json:"ad_minutes_played"`
	AdPercentageComplete                    *Percentage       `json:"ad_percentage_complete"`
	AdPlays                                 *CountPercentage  `json:"ad_plays"`
	AdRebufferingRatio                      *Ratio            `json:"ad_rebuffering_ratio"`
	AdUniqueDevices                         *Count            `json:"ad_unique_devices"`
	AdVideoPlaybackFailures                 *CountPercentage  `json:"ad_video_playback_failures"`
    AdVideoRestartTime                      *Gauge            `json:"ad_video_restart_time"`
	AdVideoStartFailures                    *CountPercentage  `json:"ad_video_start_failures"`
	AdVideoStartTime                        *Gauge            `json:"ad_video_start_time"`
	Attempts                                *Count            `json:"attempts"`
	AttemptsWithPreRoll                     *CountPercentage  `json:"attempts_with_pre_roll"`
	AttemptsWithoutPreRoll                  *CountPercentage  `json:"attempts_without_pre_roll"`
    BadSession                              *CountPercentage  `json:"bad_session"`
    BadSessionAverageLifePlayingTimeMins    *Gauge            `json:"bad_session_average_life_playing_time_mins"`
    BadUniqueDevices                        *CountPercentage  `json:"bad_unique_devices"`
    BadUniqueViewers                        *CountPercentage  `json:"bad_unique_viewers"`
	Bitrate                                 *Bitrate          `json:"bitrate"`
	ConcurrentPlays                         *Count            `json:"concurrent_plays"`
	ConnectionInducedRebufferingRatio       *Ratio            `json:"connection_induced_rebuffering_ratio"`
	EndedPlays                              *EndedPlays       `json:"ended_plays"`
	EndedPlaysWithAds                       *CountPercentage  `json:"ended_plays_with_ads"`
	EndedPlaysWithoutAds                    *CountPercentage  `json:"ended_plays_without_ads"`
	ExitBeforeVideoStarts                   *CountPercentage  `json:"exit_before_video_starts"`
	Framerate                               *Framerate        `json:"framerate"`
    GoodSession                             *Count            `json:"good_session"`
    GoodSessionAverageLifePlayingTimeMins   *Gauge            `json:"good_session_average_life_playing_time_mins"`
    GoodUniqueDevices                       *Count            `json:"good_unique_devices"`
    GoodUniqueViewers                       *Count            `json:"good_unique_viewers"`
	HighRebuffering                         *CountPercentage  `json:"high_rebuffering"`
	HighRebufferingWithAds                  *Percentage       `json:"high_rebuffering_with_ads"`
	HighRebufferingWithoutAds               *Percentage       `json:"high_rebuffering_without_ads"`
	HighStartupTime                         *CountPercentage  `json:"high_startup_time"`
	HighStartupTimeWithPreRoll              *Percentage       `json:"high_startup_time_with_pre_roll"`
	HighStartupTimeWithoutPreRoll           *Percentage       `json:"high_startup_time_without_pre_roll"`
	IntervalMinutesPlayed                   *Count            `json:"interval_minutes_played"`
	LowBitrate                              *CountPercentage  `json:"low_bitrate"`
	LowBitrateWithAds                       *Percentage       `json:"low_bitrate_with_ads"`
	LowBitrateWithoutAds                    *Percentage       `json:"low_bitrate_without_ads"`
    MinutesPlayed                           *MinutesPlayed    `json:"minutes_played"`
    NonZeroCirrEndedPlays                   *CountPercentage  `json:"non_zero_cirr_ended_plays"`
    PercentageComplete                      *Percentage       `json:"percentage_complete"`
	Plays                                   *CountPercentage  `json:"plays"`
	RebufferingRatio                        *Ratio            `json:"rebuffering_ratio"`
    SpiStreams                              *Count            `json:"spi_streams"`
    SpiUniqueDevices                        *Count            `json:"spi_unique_devices"`
    SpiUniqueViewers                        *Count            `json:"spi_unique_viewers"`
    StreamingPerformanceIndex               *Gauge            `json:"streaming_performance_index"`
    UniqueDevices                           *Count            `json:"unique_devices"`
	VideoPlaybackFailures                   *CountPercentage  `json:"video_playback_failures"`
	VideoPlaybackFailuresBusiness           *CountPercentage  `json:"video_playback_failures_business"`
	VideoPlaybackFailuresTech               *CountPercentage  `json:"video_playback_failures_tech"`
	VideoPlaybackFailuresTechWithAds        *Percentage       `json:"video_playback_failures_tech_with_ads"`
	VideoPlaybackFailuresTechWithoutAds     *Percentage       `json:"video_playback_failures_tech_without_ads"`
	VideoRestartTime                        *Gauge            `json:"video_restart_time"`
	VideoStartFailures                      *CountPercentage  `json:"video_start_failures"`
	VideoStartFailuresBusiness              *CountPercentage  `json:"video_start_failures_business"`
	VideoStartFailuresTech                  *CountPercentage  `json:"video_start_failures_tech"`
	VideoStartFailuresTechWithPreRoll       *Percentage       `json:"video_start_failures_tech_with_pre_roll"`
	VideoStartFailuresTechWithoutPreRoll    *Percentage       `json:"video_start_failures_tech_without_pre_roll"`
	VideoStartTime                          *Gauge            `json:"video_start_time"`
    ZeroCirrEndedPlays                      *CountPercentage  `json:"zero_cirr_ended_plays"`
}

/* Dimension Key:Value */

type Dimension struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

/* Dimensional Data Struct */

type DimensionalData struct {
	Dimension Dimension `json:"dimension"`
	Metrics   Metrics   `json:"metrics"`
}

/* Timestamp */

type TimeStamp struct {
	EpochMs int64  `json:"epoch_ms"`
	IsoDate string `json:"iso_date"`
}

/* Datapoint64 in a TimeSeries */

type Dimensions struct {
	TimeStamp       TimeStamp         `json:"timestamp"`
	DimensionalData []DimensionalData `json:"dimensional_data"`
}

/* Totals */

type Total struct {
	Metrics
}

/* Metric Data Response */

type MetricData struct {
	TimeSeries []Metrics `json:"time_series"`
	Total      Total       `json:"total"`
}

type DimMetricData struct {
	TimeSeries []Dimensions `json:"time_series"`
	Total      Total       `json:"total"`
}

/* Generic Logger Interface */
type Logger interface {
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}