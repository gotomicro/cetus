# README

## 配置

> 	r.GET("/login/:oauth", core.Handle(login.Oauth))

```
[app]
secretKey = "ASDFASDFASDF" # hashStatecode
rootURL = "http://localhost:9001"
baseURL = "/login/"
```