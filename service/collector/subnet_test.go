package collector

import (
	"strconv"
	"testing"
)

func TestAvailableIPPercentage(t *testing.T) {
	testCases := []struct {
		name         string
		cidr         string
		availableIPs int64

		expectedPercentage float64
		expectedError      bool
	}{
		{
			name:         "case 0",
			cidr:         "10.1.0.0/27",
			availableIPs: 27,

			expectedPercentage: 1,
			expectedError:      false,
		},
		{
			name:         "case 1",
			cidr:         "10.1.0.0/21",
			availableIPs: 1000,

			expectedPercentage: 0.48947626040137054,
			expectedError:      false,
		},
		{
			name:         "case 2",
			cidr:         "10.1.0.0/27/12",
			availableIPs: 25,

			expectedPercentage: 0,
			expectedError:      true,
		},
		{
			name:         "case 3",
			cidr:         "10.1.0.0/2.7",
			availableIPs: 25,

			expectedPercentage: 0,
			expectedError:      true,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {

			percentage, err := getAvailableIPPercentage(tc.cidr, tc.availableIPs)

			if percentage != tc.expectedPercentage {
				t.Fatalf("expected %v, got %v", tc.expectedPercentage, percentage)
			}
			if (err == nil) == tc.expectedError {
				t.Fatalf("expected error response to be %v, got %v", tc.expectedError, err)
			}
		})
	}
}
