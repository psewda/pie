package models

type Version struct {
	Version   string `json:"version"`
	Golang    string `json:"golang"`
	GitCommit string `json:"gitCommit"`
	Built     string `json:"built"`
	OsArch    string `json:"osArch"`
}
