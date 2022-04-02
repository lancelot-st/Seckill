package main

import (
	"Seckill/common"
	"Seckill/datamoudles"
	"Seckill/encrypt"
	"Seckill/rabbitmq"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

//设置集群地址， 最好内外IP
var hostArray = []string{"127.0.0.1", "127.0.0.1"}
var localHost = "127.0.0.1"
var port = "8083"
var hashConsistent *common.Consistent

//数量控制接口服务器内网IP， 或者SLB内网IP
var GetOneIp = "127.0.0.1"

var GetOnePort = "8084"

//rabbitmq
var rabbitMqValidate *rabbitmq.RabbitMQ

type AccessControl struct {
	sourcesArray map[int]interface{}
	*sync.RWMutex
}

var accessControl = &AccessControl{sourcesArray: make(map[int]interface{})}

func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	data := m.sourcesArray[uid]
	return data
}

//设置记录
func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	m.sourcesArray[uid] = time.Now()
	m.RWMutex.Unlock()
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}

	//根据用户id，判断获取具体机器

	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}
	//判断是否为本机
	if hostRequest == localHost {
		//执行本机数据读取和校验
		return m.GetDataFromMap(uid.Value)
	} else {
		//不是本机充当代理访问数据返回结果
		return GetDataFromOtherMap(hostRequest, req)
	}
}

type BlackList struct {
	listArray map[int]bool
	sync.RWMutex
}

var blackList = &BlackList{listArray: make(map[int]bool)}

//获取黑名单
func (m *BlackList) GetBlackListByID(uid int) bool {
	m.RLock()
	defer m.RUnlock()
	return m.listArray[uid]
}

//添加黑名单
func (m *BlackList) SetBlackListByID(uid int) bool {
	m.Lock()
	defer m.Unlock()
	m.listArray[uid] = true
	return true
}

func (m *AccessControl) GetDataFromMap(uid string) (isOk bool) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}

	//添加黑名单

	if blackList.GetBlackListByID(uidInt) {
		//判断是否被添加到黑名单
		return false
	}

	dataRecord := m.GetNewRecord(uidInt)
	if dataRecord == nil {
		return true
	}
	return
}

func GetDataFromOtherMap(host string, request *http.Request) bool {
	hostUrl := "http://" + host + "port" + "/checkRight"
	response, body, err := Geturl(hostUrl, request)
	if err != nil {
		return false
	}
	//判断状态
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}

//统一验证拦截器，每个接口都需要提前验证
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("执行验证")
	//添加基于cookie的验证
	return nil
}

func CheckUserInfo(r *http.Request) error {
	//获取uid, cookie
	uidCokkie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("用户uid Cookie 获取失败")
	}
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("用户加密 Cookie 获取是失败")
	}
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("加密串已被篡改")
	}
	fmt.Println("结果比对")
	fmt.Println("用户ID")
	if checkInfo(uidCokkie.Value, string(signByte)) {
		return nil
	}
	return errors.New("身份校验失败")
}

//自定义逻辑判断
func checkInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false

}

func Geturl(hostUrl string, request *http.Request) (response *http.Response, body []byte, err error) {
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return
	}
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return
	}
	//模拟接口访问
	client := &http.Client{}
	req, err := http.NewRequest("GET",
		hostUrl, nil)
	if err != nil {
		return
	}
	//手动指定，排查多余cookie
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	//添加cookie到模拟的请求中
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	response, err = client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}

func CheckRight(w http.ResponseWriter, r *http.Request) {
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}

//执行验证，验证失败就不执行Auth了
func Check(w http.ResponseWriter, r *http.Request) {
	//执行正常业务逻辑
	fmt.Println("执行check")
	queryFrom, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryFrom["productID"]) <= 0 {
		w.Write([]byte("false"))
		return
	}
	productString := queryFrom["productID"][0]
	fmt.Println(productString)
	//获取用户cookie
	userCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	//分布式权限验证

	right := accessControl.GetDistributedRight(r)
	if right == false {
		w.Write([]byte("false"))
	}
	//获取数量控制权限防止秒杀出现超卖
	hostUrl := "http://" + GetOneIp + ":" + GetOnePort + "/getOne"
	responseValidate, validateBody, err := Geturl(hostUrl, r)
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	//判断数量接口请求状态

	if responseValidate.StatusCode == 200 {
		if string(validateBody) == "true" {
			//整合下单
			productID, err := strconv.ParseInt(productString, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			//获取用户ID
			//获取商品ID
			userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			//创建消息体
			message := datamoudles.NewMessage(userID, productID)
			//类型转换
			byteMessage, err := json.Marshal(message)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			//生产消息
			err = rabbitMqValidate.PublishSimple(string(byteMessage))
			if err != nil {
				w.Write([]byte("false"))
			}
			w.Write([]byte("true"))
			return
		}
	}

	w.Write([]byte("false"))
	return
}

func main() {
	//负载均衡过滤器的设置
	//采用一致性hash算法
	hashConsistent := common.NewConsistent()
	//采用一致性hash算法添加节点
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	localIp, err := common.GetIntranceIp()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIp
	fmt.Println(localHost)

	rabbitMqValidate = rabbitmq.NewRabbitMQSimple("imoocProduct")
	defer rabbitMqValidate.Destory()

	//设置静态文件目录
	http.Handle("/html/", http.StripPrefix(
		"/html/", http.FileServer(http.Dir("./fronted/web/htmlProductShow"))))

	http.Handle("/public/", http.StripPrefix("/html/", http.FileServer(http.Dir("/fronted/web/public"))))

	//过滤器
	filter := common.NewFilter()
	//注册拦截器
	filter.RegisterFilterUri("/check", Auth)
	filter.RegisterFilterUri("/checkRight", Auth)
	//启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/check", filter.Handle(CheckRight))
	//启动服务
	http.ListenAndServe(":8083", nil)
}
