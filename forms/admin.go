package forms

type AdminRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type AdminLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}