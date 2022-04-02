package repositories

import (
	"Seckill/common"
	"Seckill/datamoudles"
	"database/sql"
	"errors"
	"strconv"
)

type IUserRepository interface {
	Conn() error
	Select(userName string) (*datamoudles.User, error)
	Insert(user *datamoudles.User) (userID int64, err error)
}

type UserManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func (u UserManagerRepository) Conn() (err error) {
	//TODO implement me
	if u.mysqlConn == nil {
		mysql, errMysql := common.NewMysqlConn()
		if errMysql != nil {
			return errMysql
		}
		u.mysqlConn = mysql
		if u.table == "" {
			u.table = "user"
		}
	}
	return
}

func (u UserManagerRepository) Select(userName string) (user *datamoudles.User, err error) {
	//TODO implement me
	if userName == "" {
		return &datamoudles.User{}, errors.New("用户名不能为空")
	}
	if err = u.Conn(); err != nil {
		return &datamoudles.User{}, err
	}
	sql := "Select form" + u.table + "where userName =?"
	rows, errRows := u.mysqlConn.Query(sql, userName)
	defer rows.Close()
	if errRows != nil {
		return &datamoudles.User{}, err
	}
	result := common.GetResultRow(rows)
	if len(result) == 0 {
		return &datamoudles.User{}, errors.New("用户不存在")
	}
	user = &datamoudles.User{}
	common.DataToStructByTagSql(result, user)
	return
}

func (u UserManagerRepository) Insert(user *datamoudles.User) (userID int64, err error) {
	//TODO implement me
	if err = u.Conn(); err != nil {
		return
	}
	sql := "INSERT" + u.table + "SET nickName = ?, userName = ?, passWord=?"
	stmt, errStmt := u.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return userID, errStmt
	}
	result, errResult := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if errResult != nil {
		return userID, errResult
	}
	return result.LastInsertId()
}

func (u *UserManagerRepository) SelectByID(userID int64) (user *datamoudles.User, err error) {
	if err = u.Conn(); err != nil {
		return &datamoudles.User{}, err
	}
	sql := "select *from" + u.table + "where ID= ?" + strconv.FormatInt(userID, 10)
	row, errRow := u.mysqlConn.Query(sql)
	if errRow != nil {
		return &datamoudles.User{}, errRow
	}

	result := common.GetResultRow(row)

	if len(result) == 0 {
		return &datamoudles.User{}, errors.New("用户不存在")
	}
	user = &datamoudles.User{}
	common.DataToStructByTagSql(result, user)
	return
}

func NewUserRepository(table string, db *sql.DB) IUserRepository {
	return &UserManagerRepository{table, db}
}
