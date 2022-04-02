package repositories

import (
	"Seckill/common"
	"Seckill/datamoudles"
	"database/sql"
	"strconv"
)

type IOrderRepository interface {
	Conn() error
	Insert(order *datamoudles.Order) (int64, error)
	Delete(int64) bool
	Update(order *datamoudles.Order) error
	SelectByKey(int64) (*datamoudles.Order, error)
	SelectAll() ([]*datamoudles.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

func NewOrderMangerRepository(table string, sql *sql.DB) IOrderRepository {
	return &OrderMangerRepository{table: table, mysqlConn: sql}
}

type OrderMangerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func (o *OrderMangerRepository) Conn() error {
	//TODO implement me

	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err

		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "order"
	}
	return nil
}

func (o *OrderMangerRepository) Insert(order *datamoudles.Order) (productID int64, err error) {
	//TODO implement me
	if err := o.Conn(); err != nil {
		return 0, err
	}
	sql := "INSERT" + o.table + "set userID=?, productID=?, orderStatus=?"
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return productID, err
	}

	result, errResult := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if errResult != nil {
		return productID, errResult
	}
	return result.LastInsertId()
}

func (o *OrderMangerRepository) Delete(orderID int64) (isOK bool) {
	//TODO implement me
	if err := o.Conn(); err != nil {
		return
	}
	sql := "delete from " + o.table + "where ID = ?"
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return
	}
	_, err := stmt.Exec(orderID)
	if err != nil {
		return
	}
	return true
}

func (o *OrderMangerRepository) Update(order *datamoudles.Order) (err error) {
	//TODO implement me
	if errConn := o.Conn(); errConn != nil {
		return errConn
	}
	sql := "Update" + o.table + "set UserID =?, productID= ?, orderStatus=?" +
		"where ID = " + strconv.FormatInt(order.ID, 10)
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return errStmt
	}
	_, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	return

}

func (o *OrderMangerRepository) SelectByKey(orderID int64) (order *datamoudles.Order, err error) {
	//TODO implement me
	if errConn := o.Conn(); errConn != nil {
		return &datamoudles.Order{}, errConn
	}
	sql := "Select from" + o.table + "where ID = ?" + strconv.FormatInt(orderID, 10)
	row, errRow := o.mysqlConn.Query(sql)
	if errRow != nil {
		return &datamoudles.Order{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamoudles.Order{}, err
	}
	order = &datamoudles.Order{}
	common.DataToStructByTagSql(result, order)
	return
}

func (o *OrderMangerRepository) SelectAll() (orderArray []*datamoudles.Order, err error) {
	//TODO implement me
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}
	sql := "Select * from" + o.table
	rows, errRows := o.mysqlConn.Query(sql)
	if errRows != nil {
		return nil, errRows
	}
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}
	for _, v := range result {
		order := &datamoudles.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return
}

func (o *OrderMangerRepository) SelectAllWithInfo() (OrderMap map[int]map[string]string, err error) {
	//TODO implement me
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}
	sql := "Select o.Id,p.productName, o.orderStatus From imooc.order as o left join product as " +
		"p on o.productID =p.ID"
	rows, errRows := o.mysqlConn.Query(sql)
	if errRows != nil {
		return nil, nil
	}
	return common.GetResultRows(rows), err
}
