package service

import (
	"Seckill/datamoudles"
	"Seckill/repositories"
)

type IOrderService interface {
	GetOrderByID(int64) (*datamoudles.Order, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(*datamoudles.Order) error
	InsertOrder(*datamoudles.Order) (int64, error)
	GetAllOrder() ([]*datamoudles.Order, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
	InsertOrderByMessage(*datamoudles.Message) (int64, error)
}

func NewOrderService(repository repositories.IOrderRepository) IOrderService {
	return &OrderService{repository}
}

type OrderService struct {
	OrderRepository repositories.IOrderRepository
}

func (o *OrderService) GetOrderByID(orderID int64) (order *datamoudles.Order, err error) {
	//TODO implement me
	return o.OrderRepository.SelectByKey(orderID)
}

func (o *OrderService) DeleteOrderByID(orderID int64) bool {
	//TODO implement me
	isOK := o.OrderRepository.Delete(orderID)
	return isOK
}

func (o *OrderService) UpdateOrder(order *datamoudles.Order) (err error) {
	//TODO implement me
	return o.OrderRepository.Update(order)
}

func (o *OrderService) InsertOrder(order *datamoudles.Order) (orderID int64, err error) {
	//TODO implement me
	return o.OrderRepository.Insert(order)
}

func (o *OrderService) GetAllOrder() ([]*datamoudles.Order, error) {
	//TODO implement me
	return o.OrderRepository.SelectAll()
}

func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
	//TODO implement me
	return o.OrderRepository.SelectAllWithInfo()
}

//根据消息创造订单
func (o *OrderService) InsertOrderByMessage(message *datamoudles.Message) (orderID int64, err error) {
	order := &datamoudles.Order{
		UserId:      message.UserID,
		ProductId:   message.ProductID,
		OrderStatus: datamoudles.OrderSuccess,
	}
	return o.InsertOrder(order)

}
