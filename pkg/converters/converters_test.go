package converters

import "testing"

func TestHumanReadable(t *testing.T) {
	var converterTests = []struct{
		in string
		expected uint64
	} {
		{"5k", 5000},
		{"500", 500},
		{"5m", 5000000},
		{"5b", 5000000000},
	}

	for _, tt := range converterTests {
		actual, err := HumanReadable(tt.in)
		if err != nil {
			t.Errorf("HumanReadable(%s) failed: %s", tt.in, err)
		} else {
			if actual != tt.expected {
				t.Errorf("HumanReadable(%s): expected %d, actual %d", tt.in, tt.expected, actual)
			}
		}
	}
}

func TestHumanReadableInvalid(t *testing.T) {
	value := "55b7"
	actual, err := HumanReadable(value)
	if err == nil {
		t.Errorf("HumanReadable(%s): expected error, actual %d", value, actual)
	}

	value = "m"
	actual, err = HumanReadable(value)
	if err == nil {
		t.Errorf("HumanReadable(%s): expected error, actual %d", value, actual)
	}

	value = "10k4b"
	actual, err = HumanReadable(value)
	if err == nil {
		t.Errorf("HumanReadable(%s): expected error, actual %d", value, actual)
	}
}
