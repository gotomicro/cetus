package xauth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gotomicro/ego/core/elog"
	"github.com/jmespath/go-jmespath"
	"github.com/spf13/cast"
	"golang.org/x/oauth2"

	"github.com/gotomicro/cetus/kerror"
)

// Connector ..
type Connector interface {
	Type() int
	UserInfo(client *http.Client, token *oauth2.Token) (*BasicUserInfo, error)
	IsEmailAllowed(email string) bool
	IsSignupAllowed() bool

	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	Client(ctx context.Context, token *oauth2.Token) *http.Client
	TokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource
}

type Base struct {
	*oauth2.Config
	log            *elog.Component
	allowSignup    bool
	allowedDomains []string
}

var (
	ConnectorMap = make(map[string]Connector)
)

func newBase(name string, config *oauth2.Config, info *OAuthInfo) *Base {
	return &Base{
		Config:         config,
		log:            elog.DefaultContainer().Build(elog.WithFileName(name + ".log")),
		allowSignup:    info.AllowSignup,
		allowedDomains: info.AllowedDomains,
	}
}

func NewOAuthService(appURL, baseURL string, auths []OAuthInfo) {
	OAuthService = &OAuths{}
	OAuthService.OAuthInfos = make(map[string]*OAuthInfo)
	for _, info := range auths {
		name := info.Typ
		if !info.Enable {
			continue
		}
		OAuthService.OAuthInfos[name] = &info
		config := oauth2.Config{
			ClientID:     info.ClientId,
			ClientSecret: info.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  info.AuthUrl,
				TokenURL: info.TokenUrl,
			},
			RedirectURL: strings.TrimSuffix(appURL, "/") + baseURL + name,
			Scopes:      info.Scopes,
		}
		// GitHub.
		if name == "github" {
			ConnectorMap["github"] = &SocialGithub{
				Base:                 newBase(name, &config, &info),
				apiUrl:               info.ApiUrl,
				teamIds:              cast.ToIntSlice(info.TeamIds),
				allowedOrganizations: info.AllowedOrganizations,
			}
		}

		// GitLab.
		if name == "gitlab" {
			ConnectorMap["gitlab"] = &SocialGitlab{
				Base:          newBase(name, &config, &info),
				apiUrl:        info.ApiUrl,
				allowedGroups: info.AllowedGroups,
			}
		}

		// Google.
		if name == "google" {
			ConnectorMap["google"] = &SocialGoogle{
				Base:         newBase(name, &config, &info),
				hostedDomain: info.HostedDomain,
				apiUrl:       info.ApiUrl,
			}
		}

		// Generic - Uses the same scheme as Github.
		if name == "generic_oauth" {
			ConnectorMap["generic_oauth"] = &SocialGenericOAuth{
				Base:                 newBase(name, &config, &info),
				apiUrl:               info.ApiUrl,
				emailAttributeName:   info.EmailAttributeName,
				emailAttributePath:   info.EmailAttributePath,
				roleAttributePath:    info.RoleAttributePath,
				teamIds:              cast.ToIntSlice(info.TeamIds),
				allowedOrganizations: info.AllowedOrganizations,
			}
		}
	}
}

func (s *Base) IsEmailAllowed(email string) bool {
	return isEmailAllowed(email, s.allowedDomains)
}

func (s *Base) IsSignupAllowed() bool {
	return s.allowSignup
}

func isEmailAllowed(email string, allowedDomains []string) bool {
	if len(allowedDomains) == 0 {
		return true
	}

	valid := false
	for _, domain := range allowedDomains {
		emailSuffix := fmt.Sprintf("@%s", domain)
		valid = valid || strings.HasSuffix(email, emailSuffix)
	}

	return valid
}

func (s *Base) searchJSONForAttr(attributePath string, data []byte) (string, error) {
	if attributePath == "" {
		return "", errors.New("no attribute path specified")
	}

	if len(data) == 0 {
		return "", errors.New("empty user info JSON response provided")
	}

	var buf interface{}
	if err := json.Unmarshal(data, &buf); err != nil {
		return "", kerror.Wrap("failed to unmarshal user info JSON response", err)
	}

	val, err := jmespath.Search(attributePath, buf)
	if err != nil {
		return "", kerror.WrapF(err, "failed to search user info JSON response with provided path: %q", attributePath)
	}

	strVal, ok := val.(string)
	if ok {
		return strVal, nil
	}

	return "", nil
}

var (
	OauthStateCookieName = "oauth_state"
)

func GenStateString() (string, error) {
	rnd := make([]byte, 32)
	if _, err := rand.Read(rnd); err != nil {
		fmt.Printf("failed to generate state string error: %v", err)
		return "", err
	}
	return base64.URLEncoding.EncodeToString(rnd), nil
}

func HashStateCode(secretKey, code, seed string) string {
	hashBytes := sha256.Sum256([]byte(code + secretKey + seed))
	return hex.EncodeToString(hashBytes[:])
}

func ParseAppAndSubURL(rootURL string) (string, string, error) {
	appURL := rootURL
	if appURL[len(appURL)-1] != '/' {
		appURL += "/"
	}
	// Check if has app suburl.
	urlParse, err := url.Parse(appURL)
	if err != nil {
		fmt.Println("invalid root url", appURL, err.Error())
		return "", "", nil
	}
	appSubURL := strings.TrimSuffix(urlParse.Path, "/")
	return appURL, appSubURL, nil
}
