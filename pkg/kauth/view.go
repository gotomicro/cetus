package kauth

type OAuthType int

const (
	GITHUB OAuthType = iota + 1
	GOOGLE
	GENERIC
	GRAFANA
	GITLAB
)

type BasicUserInfo struct {
	Id      string
	Name    string
	Email   string
	Login   string
	Company string
	Role    string
	Groups  []string
}

type OAuthInfo struct {
	Typ                  string
	ClientId             string
	ClientSecret         string
	Scopes               []string
	AuthUrl              string
	TokenUrl             string
	Enable               bool
	EmailAttributeName   string
	EmailAttributePath   string
	RoleAttributePath    string
	AllowedDomains       []string
	HostedDomain         string
	ApiUrl               string
	AllowSignup          bool
	Name                 string
	TlsClientCert        string
	TlsClientKey         string
	TlsClientCa          string
	TlsSkipVerify        bool
	TeamIds              []interface{}
	AllowedOrganizations []string
	AllowedGroups        []string
}
