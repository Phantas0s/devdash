package internal

import (
	"reflect"
	"testing"
)

func Test_addDefaultTheme(t *testing.T) {
	testCases := []struct {
		name     string
		expected Widget
		widget   Widget
		project  project
	}{
		{
			name: "project with named theme / widget without options",
			expected: Widget{Name: "ga.bar_thingy", Theme: "super_theme",
				Options: map[string]string{
					"border_color": "blue",
					"title_color":  "green",
				},
			},
			widget: Widget{
				Name:  "ga.bar_thingy",
				Theme: "super_theme",
			},
			project: project{
				themes: map[string]map[string]string{
					"super_theme": {
						"border_color": "blue",
						"title_color":  "green",
					},
				},
			},
		},
		{
			name: "project with named theme / widget with non overwriting options",
			expected: Widget{Name: "ga.bar_thingy", Theme: "super_theme",
				Options: map[string]string{
					"border_color": "blue",
					"title_color":  "green",
					"text_color":   "red",
				},
			},
			widget: Widget{Name: "ga.bar_thingy", Theme: "super_theme",
				Options: map[string]string{"text_color": "red"},
			},
			project: project{
				themes: map[string]map[string]string{
					"super_theme": {
						"border_color": "blue",
						"title_color":  "green",
					},
				},
			},
		},
		{
			name: "project with named theme / widget with overwriting options",
			expected: Widget{Name: "ga.bar_thingy", Theme: "super_theme",
				Options: map[string]string{
					"border_color": "blue",
					"title_color":  "green",
					"text_color":   "red",
				},
			},
			widget: Widget{Name: "ga.bar_thingy", Theme: "super_theme",
				Options: map[string]string{"text_color": "red"},
			},
			project: project{
				themes: map[string]map[string]string{
					"super_theme": {
						"border_color": "blue",
						"title_color":  "green",
						"text_color":   "yellow",
					},
				},
			},
		},
		{
			name: "project with widget type theme",
			expected: Widget{Name: "ga.bar_thingy",
				Options: map[string]string{
					"border_color": "blue",
					"title_color":  "green",
					"text_color":   "yellow",
				},
			},
			widget: Widget{Name: "ga.bar_thingy"},
			project: project{
				themes: map[string]map[string]string{
					"bar": {
						"border_color": "blue",
						"title_color":  "green",
						"text_color":   "yellow",
					},
					"table": {
						"border_color": "red",
						"title_color":  "red",
						"text_color":   "yellow",
					},
				},
			},
		},
		{
			name: "project with widget unknown named theme",
			expected: Widget{Name: "ga.bar_thingy", Theme: "unknown_theme",
				Options: map[string]string{
					"border_color": "blue",
					"title_color":  "green",
					"text_color":   "yellow",
				},
			},
			widget: Widget{Name: "ga.bar_thingy", Theme: "unknown_theme"},
			project: project{
				themes: map[string]map[string]string{
					"bar": {
						"border_color": "blue",
						"title_color":  "green",
						"text_color":   "yellow",
					},
					"table": {
						"border_color": "red",
						"title_color":  "red",
						"text_color":   "yellow",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.project.addDefaultTheme(tc.widget)

			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Expected %v, actual %v", tc.expected, actual)
			}
		})
	}
}
