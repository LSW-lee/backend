package rest

import (
	"bytes"
	"log"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RunAPI(address string) error {
	h, err := NewHandler()
	if err != nil {
		return err
	}
	return RunAPIWithHandler(address, h)
}

func RunAPIWithHandler(address string, h HandlerInterface) error {
	//Gin객체 초기화
	r := gin.Default()
	r.Use(MyCustomMiddleware())
	r.GET("/products", h.GetProducts)
	r.GET("/promos", h.GetPromos)

	//그룹으로 묶어서 라우팅이 가능하다.
	userGroup := r.Group("/user")
	{
		userGroup.POST("/:id/signout", h.SignOut)
		userGroup.GET("/:id/orders", h.GetOrders)
	}

	usersGroup := r.Group("/users")
	{
		usersGroup.POST("/charge", h.Charge)
		usersGroup.POST("/signin", h.SignIn)
		usersGroup.POST("", h.AddUser)
	}
	r.Static("/img", "../public/img")

	//r.Use(static.ServeRoot("/", "../public/build"))
	return r.Run(address)
}

func MyCustomMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*요청을 실행하기 전 코드 영역*/

		c.Set("v", "123") //c.Get("v")를 하면 "123"이 반환된다.
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		//요청 처리 로직 실행
		c.Next()

		//응답 코드 확인
		log.Println(c.Request.RequestURI, ":", blw.body.String())
	}
}
