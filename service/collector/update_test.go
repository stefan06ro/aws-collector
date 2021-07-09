package collector

import (
	"math"
	"strconv"
	"testing"
)

func TestCalculateUpdateMetrics(t *testing.T) {
	testCases := []struct {
		name string

		batch string
		pause string
		min   int
		max   int

		pauseTime       float64
		batchNumber     float64
		batchPercentage float64
		expectedError   bool
	}{
		{
			// default case
			name: "case 0",

			batch: DefaultBatchSize,
			pause: DefaultPauseTime,
			min:   2,
			max:   10,

			pauseTime:       900,
			batchNumber:     3,
			batchPercentage: 0.3,
			expectedError:   false,
		},
		{
			// static batch case
			name: "case 1",

			batch: "5",
			pause: DefaultPauseTime,
			min:   2,
			max:   10,

			pauseTime:       900,
			batchNumber:     5,
			batchPercentage: 2.5,
			expectedError:   false,
		},
		{
			// min 0 case
			name: "case 2",

			batch: "5",
			pause: DefaultPauseTime,
			min:   0,
			max:   10,

			pauseTime:       900,
			batchNumber:     5,
			batchPercentage: 5,
			expectedError:   false,
		},
		{
			// max 0 case
			name: "case 3",

			batch: "0.5",
			pause: DefaultPauseTime,
			min:   0,
			max:   0,

			pauseTime:       900,
			batchNumber:     0,
			batchPercentage: 0.5,
			expectedError:   false,
		},
		{
			// pause time in seconds
			name: "case 4",

			batch: DefaultBatchSize,
			pause: "PT15S",
			min:   2,
			max:   10,

			pauseTime:       15,
			batchNumber:     3,
			batchPercentage: 0.3,
			expectedError:   false,
		},
		{
			// pause time with minutes and seconds
			name: "case 5",

			batch: DefaultBatchSize,
			pause: "PT1M30S",
			min:   2,
			max:   10,

			pauseTime:       90,
			batchNumber:     3,
			batchPercentage: 0.3,
			expectedError:   false,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {

			pauseTime, batchNumber, batchPercentage, err := calculateUpdateMetrics(tc.batch, tc.pause, tc.min, tc.max)

			if batchNumber != tc.batchNumber {
				t.Fatalf("expected batch number %v, got %v", tc.batchNumber, batchNumber)
			}
			if batchPercentage != tc.batchPercentage {
				t.Fatalf("expected batch percentage %v, got %v", tc.batchPercentage, batchPercentage)
			}
			if math.Ceil(pauseTime) != tc.pauseTime {
				t.Fatalf("expected pause time %v, got %v", tc.pauseTime, pauseTime)
			}
			if (err == nil) == tc.expectedError {
				t.Fatalf("expected error response to be %v, got %v", tc.expectedError, err)
			}
		})
	}
}
