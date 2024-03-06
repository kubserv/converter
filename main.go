package main

import (
	"encoding/json"
	"flag"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kubserv/converter/internal/grafana"
	"github.com/kubserv/converter/internal/openobserve"
)

var (
	inputFileFlag  = flag.String("input", "", "Input file")
	outputFileFlag = flag.String("output", "", "Output file")

	// extract dynamic variables from grafana dashboard: label_values(stream, field)
	// since in OpenObserve we can only select a stream right now and not build upon
	// a previous variable, we strip out the labels that might be in there.
	// Example:
	// label_values(redis_up{namespace=\"$namespace\"}, pod)
	// extract the stream "redis_up" and its tag/field "pod"
	templateListRegexp = regexp.MustCompile(`label_values\((\w+)(?:\{[^}]*\})?,\s*(\w+)\)`)
)

func main() {
	flag.Parse()

	// read grafana dashboard data
	data, err := os.ReadFile(*inputFileFlag)
	if err != nil {
		panic(err)
	}
	var grafanaDashboard grafana.Dashboard
	err = json.Unmarshal(data, &grafanaDashboard)
	if err != nil {
		panic(err)
	}

	// initial openobserve setup
	dashboard := openobserve.Dashboard{
		Version:     grafanaDashboard.Version,
		DashboardID: grafanaDashboard.UID,
		Title:       grafanaDashboard.Title,
		Description: grafanaDashboard.Description,
		Created:     time.Now().Format(time.RFC3339),
		// right now let's only support a single tab
		// I think converting a Grafana "row" panel to a tab could be a possibility
		Tabs: []openobserve.Tab{{
			TabID: "default",
			Name:  "Default",
		}},
		Variables: openobserve.Variable{
			List:               []openobserve.VariableListItem{},
			ShowDynamicFilters: true,
		},
	}

	for _, panel := range grafanaDashboard.Panels {
		if panel.Type == "row" {
			// we don't care about rows, later those could be tabs
			if len(panel.Panels) == 0 {
				continue
			}

			for _, subPanel := range panel.Panels {
				dashboard.Tabs[0].Panels = append(dashboard.Tabs[0].Panels, grafanaPanelToOpenobservePanel(subPanel))
			}
		} else {
			dashboard.Tabs[0].Panels = append(dashboard.Tabs[0].Panels, grafanaPanelToOpenobservePanel(panel))
		}
	}

	// get variables
	for _, variable := range grafanaDashboard.Templating.List {
		matches := templateListRegexp.FindStringSubmatch(variable.Definition)
		if len(matches) == 0 {
			continue
		}

		dashboard.Variables.List = append(dashboard.Variables.List, openobserve.VariableListItem{
			Type: "query_values",
			Name: variable.Name,
			QueryData: openobserve.VariableQueryData{
				StreamType: "metrics",
				Stream:     matches[1],
				Field:      matches[2],
			},
			Options: []any{},
		})
	}

	// now we go through each panel and increment the I from the last + 1
	for i, panel := range dashboard.Tabs[0].Panels {
		panel.Layout.I = i + 1
		dashboard.Tabs[0].Panels[i] = panel
	}

	js, err := json.Marshal(dashboard)
	if err != nil {
		panic(err)
	}

	if *outputFileFlag != "" {
		err = os.WriteFile(*outputFileFlag, js, 0644)
		if err != nil {
			panic(err)
		}
	} else {
		os.Stdout.Write(js)
	}
}

func grafanaPanelToOpenobservePanel(panel grafana.Panel) openobserve.Panel {
	// very basic right now, need to do some more
	// investigation
	panelType := "line"
	switch panel.Type {
	case "graph":
		panelType = "line"
	case "singlestat":
		panelType = "metric"
	case "table":
		panelType = "table"
	case "table-old":
		panelType = "table"
	case "gauge":
		panelType = "gauge"
	}

	// small hack
	if panel.Stack {
		panelType = "area-stacked"
	}

	p := openobserve.Panel{
		ID:          strconv.Itoa(panel.ID),
		Type:        panelType,
		Title:       panel.Title,
		Description: panel.Format,
		QueryType:   "promql",
		// OpenObserve uses a grid twice as big as Grafana
		Layout: openobserve.PanelLayout{
			X: panel.GridPos.X * 2,
			Y: panel.GridPos.Y * 2,
			W: panel.GridPos.W * 2,
			H: panel.GridPos.H * 2,
			I: 1,
		},
		Config: openobserve.Config{
			ShowLegends:    true,
			Decimals:       2,
			AxisBorderShow: false,
		},
		Queries: []openobserve.Query{},
	}

	if len(panel.YAxes) > 0 {
		p.Config.Unit = &panel.YAxes[0].Format
	}

	for _, target := range panel.Targets {

		// some clean-up for legend, because OpenObserve uses single
		// curly braces for templating, while Grafana uses double...
		// and some other clean-up
		legendFormat := target.LegendFormat
		legendFormat = strings.ReplaceAll(legendFormat, "{{", "{")
		legendFormat = strings.ReplaceAll(legendFormat, "}}", "}")
		legendFormat = strings.ReplaceAll(legendFormat, "{ ", "{")
		legendFormat = strings.ReplaceAll(legendFormat, " }", "}")

		// replace $interval with $__interval, because that was
		// in some templates
		target.Expr = strings.ReplaceAll(target.Expr, "$interval", "$__interval")

		p.Queries = append(p.Queries, openobserve.Query{
			Query:       target.Expr,
			CustomQuery: true,
			Fields: openobserve.QueryFields{
				Stream:     "",
				StreamType: "metrics",
				X:          []string{},
				Y:          []string{},
				Z:          []string{},
				Filter:     []string{},
			},
			Config: openobserve.QueryConfig{
				PromqlLegend: legendFormat,
				LayerType:    "line",
				WeightFixed:  0,
			},
		})
	}

	return p
}
