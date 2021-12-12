package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	//"time"
)

// segment struct, makes templating easier
type brewfather struct {
	props properties
	env   environmentInfo

	Batch
	TemperatureTrendIcon string
	StatusIcon           string
}

const (
	BFUserID  Property = "user_id"
	BFAPIKey  Property = "api_key"
	BFBatchId Property = "batch_id"

	BFDoubleUpIcon      Property = "doubleup_icon"
	BFSingleUpIcon      Property = "singleup_icon"
	BFFortyFiveUpIcon   Property = "fortyfiveup_icon"
	BFFlatIcon          Property = "flat_icon"
	BFFortyFiveDownIcon Property = "fortyfivedown_icon"
	BFSingleDownIcon    Property = "singledown_icon"
	BFDoubleDownIcon    Property = "doubledown_icon"

	BFPlanningStatusIcon     Property = "planning_status_icon"
	BFBrewingStatusIcon      Property = "brewing_status_icon"
	BFFermentingStatusIcon   Property = "fermenting_status_icon"
	BFConditioningStatusIcon Property = "conditioning_status_icon"
	BFCompletedStatusIcon    Property = "completed_status_icon"
	BFArchivedStatusIcon     Property = "archived_status_icon"

	BFCacheTimeout Property = "cache_timeout"
)

// Returned from https://api.brewfather.app/v1/batches/batch_id/readings
type BatchReading struct {
	Comment     string  `json:"comment"`
	Gravity     float64 `json:"sg"`
	DeviceType  string  `json:"type"`
	DeviceId    string  `json:"id"`
	Temperature float64 `json:"temp"`      // celsius
	Timepoint   int64   `json:"timepoint"` // << check what these are...
	Time        int64   `json:"time"`      // <<
}

// Returned from https://api.brewfather.app/v1/batches/batch_id
type Batch struct {
	Status string `json:"status"`
	Recipe struct {
		Name string `json:"name"`
	} `json:"recipe"`

	Reading          *BatchReading
	TemperatureTrend float64 // diff between this and last, short term trend
}

func (bf *brewfather) enabled() bool {
	data, err := bf.getResult()
	if err != nil {
		return false
	}
	bf.Batch = *data
	bf.TemperatureTrendIcon = bf.getTrendIcon(bf.TemperatureTrend)
	bf.StatusIcon = bf.getBatchStatusIcon(data.Status)

	return true
}

func (bf *brewfather) getTrendIcon(trend float64) string {
	if trend > 0 {
		if trend > 4 {
			return bf.props.getString(BFDoubleUpIcon, "↑↑")
		}

		if trend > 2 {
			return bf.props.getString(BFSingleUpIcon, "↑")
		}

		if trend > 0.5 {
			return bf.props.getString(BFFortyFiveUpIcon, "↗")
		}
	}

	if trend < -4 {
		return bf.props.getString(BFDoubleDownIcon, "↓↓")
	}

	if trend < -2 {
		return bf.props.getString(BFSingleDownIcon, "↓")
	}

	if trend < -0.5 {
		return bf.props.getString(BFFortyFiveDownIcon, "↘")
	}

	return bf.props.getString(BFFlatIcon, "→")
}

func (bf *brewfather) getBatchStatusIcon(batchStatus string) string {
	switch batchStatus {
	case "Planning":
		return bf.props.getString(BFPlanningStatusIcon, "") //辰?
	case "Brewing":
		return bf.props.getString(BFBrewingStatusIcon, "")
	case "Fermenting": //ﭙ
		return bf.props.getString(BFFermentingStatusIcon, "")
	case "Conditioning":
		return bf.props.getString(BFConditioningStatusIcon, "צּ") //?
	case "Completed":
		return bf.props.getString(BFCompletedStatusIcon, "祖")
	case "Archived":
		return bf.props.getString(BFArchivedStatusIcon, "")
	default:
		return ""
	}
}

func (bf *brewfather) string() string {
	segmentTemplate := bf.props.getString(SegmentTemplate, "{{.StatusIcon}} {{.Recipe.Name}}{{ if and (.Reading) (eq .Status \"Fermenting\")}}: {{.Reading.Gravity}} {{.Reading.Temperature}}\ue33e {{.TemperatureTrendIcon}} {{end}}")
	template := &textTemplate{
		Template: segmentTemplate,
		Context:  bf,
		Env:      bf.env,
	}
	text, err := template.render()
	if err != nil {
		return err.Error()
	}

	return text
}

func (bf *brewfather) getResult() (*Batch, error) {

	getFromCache := func(key string) (*Batch, error) {
		val, found := bf.env.cache().get(key)
		// we got something from the cache
		if found {
			var result Batch
			err := json.Unmarshal([]byte(val), &result)
			if err == nil {
				return &result, nil
			}
		}
		return nil, errors.New("no data in cache")
	}

	putToCache := func(key string, batch *Batch, cacheTimeout int) error {
		cacheJson, err := json.Marshal(batch)
		if err != nil {
			return err
		}

		bf.env.cache().set(key, string(cacheJson), cacheTimeout)

		return nil
	}

	userId := bf.props.getString(BFUserID, "")
	if len(userId) == 0 {
		return nil, errors.New("missing Brewfather user id (user_id)")
	}

	apiKey := bf.props.getString(BFAPIKey, "")
	if len(apiKey) == 0 {
		return nil, errors.New("missing Brewfather api key (api_key)")
	}

	batchId := bf.props.getString(BFBatchId, "")
	if len(batchId) == 0 {
		return nil, errors.New("missing Brewfather batch id (batch_id)")
	}

	authString := fmt.Sprintf("%s:%s", userId, apiKey)
	authStringb64 := base64.StdEncoding.EncodeToString([]byte(authString))
	authHeader := fmt.Sprintf("Basic %s", authStringb64)
	batchUrl := fmt.Sprintf("https://api.brewfather.app/v1/batches/%s", batchId)
	batchReadingsUrl := fmt.Sprintf("https://api.brewfather.app/v1/batches/%s/readings", batchId)

	httpTimeout := bf.props.getInt(HTTPTimeout, DefaultHTTPTimeout)
	cacheTimeout := bf.props.getInt(BFCacheTimeout, 5)

	if cacheTimeout > 0 {
		if data, err := getFromCache(batchUrl); err == nil {
			return data, nil
		}
	}

	// batch
	body, err := bf.env.doGetWithAuth(batchUrl, httpTimeout, authHeader)
	if err != nil {
		return nil, err
	}
	var batch Batch
	err = json.Unmarshal(body, &batch)
	if err != nil {
		return nil, err
	}

	// readings
	body, err = bf.env.doGetWithAuth(batchReadingsUrl, httpTimeout, authHeader)
	if err != nil {
		return nil, err
	}

	var arr []*BatchReading
	err = json.Unmarshal(body, &arr)
	if err != nil {
		return nil, err
	}

	if len(arr) > 0 {
		// could just take latest reading using their API, but that won't allow us to see trend - get 'em all and sort by time instead
		sort.Slice(arr[:], func(i, j int) bool {
			return arr[i].Time > arr[j].Time
		})

		// Keep the latest one
		batch.Reading = arr[0]

		if len(arr) > 1 {
			batch.TemperatureTrend = arr[0].Temperature - arr[1].Temperature
		}
	}

	if cacheTimeout > 0 {
		putToCache(batchUrl, &batch, cacheTimeout)
	}

	return &batch, nil
}

func (bf *brewfather) init(props properties, env environmentInfo) {
	bf.props = props
	bf.env = env
}
