```azure
	e.HTMLRender = xpongo.New(
		xpongo.WithDebug(gin.IsDebugging()),
		xpongo.WithGlobalContext(map[string]interface{}{
			"settings": configs,
		}),
		xpongo.WithFS(template.FS),
	)
```
