package up

import (
	"bufio"
	"fmt"
	"net"
	"os/user"
	"regexp"
	"strings"
)

// App is a struct containing necessary fields for an application
type App struct {
	name       string
	identifier string
	routes     []Route
	server     *net.UnixListener
	socketName string
}

func (a *App) route(req *Request) Handler {
	for _, val := range a.routes {
		r := regexp.MustCompile(val.pattern)
		if r.MatchString(req.URI) {
			req.Parameters = extractURIParameters(r, req.URI)
			return val.handler
		}
	}
	return func404
}

// (a *App) serve ...
func (a *App) serve(conn net.Conn) {
	buf := bufio.NewReader(conn)
	rawreq := ""
	for {
		temp, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		if temp == "\n" && rawreq != "" {
			req, err := NewRequestFromRawString(rawreq)
			if err != nil || req.URI == "" || req.Method == "" {
				fmt.Println("Error in parsing request : ", err)
				conn.Write([]byte("Wrong Request\n"))
				rawreq = ""
				continue
			}
			conn.Write([]byte(a.route(&req)(&req) + "\n"))
			rawreq = ""
			continue
		}
		rawreq += temp
	}
}

// Serve runs the app with the supplied configuration
func (a *App) Run() {
	for {
		conn, err := a.server.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go a.serve(conn)
	}
}

// NewApp returns a App struct denoting the application
func NewApp(name, identifier string) (*App, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}
	socketName := "/tmp/" + user.Username + "_" + identifier + ".app"
	server, err := net.ListenUnix("unix", &net.UnixAddr{socketName, "unix"})
	if err != nil {
		return nil, err
	}
	return &App{
		name:       name,
		identifier: identifier,
		routes:     make([]Route, 0),
		server:     server,
		socketName: socketName,
	}, err
}

// HandleFunc function is used to add a route to the app.
// Route can contain variables as uri parameters.
func (a *App) HandleFunc(pattern string, handler Handler) {
	uriParts := strings.Split(pattern, "/")
	for i, val := range uriParts {
		if len(val) < 2 {
			continue
		}
		if val[0] == ':' {
			parameter := val[1:]
			if checkParamName(parameter) {
				uriParts[i] = "(?P<" + parameter + ">[[:alnum:]]+)"
			}
		}
	}
	a.routes = append(a.routes, Route{join(uriParts), handler})
}

// checkParamName returns a boolean based on whether the passed `name` is valid(only contains letters)
func checkParamName(name string) bool {
	r := regexp.MustCompile("[[:alpha:]]+")
	return r.MatchString(name)
}

// join returns a uri formed from joining the parts(uri splitted) passed
func join(parts []string) string {
	uri := "^/"
	for _, val := range parts {
		if val != "" {
			uri += val + "/"
		}
	}
	return uri[:len(uri)-1] + "$"
}
