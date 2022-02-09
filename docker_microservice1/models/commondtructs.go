package models

type DataTablesRequest struct {
	Draw    int `json:"draw"`
	Columns []struct {
		Data       string `json:"data"`
		Name       string `json:"name"`
		Searchable bool   `json:"searchable"`
		Orderable  bool   `json:"orderable"`
		Search     struct {
			Value string `json:"value"`
			Regex bool   `json:"regex"`
		} `json:"search"`
	} `json:"columns"`
	Order []struct {
		Column int    `json:"column"`
		Dir    string `json:"dir"`
	} `json:"order"`
	Start  int `json:"start"`
	Length int `json:"length"`
	Search struct {
		Value string `json:"value"`
		Regex bool   `json:"regex"`
	} `json:"search"`
}