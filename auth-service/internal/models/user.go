package models

type User struct {
    ID       int32  `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}