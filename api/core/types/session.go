package types

const LoginUserID = "LOGIN_USER_ID"
const LoginUserCache = "LOGIN_USER_CACHE"

const UserAuthHeader = "Authorization"
const AdminAuthHeader = "Admin-Authorization"
const ChatTokenHeader = "Chat-Token"

// Session configs struct
type Session struct {
	SecretKey string
	MaxAge    int
}
