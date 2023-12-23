package request

type UserRequest struct {
	UserNameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
}
