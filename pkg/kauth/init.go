package kauth

var OAuthService *OAuths

type OAuths struct {
	OAuthInfos map[string]*OAuthInfo
}
