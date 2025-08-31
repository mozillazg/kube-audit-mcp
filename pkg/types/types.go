package types

import (
	"strconv"
	"strings"
	"time"
)

type QueryAuditLogParams struct {
	ClusterName   string    `json:"cluster_name"`
	StartTime     TimeParam `json:"start_time"`
	EndTime       TimeParam `json:"end_time"`
	User          string    `json:"user"`
	Namespace     string    `json:"namespace"`
	Verbs         []string  `json:"verbs"`
	ResourceTypes []string  `json:"resource_types"`
	ResourceName  string    `json:"resource_name"`
	Limit         int       `json:"limit"`
}

type TimeParam struct {
	time.Time
	rawInput []byte
}

func NewTimeParam(t time.Time) TimeParam {
	return TimeParam{Time: t.UTC()}
}

func (t *TimeParam) MarshalJSON() ([]byte, error) {
	if len(t.rawInput) > 0 {
		return t.rawInput, nil
	}

	return []byte(time.Since(t.Time).String()), nil
}

func (t *TimeParam) UnmarshalJSON(bytes []byte) error {
	var err error
	s := string(bytes)
	s = strings.TrimSpace(s)
	s = strings.Trim(s, `"'`)
	s = strings.TrimSpace(s)
	t.rawInput = bytes

	// Try to parse as RFC3339 first
	t.Time, err = time.Parse(time.RFC3339, s)
	if err == nil {
		return nil
	}

	var dt time.Duration
	switch {
	case strings.HasSuffix(s, "w"):
		weeks := strings.TrimSuffix(s, "w")
		d, err := strconv.Atoi(weeks)
		if err != nil {
			return err
		}
		dt = time.Duration(d) * 7 * 24 * time.Hour
		break
	case strings.HasSuffix(s, "d"):
		days := strings.TrimSuffix(s, "d")
		d, err := strconv.Atoi(days)
		if err != nil {
			return err
		}
		dt = time.Duration(d) * 24 * time.Hour
		break
	default:
		dt, err = time.ParseDuration(s)
		if err != nil {
			return err
		}
	}
	t.Time = time.Now().Add(-dt)

	return nil
}
