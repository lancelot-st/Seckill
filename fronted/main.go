package main

import (
	"Seckill/common"
	"Seckill/fronted/middlerware"
	"Seckill/fronted/web/controllers"
	"Seckill/rabbitmq"
	"Seckill/repositories"
	"Seckill/service"
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"log"
	"time"
)

func main() {
	//创建iris实例
	app := iris.New()
	//设置错误模式
	app.Logger().SetLevel("debug")
	//注册模板
	tmplate := iris.HTML("./fronted/web/views", ".html").Layout(
		"shared/layout.html").Reload(
		true)
	app.RegisterView(tmplate)
	//设置模板目标
	app.HandleDir("/public", "./fronted/web/public")
	//访问生成好的静态文件
	app.HandleDir("/html", "./fronted/web/htmlProductShow")
	//出现异常跳转到指定页面
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.ViewData("message",
			ctx.Values().GetStringDefault("message", "访问的页面出错"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	//连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Println(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sess := sessions.New(sessions.Config{
		Cookie:  "helloword",
		Expires: 600 * time.Minute,
	})
	user := repositories.NewUserRepository("user", db)
	userService := service.NewService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx, sess, sess.Start)
	userPro.Handle(new(controllers.UserController))

	_ = rabbitmq.NewRabbitMQSimple("imoocProduct")

	product := repositories.NewProductManager("product", db)
	productService := service.NewProductService(product)
	order := repositories.NewOrderMangerRepository("table", db)
	orderService := service.NewOrderService(order)
	proProduct := app.Party("/product")
	proProduct.Use(middlerware.AuthConProduct)
	pro := mvc.New(proProduct)
	pro.Register(productService, orderService)
	pro.Handle(new(controllers.ProductController))

	app.Run(
		iris.Addr("localhost:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
