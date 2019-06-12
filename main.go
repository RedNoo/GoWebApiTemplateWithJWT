package main

import (
	"GoWebApiTemplateWithJWT/lib"
	user "GoWebApiTemplateWithJWT/model"
	"fmt"

	"github.com/buaazp/fasthttprouter"
	logger "github.com/savsgio/go-logger"
	"github.com/valyala/fasthttp"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func init() {
	logger.SetLevel(logger.DEBUG)
}

func Index(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")
	fmt.Fprint(ctx, "<h1>Hola, estas en el Index...<h1>")
}

func Login(ctx *fasthttp.RequestCtx) {
	qUser := []byte("savsgio")
	qPasswd := []byte("mypasswd")
	fasthttpJwtCookie := ctx.Request.Header.Cookie("fasthttp_jwt")

	// for example, server receive token string in request header.
	if len(fasthttpJwtCookie) == 0 {
		tokenString, expireAt := lib.CreateToken(qUser, qPasswd)

		// Set cookie for domain
		cookie := fasthttp.AcquireCookie()
		cookie.SetKey("fasthttp_jwt")
		cookie.SetValue(tokenString)
		cookie.SetExpire(expireAt)
		ctx.Response.Header.SetCookie(cookie)
	}

	ctx.Redirect("/", ctx.Response.StatusCode())
}

func main() {	
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=onder dbname=testdb password=1234")
	if err != nil {
		panic("failed to connect database")
	}

	db.DropTableIfExists(&user.User{}, "users")
	db.AutoMigrate(&user.User{})

	router := fasthttprouter.New()
	router.GET("/login", Login)
	router.GET("/", lib.Middleware(Index))

	server := &fasthttp.Server{
		Name:    "JWTTestServer",
		Handler: router.Handler,
	}

	logger.Debug("Listening in http://localhost:12345...")
	logger.Fatal(server.ListenAndServe(":12345"))
}
