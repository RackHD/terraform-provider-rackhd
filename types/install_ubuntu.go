package types

type InstallUbuntu struct {
	Comport       string `json:"comport"`
	CompletionUri string `json:"completionUri"`
	Domain        string `json:"domain"`
	Hostname      string `json:"hostname"`
	Password      string `json:"password"`
	Profile       string `json:"profile"`
	UID           int    `json:"uid"`
	Username      string `json:"username"`
}
