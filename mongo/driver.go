package mongo

import (
	"fmt"
	"github.com/shelmesky/rconsole/utils"
	"gopkg.in/mgo.v2"
	"time"
)

var (
	MgoTimeout time.Duration /* 全局的MongoDB连接超时时间 */

	MgoURL string /* 全局的MongoDB连接字符串 */

	MgoSession  *mgo.Session  /* 全局的MongoDB会话对象 */
	MgoDatabase *mgo.Database /* 全局的MongoDB数据库对象 */

	/*
	   全局标志，指示是否连接到MongoDB.
	   由TouchMongoDB或者封装的Find/Update等函数设置
	   当调用mgo的任何函数/方法失败时，应尽早设置此标志
	*/
	MgoConnected bool
)

/*
设置MongoDB为已连接状态
*/
func SetMgoOK() {
	MgoConnected = true
}

/*
设置MongoDB为未连接状态
*/
func SetMgoFailed() {
	MgoConnected = false
}

/*
根据参数中指定的MongoDB连接字符串
连接到MongoDB

@URL: MongoDB的连接字符串

返回:
@*mgo.Session 连接成功后MongoDB的会话对象
@error 错误信息
*/
func ConnectMongoDB(URL string, duration time.Duration) (*mgo.Session, error) {
	var session *mgo.Session
	var err error

	session, err = mgo.DialWithTimeout(URL, duration)

	return session, err
}

/*
切换到参数中指定的数据库

@DBName: 数据库的名称
*/
func ChangeDatabase(DBName string) {
	MgoDatabase = MgoSession.DB(DBName)
}

/*
获取参数中指定的Collection对象

@CollectionName: Collection的名称

返回:
Collection的对象
*/
func GetCollection(CollectionName string) (*mgo.Collection, error) {
	var coll *mgo.Collection

	if MgoConnected == false {
		return coll, fmt.Errorf("MongoDB connection has lost")
	}

	coll = MgoDatabase.C(CollectionName)

	return coll, nil
}

/*
周期性的检查到MongoDB的连接
如果检查失败，则尝试重新连接
*/
func TouchMongoDB(show_check_info bool) {
	c := time.Tick(10 * time.Second)
	for _ = range c {
		if MgoSession != nil && MgoDatabase != nil {
			err := MgoSession.Ping()
			if err != nil {
				// 当MongoDB无法连接时，关闭Session
				if MgoSession != nil && MgoDatabase == nil {
					utils.Println("MongoDB connection lost, clear resources and close session.")
					MgoDatabase.Logout()
					MgoSession.LogoutAll()
					MgoSession.Close()
					MgoDatabase = nil
					MgoSession = nil
				}
				SetMgoFailed()
				utils.Println("Check connection: MongoDB has gone, try to reconnect.")
				InitMongoDB(MgoURL, MgoTimeout)
			} else {
				SetMgoOK()

				if show_check_info == true {
					utils.Println("Check MongoDB connection... OK")
				}
			}
		} else {
			SetMgoFailed()
			utils.Println("Check connection: MongoDB has gone, try to reconnect.")
			InitMongoDB(MgoURL, MgoTimeout)
		}
	}
}

/*
初始化到MongoDB的连接
连接成功后，切换到URL中指定的数据库

@URL: MongoDB的连接字符串
*/
func InitMongoDB(URL string, timeout time.Duration) error {
	MgoURL = URL
	MgoTimeout = timeout

	session, err := ConnectMongoDB(URL, timeout)
	if err != nil {
		SetMgoFailed()
		return fmt.Errorf("Can not connect to MongoDB: %s", err)
	}

	MgoSession = session
	ChangeDatabase("")

	err = MgoSession.Ping()
	if err != nil {
		SetMgoFailed()
		return fmt.Errorf("Check connection: MongoDB has gone, try to reconnect.")
	}

	SetMgoOK()
	return nil
}
