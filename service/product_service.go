package service

import (
	"Seckill/datamoudles"
	"Seckill/repositories"
)

type IProductService interface {
	GetProductByID(int64) (*datamoudles.Product, error)
	GetAllProduct() ([]*datamoudles.Product, error)
	DeleteProductByID(int64) bool
	InsertProduct(product *datamoudles.Product) (int64, error)
	UpdateProduct(product *datamoudles.Product) error
	SubNumberOne(productID int64) error
}

type ProductService struct {
	productRepositories repositories.IProduct
}

func (p *ProductService) GetProductByID(productID int64) (*datamoudles.Product, error) {
	return p.productRepositories.SelectByKey(productID)
}

func (p *ProductService) GetAllProduct() ([]*datamoudles.Product, error) {
	return p.productRepositories.SelectAll()
}

func (p *ProductService) DeleteProductByID(productID int64) bool {
	return p.productRepositories.Delete(productID)
}

func (p *ProductService) InsertProduct(product *datamoudles.Product) (int64, error) {
	return p.productRepositories.Insert(product)
}

func (p *ProductService) UpdateProduct(product *datamoudles.Product) error {
	return p.productRepositories.Update(product)
}

func NewProductService(repository repositories.IProduct) IProductService {
	return &ProductService{repository}
}

func (p *ProductService) SubNumberOne(productID int64) error {
	return p.productRepositories.SubProductNum(productID)
}
