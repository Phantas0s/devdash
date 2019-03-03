package main

import (
	"reflect"
	"testing"

	"github.com/Phantas0s/devdash/internal"
)

// Each slice dimension represent row, col or widgets.
// For rows slice: [RowIndex][ColIndex][WidgetIndex]widget
// For size slice: [RowIndex][ColIndex]widget - no need the index of widgets
func Test_OrderWidgets(t *testing.T) {
	testCases := []struct {
		name          string
		project       Project
		expectedSizes [][]string
		expectedRows  [][][]internal.Widget
		wantErr       bool
	}{
		{
			name: "1 row 2 col 3 widgets",
			project: Project{
				Name: "test",
				Widgets: []Row{
					Row{
						Row: []Column{
							Column{
								Col: []Widgets{
									Widgets{
										Size: "XL",
										Elements: []internal.Widget{
											internal.Widget{
												Name: "row_1_col_1_widget1",
												Size: "XL",
												Options: map[string]string{
													"start_date": "7_days_ago",
													"end_date":   "today",
												},
											},
										},
									},
									Widgets{
										Size: "L",
										Elements: []internal.Widget{
											internal.Widget{
												Name: "row_1_col_1_widget2",
												Size: "L",
											},
										},
									},
								},
							},
							Column{
								Col: []Widgets{
									Widgets{
										Size: "M",
										Elements: []internal.Widget{
											internal.Widget{
												Name: "row_1_col_2_widget1",
												Size: "M",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedSizes: [][]string{
				[]string{ // row 1
					"XL",
					"L",
					"M",
				},
			},
			expectedRows: [][][]internal.Widget{
				[][]internal.Widget{ // row 1
					[]internal.Widget{ // column 1
						internal.Widget{
							Name: "row_1_col_1_widget1",
							Size: "XL",
							Options: map[string]string{
								"start_date": "7_days_ago",
								"end_date":   "today",
							},
						},
						internal.Widget{
							Name: "row_1_col_1_widget2",
							Size: "L",
						},
					},
					[]internal.Widget{ //column 2
						internal.Widget{
							Name: "row_1_col_2_widget1",
							Size: "M",
						},
					},
				},
			},
		},
		{
			name: "3 row 3 col 5 widgets",
			project: Project{
				Name: "test",
				Widgets: []Row{
					Row{
						Row: []Column{
							Column{
								Col: []Widgets{
									Widgets{
										Size: "S",
										Elements: []internal.Widget{
											internal.Widget{
												Name: "row_1_col_1_widget1",
												Size: "S",
											},
										},
									},
									Widgets{
										Size: "XS",
										Elements: []internal.Widget{
											internal.Widget{
												Name: "row_1_col_1_widget2",
												Size: "XS",
											},
										},
									},
								},
							},
							Column{
								Col: []Widgets{
									Widgets{
										Size: "M",
										Elements: []internal.Widget{
											internal.Widget{
												Name: "row_1_col_2_widget1",
												Size: "M",
											},
										},
									},
								},
							},
						},
					},
					Row{
						Row: []Column{
							Column{
								Col: []Widgets{
									Widgets{
										Size: "S",
										Elements: []internal.Widget{
											internal.Widget{
												Name: "row_2_col_1_widget1",
												Size: "S",
											},
										},
									},
								},
							},
							Column{
								Col: []Widgets{
									Widgets{
										Size: "L",
										Elements: []internal.Widget{
											internal.Widget{
												Name: "row_2_col_2_widget1",
												Size: "L",
											},
										},
									},
									Widgets{
										Size: "XL",
										Elements: []internal.Widget{
											internal.Widget{
												Name: "row_2_col_2_widget2",
												Size: "XL",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedSizes: [][]string{
				[]string{ // row 1
					"S",
					"XS",
					"M",
				},
				[]string{
					"S",
					"L",
					"XL",
				},
			},
			expectedRows: [][][]internal.Widget{
				[][]internal.Widget{ // row 1
					[]internal.Widget{ // column 1
						internal.Widget{
							Name: "row_1_col_1_widget1",
							Size: "S",
						},
						internal.Widget{
							Name: "row_1_col_1_widget2",
							Size: "XS",
						},
					},
					[]internal.Widget{ //column 2
						internal.Widget{
							Name: "row_1_col_2_widget1",
							Size: "M",
						},
					},
				},
				[][]internal.Widget{ // row 2
					[]internal.Widget{ // column 1
						internal.Widget{
							Name: "row_2_col_1_widget1",
							Size: "S",
						},
					},
					[]internal.Widget{ // column 1
						internal.Widget{
							Name: "row_2_col_2_widget1",
							Size: "L",
						},
						internal.Widget{
							Name: "row_2_col_2_widget2",
							Size: "XL",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualRows, actualSizes := tc.project.OrderWidgets()

			if tc.wantErr == false && !reflect.DeepEqual(actualSizes, tc.expectedSizes) {
				t.Errorf("Expected sizes %v, actual %v", tc.expectedSizes, actualSizes)
			}

			if tc.wantErr == false && !reflect.DeepEqual(actualRows, tc.expectedRows) {
				t.Errorf("Expected rows %v, actual %v", tc.expectedRows, actualRows)
			}
		})
	}
}
