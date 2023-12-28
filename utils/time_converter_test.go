package utils

import (
	"testing"
)

type TimeTest struct {
	input         interface{}
	expectedMs    int
	expectedTime  string
	expectedTicks int
}

func TestTimeConversions(t *testing.T) {
	testCases := []TimeTest{
		{"23:43.400", 1423400, "23:43.80", 2373},
		{"01:23:56", 5036000, "01:23:56.40", 8394},
		{"01:23.3", 83300, "01:23.40", 139},
		{"01:03:58", 3838000, "01:03:58.20", 6397},
		{"02:09.00", 129000, "02:09.00", 215},
		{"06:17.40", 377400, "06:17.40", 629},
		{"14:32.40", 872400, "14:32.40", 1454},
		{"00:34.80", 34800, "34.80", 58},
	}

	for _, tc := range testCases {
		testTimeToMs(t, tc)
		testTimeToTicks(t, tc)
		testTicksToTime(t, tc)
	}
}

func testTimeToMs(t *testing.T, tc TimeTest) {
	result := TimeToMs(tc.input.(string))
	if result != tc.expectedMs {
		t.Errorf("TimeToMs(%v): expected %v, got %v", tc.input, tc.expectedMs, result)
	}
}

func testTimeToTicks(t *testing.T, tc TimeTest) {
	result := TimeToTicks(tc.input.(string))
	if result != tc.expectedTicks {
		t.Errorf("TimeToTicks(%v): expected %v, got %v", tc.input, tc.expectedTicks, result)
	}
}

func testTicksToTime(t *testing.T, tc TimeTest) {
	result := TicksToTime(tc.expectedTicks)
	if result != tc.expectedTime {
		t.Errorf("TicksToTime(%v): expected %v, got %v", tc.expectedTicks, tc.expectedTime, result)
	}
}
