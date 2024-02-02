package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	FIFTEEN_MINUTES = 15 * time.Minute
)

type ConvivaCollector struct {
	URL             string
	ClientId        string
	ClientSecret    string
	StartOffset     time.Duration
	EndOffset		time.Duration
	Granularity     string
	RealTime		*bool
	log             Logger
}

func NewConvivaCollector(
	URL             string,
	ClientId        string,
	ClientSecret    string,
	StartOffset     string,
	EndOffset		string,
	Granularity     string,
	RealTime		*bool,
	log             Logger,
) (*ConvivaCollector, error) {
	var err error

	startOffset := time.Duration(0)
	endOffset := time.Duration(0)

	if StartOffset != "" {
		if startOffset, err = time.ParseDuration(StartOffset); err != nil {
			return nil, err
		}
	}

	if EndOffset != "" {
		if endOffset, err = time.ParseDuration(EndOffset); err != nil {
			return nil, err
		}
	}

	return &ConvivaCollector{
		URL,
		ClientId,
		ClientSecret,
		startOffset,
        endOffset,
		Granularity,
		RealTime,
		log,
	}, nil
}

func (c *ConvivaCollector) CollectMetricsByDimension(
	metricNames []string,
	dimension string,
	filters map[string][]string,
	startOffset string,
	endOffset string,
	granularity string,
	realTime *bool,
) (*DimMetricData, error) {
	url, err := c.makeUrl(
		c.makePath(metricNames, "", dimension),
		metricNames,
		filters,
		startOffset,
		endOffset,
		granularity,
		realTime,
	)
	if err != nil {
		return nil, err
	}

	return c.getMetricDataByDimension(url)
}

func (c *ConvivaCollector) CollectMetricGroupByDimension(
	metricGroup string,
	dimension string,
	filters map[string][]string,
	startOffset string,
	endOffset string,
	granularity string,
	realTime *bool,
) (*DimMetricData, error) {
	url, err := c.makeUrl(
		c.makePath(nil, metricGroup, dimension),
		nil,
		filters,
		startOffset,
		endOffset,
		granularity,
		realTime,
	)
	if err != nil {
		return nil, err
	}

	return c.getMetricDataByDimension(url)
}

func (c *ConvivaCollector) CollectMetrics(
	metricNames []string,
	filters map[string][]string,
	startOffset string,
	endOffset string,
	granularity string,
	realTime *bool,
) (*MetricData, error) {
	url, err := c.makeUrl(
		c.makePath(metricNames, "", ""),
		metricNames,
		filters,
		startOffset,
		endOffset,
		granularity,
		realTime,
	)
	if err != nil {
		return nil, err
	}

	return c.getMetricData(url)
}

func (c *ConvivaCollector) CollectMetricGroup(
	metricGroup string,
	filters map[string][]string,
	startOffset string,
	endOffset string,
	granularity string,
	realTime *bool,
) (*MetricData, error) {
	url, err := c.makeUrl(
		c.makePath(nil, metricGroup, ""),
		nil,
		filters,
		startOffset,
		endOffset,
		granularity,
		realTime,
	)
	if err != nil {
		return nil, err
	}

	return c.getMetricData(url)
}

func (c ConvivaCollector) makePath(
	metricNames []string,
	metricGroup string,
	dimension string,
) string {
	var (
		s string
		l = len(metricNames)
	)

	if metricGroup != "" {
		s = metricGroup
	} else if l > 1 {
		s = "custom-selection"
	} else if l == 1 {
		s = metricNames[0]
	}

	if dimension != "" {
		s += fmt.Sprintf("/group-by/%s", dimension)
	}

	return s
}

func useRealTime(
	start time.Duration,
	r1 *bool,
	r2 *bool,
) bool {
	if r1 != nil && !*r1 {
		return false
	}

	if r2 != nil && !*r2 {
		return false
	}

	if start != 0 && start > FIFTEEN_MINUTES {
		return false
	}

	return true
}

func (c ConvivaCollector) makeUrl(
	path string,
	metricNames []string,
	filters map[string][]string,
	startOffset string,
	endOffset string,
	granularity string,
	realTime *bool,
) (string, error) {
	var params []string

	start, err := getDuration(startOffset, c.StartOffset)
	if err != nil {
		return "", err
	}

	end, err := getDuration(endOffset, c.EndOffset)
	if err != nil {
		return "", err
	}

	if (start != 0) {
		c.log.Debugf("start: %d, end: %d", start, end)

		if end > start {
			return "", fmt.Errorf(
				"end offset %d is more than start offset %d",
				end,
				start,
			)
		}

		params = addTimeRange(params, start, end)
	}

	params = addGranularity(params, granularity, c.Granularity)

	endpoint := "real-time-metrics"
	if !useRealTime(start, realTime, c.RealTime) {
		endpoint = "metrics"
	}

	if len(filters) > 0 {
		for k, v := range filters {
			for _, u := range v {
				params = append(params, fmt.Sprintf("%s=%s", k, u))
			}
		}
	}

	if len(metricNames) > 1 {
		for _, u := range metricNames {
			params = append(params, "metric=" + u)
		}
	}

	if len(params) == 0 {
		return fmt.Sprintf(
			"%s/%s/%s",
			c.URL,
			endpoint,
			path,
		), nil
	}

	return fmt.Sprintf(
		"%s/%s/%s?%s",
		c.URL,
		endpoint,
		path,
		strings.Join(params, "&"),
	), nil
}

func getDuration(offset1 string, offset2 time.Duration) (time.Duration, error) {
	d := offset2

	if offset1 != "" {
		if dur, err := time.ParseDuration(offset1); err != nil {
			return 0, err
		} else {
			d = dur
		}
	}

	return d, nil
}

func addTimeRange(params []string, start, end time.Duration) []string {
	params = append(params, fmt.Sprintf(
		"start_epoch=%d",
		time.Now().Add(-start).UnixMilli() / 1000,
	))

	params = append(params, fmt.Sprintf(
		"end_epoch=%d",
		time.Now().Add(-end).UnixMilli() / 1000,
	))

	return params
}

func addGranularity(params []string, g1, g2 string) []string {
	if g1 != "" {
		return append(params, "granularity=" + g1)
	} else if g2 != "" {
		return append(params, "granularity=" + g2)
	}

	return params
}

func (c ConvivaCollector) makeRequest(url string) ([]byte, error) {
	then := time.Now()
	client := &http.Client{}

	c.log.Debugf("making metrics request using URL %s...", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.ClientId, c.ClientSecret)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	c.log.Debugf(
		"request for %s took %dms",
		url,
		time.Since(then).Milliseconds(),
	)

	return body, nil
}

func (c ConvivaCollector) getMetricData(url string) (*MetricData, error) {
	body, err := c.makeRequest(url)
	if err != nil {
		return nil, err
	}

	then := time.Now()
	metricData := &MetricData{}

	c.log.Debugf("unmarshalling data...")

	err = json.Unmarshal(body, metricData)
	if err != nil {
		return nil, err
	}

	c.log.Debugf(
		"unmarshalling took %dms",
		time.Since(then).Milliseconds(),
	)

	return metricData, nil
}

func (c ConvivaCollector) getMetricDataByDimension(
	url string,
)(*DimMetricData, error) {
	body, err := c.makeRequest(url)
	if err != nil {
		return nil, err
	}

	then := time.Now()
	metricData := &DimMetricData{}

	c.log.Debugf("unmarshalling data...")

	err = json.Unmarshal(body, metricData)
	if err != nil {
		return nil, err
	}

	c.log.Debugf(
		"unmarshalling took %dms",
		time.Since(then).Milliseconds(),
	)

	return metricData, nil
}
