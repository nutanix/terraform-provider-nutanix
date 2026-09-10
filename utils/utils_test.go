package utils

import (
	"testing"
	"time"
)

type fakeWaitTimeoutClient struct {
	wt int64
}

func (f *fakeWaitTimeoutClient) GetWaitTimeout() int64 { return f.wt }

func TestResolveWaitTimeout(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		meta           interface{}
		defaultTimeout time.Duration
		expected       time.Duration
	}{
		{
			name:           "uses provider wait_timeout when positive",
			meta:           &fakeWaitTimeoutClient{wt: 90},
			defaultTimeout: 30 * time.Minute,
			expected:       90 * time.Minute,
		},
		{
			name:           "falls back to default timeout when wait_timeout is zero",
			meta:           &fakeWaitTimeoutClient{wt: 0},
			defaultTimeout: 45 * time.Minute,
			expected:       45 * time.Minute,
		},
		{
			name:           "falls back to default timeout when wait_timeout is negative",
			meta:           &fakeWaitTimeoutClient{wt: -5},
			defaultTimeout: 45 * time.Minute,
			expected:       45 * time.Minute,
		},
		{
			name:           "wait_timeout of 1 minute is honoured",
			meta:           &fakeWaitTimeoutClient{wt: 1},
			defaultTimeout: 60 * time.Minute,
			expected:       1 * time.Minute,
		},
		{
			name:           "meta without wait_timeout accessor falls back to default",
			meta:           struct{}{},
			defaultTimeout: 20 * time.Minute,
			expected:       20 * time.Minute,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := ResolveWaitTimeout(tc.meta, tc.defaultTimeout)
			if actual != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
