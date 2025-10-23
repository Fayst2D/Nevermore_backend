package auth

type Token struct {
	At string `json:"access_token"`
	Rt string `json:"refresh_token"`
}
