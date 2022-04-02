package main

import (
	"Seckill/common"
	"Seckill/rabbitmq"
	"Seckill/repositories"
	"Seckill/service"
	"fmt"
)

func main() {
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}
	//创建product数据库操作实例
	product := repositories.NewProductManager("product", db)
	//创建 product service
	productService := service.NewProductService(product)
	//创建Order数据库实例
	order := repositories.NewOrderMangerRepository("order", db)

	//创建order service
	orderService := service.NewOrderService(order)

	rabbitmqConsumeSimple := rabbitmq.NewRabbitMQSimple("imoocProduct")

	rabbitmqConsumeSimple.ConsumeSimple(orderService, productService)

}
