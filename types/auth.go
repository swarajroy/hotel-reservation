package types

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateAuthParams(email, password string) AuthParams {
	return AuthParams{
		Email:    email,
		Password: password,
	}
}
