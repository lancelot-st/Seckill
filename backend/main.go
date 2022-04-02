package main

import (
	"Seckill/backend/web/controllers"
	"Seckill/common"
	"Seckill/repositories"
	"Seckill/service"
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"log"
)

func main() {
	//创建iris实例
	app := iris.New()
	//设置错误模式
	app.Logger().SetLevel("debug")
	//注册模板
	tmplate := iris.HTML("./backend/web/views", ".html").Layout(
		"shared/layout.html").Reload(
		true)
	app.RegisterView(tmplate)
	//设置模板目标
	app.HandleDir("/assets", "./backend/web/assets")
	//出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message",
			ctx.Values().GetStringDefault("message", "访问的页面出错"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Println(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productService := service.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderMangerRepository("order", db)
	orderService := service.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))
	//启动服务

	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
