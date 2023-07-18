package xauth

var OAuthService *OAuths

type OAuths struct {
	OAuthInfos map[string]*OAuthInfo
}
