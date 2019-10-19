package internal

import "strings"

const (
	// Widget config options
	optionTitle      = "title"
	optionTitleColor = "title_color"

	// Time
	optionStartDate  = "start_date"
	optionEndDate    = "end_date"
	optionTimePeriod = "time_period"
	optionGlobal     = "global"

	// Tables
	optionRowLimit  = "row_limit"
	optionCharLimit = "character_limit"

	// Metrics
	optionDimension  = "dimension"
	optionDimensions = "dimensions"
	optionMetrics    = "metrics"
	optionMetric     = "metric"

	// Ordering
	optionOrder = "order"

	// Filtering
	optionFilters = "filters"

	// Display
	optionContent = "content"

	// Repository
	optionRepository = "repository"
	optionOwner      = "owner"

	// Owner / all
	optionScope = ownerScope

	ownerScope = "owner"
)

// A widget is a representation of a set of data in the screen.
type Widget struct {
	Name    string            `mapstructures:"name"`
	Size    string            `mapstructures:"size"`
	Options map[string]string `mapstructures:"options"`
	Theme   string            `mapstructures:"theme"`
}

func (w *Widget) typeID() string {
	n := strings.Split(w.Name, ".")[1]
	return strings.Split(n, "_")[0]
}

func (w *Widget) serviceID() string {
	return strings.Split(w.Name, ".")[0]
}
