package openobserve

type Dashboard struct {
	Version     int      `json:"version"`
	DashboardID string   `json:"dashboardId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Role        string   `json:"role"`
	Owner       string   `json:"owner"`
	Created     string   `json:"created"`
	Tabs        []Tab    `json:"tabs"`
	Variables   Variable `json:"variables"`
}

type Tab struct {
	TabID  string  `json:"tabId"`
	Name   string  `json:"name"`
	Panels []Panel `json:"panels"`
}

type Panel struct {
	ID              string      `json:"id"`
	Type            string      `json:"type"`
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	QueryType       string      `json:"queryType"`
	Queries         []Query     `json:"queries"`
	Layout          PanelLayout `json:"layout"`
	HTMlContent     string      `json:"htmlContent"`
	MarkdownContent string      `json:"markdownContent"`
	Config          Config      `json:"config"`
}

type Config struct {
	ShowLegends    bool    `json:"show_legends"`
	Decimals       int     `json:"decimals"`
	AxisBorderShow bool    `json:"axis_border_show"`
	Unit           *string `json:"unit,omitempty"`
}

type Query struct {
	Query       string      `json:"query"`
	CustomQuery bool        `json:"customQuery"`
	Config      QueryConfig `json:"config"`
	Fields      QueryFields `json:"fields"`
}

type QueryConfig struct {
	PromqlLegend string `json:"promql_legend"`
	LayerType    string `json:"layer_type"`
	WeightFixed  int    `json:"weight_fixed"`
}

type QueryFields struct {
	Stream     string   `json:"stream"`
	StreamType string   `json:"stream_type"`
	X          []string `json:"x"`
	Y          []string `json:"y"`
	Z          []string `json:"z"`
	Filter     []string `json:"filter"`
}

type PanelLayout struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
	I int `json:"i"`
}

type Variable struct {
	List               []VariableListItem `json:"list"`
	ShowDynamicFilters bool               `json:"showDynamicFilters"`
}

type VariableListItem struct {
	Type      string            `json:"type"`
	Name      string            `json:"name"`
	Label     string            `json:"label"`
	QueryData VariableQueryData `json:"query_data"`
	Value     string            `json:"value"`
	Options   []any             `json:"options"`
}

type VariableQueryData struct {
	StreamType string `json:"stream_type"`
	Stream     string `json:"stream"`
	Field      string `json:"field"`
}
