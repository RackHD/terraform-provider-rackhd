package types

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	UID      int    `json:"uid"`
}
