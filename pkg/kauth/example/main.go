package main

import (
	"github.com/gotomicro/cetus/pkg/kauth"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egin"

	"github.com/gotomicro/cetus/pkg/kauth/example/login"
)

func main() {
	err := ego.New().
		Invoker(
			AuthInit,
		).
		Serve(
			GetRouter(),
		).Run()
	if err != nil {
		elog.Panic("start up error: " + err.Error())
	}
}

func AuthInit() error {
	oauthInfos := make([]kauth.OAuthInfo, 0)

	_ = econf.UnmarshalKey("auth.tps", &oauthInfos)
	elog.Info("AuthInit", elog.Any("step", "UnmarshalKey"), elog.Any("oauthInfos", oauthInfos))

	appURL, appSubURL, _ := kauth.ParseAppAndSubURL(econf.GetString("app.rootURL"))
	baseURL := econf.GetString("app.baseURL")

	elog.Info("AuthInit", elog.Any("step", "ParseAppAndSubURL"), elog.Any("appURL", appURL), elog.Any("appSubURL", appSubURL))

	kauth.NewOAuthService(appURL, baseURL, oauthInfos)
	return nil
}

func GetRouter() *egin.Component {
	r := egin.Load("server.http").Build()
	r.GET("/login/:oauth", login.Oauth)
	return r
}
