环境：`docker run -d --entrypoint /bin/sh -v ./:/app -w /app --network host caddy:builder -c "sleep infinity"`

命令行：`docker exec -it <name> /bin/sh`

---

运行：`xcaddy run`

编译：`xcaddy build`

> `xcaddy build` 只能有一个插件 所以多个插件只能共用一个init（似乎是这样的）

---

```golang
func (Clipboard) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.clipboard", // here
		New: func() caddy.Module { return new(Clipboard) },
	}
}
```
要对应 Caddyfile directive

---

Caddyfile
```conf
:81 {
    handle * {
        route { // here
            clipboard 
        }
    }
}
```
要写route 因为他不属于[build-in order](https://caddyserver.com/docs/caddyfile/directives#directive-order)
