package grafana

type Dashboard struct {
	Description string     `json:"description"`
	Title       string     `json:"title"`
	UID         string     `json:"uid"`
	Version     int        `json:"version"`
	Panels      []Panel    `json:"panels"`
	Templating  Templating `json:"templating"`
}

type Panel struct {
	ID      int      `json:"id"`
	GridPos GridPos  `json:"gridPos"`
	Title   string   `json:"title"`
	Type    string   `json:"type"`
	Panels  []Panel  `json:"panels"`
	Format  string   `json:"format"`
	Targets []Target `json:"targets"`
	Stack   bool     `json:"stack"`
	YAxes   []YAxes  `json:"yaxes"`
}

type GridPos struct {
	H int `json:"h"`
	W int `json:"w"`
	X int `json:"x"`
	Y int `json:"y"`
}

type Target struct {
	Expr         string `json:"expr"`
	Format       string `json:"format"`
	LegendFormat string `json:"legendFormat"`
}

type Templating struct {
	List []TemplatingListItem `json:"list"`
}

type TemplatingListItem struct {
	Type string `json:"type"`
	// Query      string `json:"query"`
	Name       string `json:"name"`
	Definition string `json:"definition"`
}

type YAxes struct {
	Format string `json:"format"`
}
