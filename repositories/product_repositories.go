package repositories

import (
	"Seckill/common"
	"Seckill/datamoudles"
	"database/sql"
	"strconv"
)

type IProduct interface {
	//连接数据库
	Conn() error
	Insert(*datamoudles.Product) (int64, error)
	Delete(int64) bool
	Update(*datamoudles.Product) error
	SelectByKey(int64) (*datamoudles.Product, error)
	SelectAll() ([]*datamoudles.Product, error)
	SubProductNum(productID int64) error
}

type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

//数据库连接
func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
		if p.table == "" {
			p.table = "product"
		}
	}
	return
}

//数据库插入
func (p *ProductManager) Insert(product *datamoudles.Product) (productId int64, err error) {
	//判断连接是否存在
	if err = p.Conn(); err != nil {

		return
	}
	//准备SQL
	sql := "INSERT Product SET productName= ?, productNum=?, productImage= ?,  productUrl= ?"
	stmt, errSql := p.mysqlConn.Prepare(sql)
	if err != nil {
		return 0, errSql
	}
	//传入参数
	result, errStmt := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if errStmt != nil {
		return 0, errStmt
	}

	return result.LastInsertId()
}

//数据库的删除
func (p *ProductManager) Delete(productId int64) bool {
	if err := p.Conn(); err != nil {
		return false
	}
	sql := "DELETE from product where ID =?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(productId)
	if err != nil {
		return false
	}
	return true
}

//商品的更新
func (p *ProductManager) Update(product *datamoudles.Product) error {
	if err := p.Conn(); err != nil {
		return err
	}
	sql := "Update product Set productName = ? ,productNum = ?, productImage = ?, productUrl = ? where ID =" +
		strconv.FormatInt(product.ID, 10)

	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return err
	}
	return nil

}

//根据ID查询商品
func (p *ProductManager) SelectByKey(productID int64) (productResult *datamoudles.Product, err error) {
	if err = p.Conn(); err != nil {
		return &datamoudles.Product{}, err
	}
	sql := "Select from" + p.table + "where ID=?" + strconv.FormatInt(productID, 10)

	row, errRow := p.mysqlConn.Query(sql)
	defer row.Close()
	if errRow != nil {
		return &datamoudles.Product{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamoudles.Product{}, err
	}
	common.DataToStructByTagSql(result, productResult)
	return

}

//获取所有商品
func (p *ProductManager) SelectAll() (productArray []*datamoudles.Product, errProduct error) {
	if err := p.Conn(); err != nil {
		return nil, err
	}
	sql := "Select *from" + p.table
	rows, err := p.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err

	}
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, nil
	}
	for _, v := range result {
		product := &datamoudles.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)

	}
	return
}
func NewProductManager(table string, Db *sql.DB) IProduct {
	return &ProductManager{table, Db}
}

func (p *ProductManager) SubProductNum(productID int64) error {
	if err := p.Conn(); err != nil {
		return err
	}
	sql := "update " + p.table + " set " + " productNum=productNum-1 where ID =" + strconv.FormatInt(productID, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}
