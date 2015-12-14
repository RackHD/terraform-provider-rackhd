package types

type InstallOS struct {
	Domain       string `json:"domain"`
	Hostname     string `json:"hostname"`
	RootPassword string `json:"rootPassword"`
	Users        []User `json:"users"`
	Version      string `json:"version"`
}
