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
要写route 因为他不属于[build-on order](https://caddyserver.com/docs/caddyfile/directives#directive-order)
