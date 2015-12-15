package types

type InstallCoreos struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Profile       string `json:"profile,omitempty"`
	Comport       string `json:"comport,omitempty"`
	Hostname      string `json:"hostname"`
	SshKey        string `json:"sshkey"`
	EtcdToken     string `json:"etcdToken"`
	CompletionUri string `json:"completionUri,omitempty"`
}
