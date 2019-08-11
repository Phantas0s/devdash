package internal

import "testing"

func Test_typeID(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
		widget   Widget
	}{
		{
			name:     "happy case",
			expected: "bar",
			widget: Widget{
				Name: "ga.bar_chart",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.widget.typeID()

			if actual != tc.expected {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}

func Test_serviceID(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
		widget   Widget
	}{
		{
			name:     "happy case",
			expected: "ga",
			widget: Widget{
				Name: "ga.bar_chart",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.widget.serviceID()

			if actual != tc.expected {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}
