package controllers

import (
	"Seckill/datamoudles"
	"Seckill/rabbitmq"
	"Seckill/service"
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService service.IProductService
	OrderService   service.IOrderService
	RabbitMQ       *rabbitmq.RabbitMQ
	Session        *sessions.Session
}

var (
	htmlOutPath  = "./fronted/web/htmlProductShow" //生成HTML保存目录
	templatePath = "/fronted/web/views/template"   //静态文件
)

func (p *ProductController) GetGenerateHtml() {
	//获取模板文件地址
	contentTmp, err := template.ParseFiles(filepath.Join(templatePath), "product.html")
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//获取HTML生成路径
	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")
	//获取模板渲染路径
	productString := p.Ctx.URLParam("productID")
	productID, err := strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(int64(productID))

	//生成静态文件
	generateStaticHtml(p.Ctx, contentTmp, fileName, product)
}

//生成HTML静态文件
func generateStaticHtml(ctx iris.Context, template *template.Template,
	fileName string, product *datamoudles.Product) {
	if Exist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			ctx.Application().Logger().Debug(err)
		}
	}
	//生成静态文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		ctx.Application().Logger().Debug(err)
	}
	defer file.Close()
	template.Execute(file, &product)
}

//判断静态文件是否存在
func Exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func (p *ProductController) GetDetail() mvc.View {
	product, err := p.ProductService.GetProductByID(1)
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func (p *ProductController) GetOrder() []byte {
	productString := p.Ctx.URLParam("productID")
	userString := p.Ctx.GetCookie("uid")
	productID, err := strconv.ParseInt(productString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	userID, err := strconv.ParseInt(userString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)

	}
	//创建消息体
	message := datamoudles.NewMessage(userID, productID)
	//类型转换
	byteMessage, err := json.Marshal(message)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	err = p.RabbitMQ.PublishSimple(string(byteMessage))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)

	}
	return []byte("true")
	/*userID, err := strconv.Atoi(userString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	var orderID int64
	showMessage := "抢购失败！"
	if product.ProductNum > 0 {
		//创建订单
		product.ProductNum -= 1
		err := p.ProductService.UpdateProduct(product)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		order := &datamoudles.Order{
			UserId:      int64(userID),
			ProductId:   int64(productID),
			OrderStatus: datamoudles.OrderSuccess,
		}
		orderID, err = p.OrderService.InsertOrder(order)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		} else {
			showMessage = "抢购成功！"
		}
	}*/
}
