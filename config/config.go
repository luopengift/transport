package config

type InputConfig struct {
	Type string `json:"type"`
}
type FilterConfig struct{}
type OutputConfig struct{}
type RuntimeConfig struct {
	DEBUG    bool `json:"DEBUG"`
	MAXPROCS int  `json:"MAXPROCS"`
}
type ApiConfig struct {
}

type Config struct {
	Runtime RuntimeConfig `json:"runtime"`
	Input   InputConfig   `json:"input"`
	Filter  FilterConfig  `json:"filter"`
	Output  OutputConfig  `json:"output"`
	Api     ApiConfig     `json:"api"`
}
