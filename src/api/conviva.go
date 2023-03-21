package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ConvivaCollector struct {
	URL             string
	ClientId        string
	ClientSecret    string
	StartOffset     time.Duration
	EndOffset		time.Duration
	Granularity     string
	log             Logger
}

func NewConvivaCollector(
	URL             string,
	ClientId        string,
	ClientSecret    string,
	StartOffset     string,
	EndOffset		string,
	Granularity     string,
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
) (*DimMetricData, error) {
	url, err := c.makeUrl(
		c.makePath(metricNames, "", dimension),
		metricNames,
		filters,
		startOffset,
		endOffset,
		granularity,
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
) (*DimMetricData, error) {
	url, err := c.makeUrl(
		c.makePath(nil, metricGroup, dimension),
		nil,
		filters,
		startOffset,
		endOffset,
		granularity,
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
) (*MetricData, error) {
	url, err := c.makeUrl(
		c.makePath(metricNames, "", ""),
		metricNames,
		filters,
		startOffset,
		endOffset,
		granularity,
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
) (*MetricData, error) {
	url, err := c.makeUrl(
		c.makePath(nil, metricGroup, ""),
		nil,
		filters,
		startOffset,
		endOffset,
		granularity,
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
func (c ConvivaCollector) makeUrl(
	path string,
	metricNames []string,
	filters map[string][]string,
	startOffset string,
	endOffset string,
	granularity string,
) (string, error) {
	var params []string

	params, err := addTimeParam(
		params,
		"start_epoch",
		startOffset,
		c.StartOffset,
	)
	if err != nil {
		return "", err
	}

	params, err = addTimeParam(
		params,
		"end_epoch",
		endOffset,
		c.EndOffset,
	)
	if err != nil {
		return "", err
	}

	params = addGranularity(params, granularity, c.Granularity)

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

	return fmt.Sprintf(
		"%s/%s?%s",
		c.URL,
		path,
		strings.Join(params, "&"),
	), nil
}

func addTimeParam(
	params []string,
	paramName, offset1 string,
	offset2 time.Duration,
) ([]string, error) {
	if offset1 == "" && offset2 == 0 {
		return params, nil
	}
	
	d := offset2

	if offset1 != "" {
		if dur, err := time.ParseDuration(offset1); err != nil {
			return nil, err
		} else {
			d = dur
		}
	}

	return append(
		params,
		fmt.Sprintf(
			"%s=%d",
			paramName,
			time.Now().Add(-d).UnixMilli() / 1000,
		),
	), nil
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