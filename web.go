package main

/*
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS: -L. -lsqlite3 -lsql
#include "sql.h"
*/
import "C"
import (
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"
	"os"
	"sync"
	"syscall"
	"time"
)

var LOCK = sync.Mutex{}

func requestHandler(ctx *fasthttp.RequestCtx) {
	fmt.Printf("Handler running at %v\n", syscall.Getpid())
	fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
	fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
	fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
	fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
	fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
	fmt.Fprintf(ctx, "Post content type is %q\n", ctx.PostArgs())
	fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
	fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
	fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
	fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
	fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())
	LOCK.Lock()
	var i = C.insert(C.CString("test.db"), C.CString(fmt.Sprintf("insert into test values(%q, %q, %q, %q, %q, %q, %q)", ctx.Method(), ctx.RequestURI(), ctx.Path(), ctx.PostArgs(), ctx.UserAgent(), ctx.RemoteIP(), time.Now().Format(time.RFC850))))
	LOCK.Unlock()
	if i == 0 {
		fmt.Println("data recorded!")
	} else {
		fmt.Println("something wrong!")
	}
	fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)
	ctx.SetContentType("text/plain; charset=utf8")
	// Set arbitrary headers
	ctx.Response.Header.SetStatusCode(200)
	ctx.Response.Header.Set("X-My-Header", "my-header-value")
	// Set cookies
	var c fasthttp.Cookie
	c.SetKey("This is a cookie")
	c.SetValue("You get yourself a cookie")
	ctx.Response.Header.SetCookie(&c)
}
func main() {
	fmt.Printf("main function running at %v \n", syscall.Getpid())
	var port string
	flag.StringVar(&port, "p", "8080", "the port that server runs on")
	var mode bool
	flag.BoolVar(&mode, "m", false, "If m is set to true, then it listens on 0.0.0.0, else it listens on 127.0.0.1")
	flag.Parse()
	var debug, exists = os.LookupEnv("WEB_DEBUG")
	if exists {
		fmt.Printf("debug mode is %q\n", debug)
	}
	var ip = "127.0.0.1"
	if mode {
		ip = "0.0.0.0"
	}
	var path = ip + ":" + port
	var err = fasthttp.ListenAndServe(path, requestHandler)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
