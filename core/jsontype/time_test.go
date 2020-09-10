package jsontype

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMilliTime_GetBSON(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
	}{
		{
			name: "now",
			tm:   time.Now(),
		},
		{
			name: "future",
			tm:   time.Now().Add(time.Hour),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := MilliTime{test.tm}.GetBSON()
			assert.Nil(t, err)
			assert.Equal(t, test.tm, got)
		})
	}
}

func TestMilliTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
	}{
		{
			name: "now",
			tm:   time.Now(),
		},
		{
			name: "future",
			tm:   time.Now().Add(time.Hour),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := MilliTime{test.tm}.MarshalJSON()
			assert.Nil(t, err)
			assert.Equal(t, strconv.FormatInt(test.tm.UnixNano()/1e6, 10), string(b))
		})
	}
}

func TestMilliTime_Milli(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
	}{
		{
			name: "now",
			tm:   time.Now(),
		},
		{
			name: "future",
			tm:   time.Now().Add(time.Hour),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			n := MilliTime{test.tm}.Milli()
			assert.Equal(t, test.tm.UnixNano()/1e6, n)
		})
	}
}

func TestMilliTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
	}{
		{
			name: "now",
			tm:   time.Now(),
		},
		{
			name: "future",
			tm:   time.Now().Add(time.Hour),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mt MilliTime
			s := strconv.FormatInt(test.tm.UnixNano()/1e6, 10)
			err := mt.UnmarshalJSON([]byte(s))
			assert.Nil(t, err)
			s1, err := mt.MarshalJSON()
			assert.Nil(t, err)
			assert.Equal(t, s, string(s1))
		})
	}
}
