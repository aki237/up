# up

Server framework on top of Unix Domain Sockets (implementation in Golang).
This can be used to run desktop apps as a separate process serving the HTML, 
CSS, Js files and a frontend WebKit or Gecko Process can render the view of the app.

## Why this can be better?

Nowadays a lot of desktop apps are being written on top of electron framework. First of, 
bundling a whopping 40 Megs of stuff for each app is just space killer. And This can be implemented in any
insane language or scripting available thus lifting some the limitations of Javascript. There is a limitation 
that this can only be used in Unix (or Linux).


## Sample Usage of this package

```go
package main

import (
	"fmt"

	"github.com/aki237/up"
)

func main() {
	app, err := up.NewApp("Test", "org.testorg.test")
	if err != nil {
		fmt.Println(err)
		return
	}
	app.HandleFunc("/welcome/:name", f)
	app.Run()
}

func f(r *up.Request) up.Response {
	return up.Response("Hello, " + r.Parameters["name"] + "!!")
}
```

Test this using the `socat` tool : (`socat` is like `telnet` for Unix Domain Sockets)

```shell
$ socat - UNIX:/tmp/$USER_org.testorg.test.app
{"method" : "GET", "uri" : "/welcome/John"}

Hello, John!!
```

### Thoughts

+ Use `msgpack` for data format
+ Modify electron framework or WebKit to work with this format natively without
  any dirty signal interupts and data feeds.
