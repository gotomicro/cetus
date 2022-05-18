package login

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/gotomicro/cetus/pkg/kauth"
)

func Oauth(c *gin.Context) {
	if kauth.OAuthService == nil {
		c.JSON(1, "oauth not enabled")
		return
	}

	name := c.Param("oauth")
	connect, ok := kauth.ConnectorMap[name]
	if !ok {
		c.JSON(1, fmt.Sprintf("No OAuth with name %s configured", name))
		return
	}
	state := c.Query("state")
	errorParam := c.Query("error")
	if errorParam != "" {
		errorDesc := c.Query("error_description")
		elog.Error("failed to login ", zap.Any("error", errorParam), zap.String("errorDesc", errorDesc))
		c.JSON(2, fmt.Sprintf("failed to login, errorParam: %s", errorParam))
		return
	}

	code := c.Query("code")
	if code == "" {
		state, err := kauth.GenStateString()
		if err != nil {
			elog.Error("Generating state string failed", zap.Error(err))
			c.JSON(3, "An internal error occurred")
			return
		}
		hashedState := kauth.HashStateCode(state, econf.GetString("app.secretKey"), kauth.OAuthService.OAuthInfos[name].ClientSecret)
		c.SetCookie(
			kauth.OauthStateCookieName,
			url.QueryEscape(hashedState),
			econf.GetInt("auth.OauthStateCookieMaxAge"),
			"/",
			"",
			false, // todo
			true,
		)

		if kauth.OAuthService.OAuthInfos[name].HostedDomain == "" {
			c.Redirect(http.StatusFound, connect.AuthCodeURL(state, oauth2.AccessTypeOnline))
			return
		} else {
			c.Redirect(http.StatusFound, connect.AuthCodeURL(state, oauth2.SetAuthURLParam("hd", kauth.OAuthService.OAuthInfos[name].HostedDomain), oauth2.AccessTypeOnline))
			return
		}
	}

	cookie, err := c.Cookie(kauth.OauthStateCookieName)
	if err != nil {
		c.JSON(4, "system error: "+err.Error())
		return
	}
	cookieState, _ := url.QueryUnescape(cookie)

	// delete cookie
	c.SetCookie(
		kauth.OauthStateCookieName,
		"",
		-1,
		"/",
		"",
		false, // todo
		true,
	)

	if cookieState == "" {
		c.JSON(5, "login.OAuthLogin(missing saved state)")
		return
	}

	queryState := kauth.HashStateCode(state, econf.GetString("app.secretKey"), kauth.OAuthService.OAuthInfos[name].ClientSecret)
	elog.Info("state check", zap.Any("queryState", queryState), zap.Any("cookieState", cookieState))
	if cookieState != queryState {
		c.JSON(6, "login.OAuthLogin(state mismatch)")
		return
	}

	// handle call back
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: kauth.OAuthService.OAuthInfos[name].TlsSkipVerify,
		},
	}
	oauthClient := &http.Client{
		Transport: tr,
	}

	if kauth.OAuthService.OAuthInfos[name].TlsClientCert != "" || kauth.OAuthService.OAuthInfos[name].TlsClientKey != "" {
		cert, err := tls.LoadX509KeyPair(kauth.OAuthService.OAuthInfos[name].TlsClientCert, kauth.OAuthService.OAuthInfos[name].TlsClientKey)
		if err != nil {
			elog.Error("Failed to setup TlsClientCert", zap.String("oauth", name), zap.Error(err))
			c.JSON(7, "login.OAuthLogin(Failed to setup TlsClientCert)")
			return
		}

		tr.TLSClientConfig.Certificates = append(tr.TLSClientConfig.Certificates, cert)
	}

	if kauth.OAuthService.OAuthInfos[name].TlsClientCa != "" {
		caCert, err := ioutil.ReadFile(kauth.OAuthService.OAuthInfos[name].TlsClientCa)
		if err != nil {
			elog.Error("Failed to setup TlsClientCa", zap.String("oauth", name), zap.Error(err))
			c.JSON(8, "login.OAuthLogin(Failed to setup TlsClientCa)")
			return
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tr.TLSClientConfig.RootCAs = caCertPool
	}

	oauthCtx := context.WithValue(context.Background(), oauth2.HTTPClient, oauthClient)

	// get token from provider
	token, err := connect.Exchange(oauthCtx, code)
	if err != nil {
		c.JSON(9, "login.OAuthLogin(NewTransportWithCode)"+err.Error())
		return
	}
	// token.TokenType was defaulting to "bearer", which is out of spec, so we explicitly set to "Bearer"
	token.TokenType = "Bearer"

	elog.Debug("OAuthLogin Got token", zap.Any("token", token))

	// set up oauth2 client
	client := connect.Client(oauthCtx, token)

	_, appSubURL, _ := kauth.ParseAppAndSubURL(econf.GetString("app.rootURL"))

	// get user info
	userInfo, err := connect.UserInfo(client, token)
	if err != nil {
		if _, ok := err.(*kauth.Error); ok {
			// todo
			c.Redirect(http.StatusFound, appSubURL+"/login")
			return
		} else {
			c.JSON(10, fmt.Sprintf("login.OAuthLogin(get info from %s error %s)", name, err.Error()))
			return
		}
	}

	elog.Debug("OAuthLogin got user info", zap.Any("userInfo", userInfo))

	// validate that we got at least an email address
	if userInfo.Email == "" {
		c.Redirect(http.StatusFound, appSubURL+"/login")
		return
	}

	// validate that the email is allowed to login to juno
	if !connect.IsEmailAllowed(userInfo.Email) {
		c.Redirect(http.StatusFound, appSubURL+"/login")
		return
	}

	// // TODO 存储用户数据
	// mysqlUser := &db.User{
	// 	Username:   userInfo.Name + "_" + name,
	// 	Nickname:   userInfo.Login + "_" + name,
	// 	Email:      userInfo.Email,
	// 	Oauth:      "oauth_" + name,
	// 	OauthId:    userInfo.Id,
	// 	OauthToken: db.OAuthToken{Token: token},
	// }
	// // create or update oauth user
	// err = user.User.CreateOrUpdateOauthUser(mysqlUser)
	// if err != nil {
	// 	c.JSONE(11, "create or update oauth user error", err.Error())
	// 	return
	// }
	// elog.Debug("OAuthLogin got user info", zap.Any("mysqlUserInfo", mysqlUser))
	//
	// err = user.Session.Save(c, mysqlUser)
	// if err != nil {
	// 	return elog.JSON(c, 12, "create or update oauth user error", err.Error())
	// }

	toURULCookie, err := c.Cookie("redirect_k_auth_to")
	toURL := appSubURL + "/"
	if err != nil {
		c.Redirect(http.StatusFound, toURL)
		return
	}

	toURL, err = url.QueryUnescape(toURULCookie)
	if err != nil {
		c.Redirect(http.StatusFound, toURL)
		return
	}
	c.Redirect(http.StatusFound, toURL)
	return
}
