package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"server/config"
	"server/utils/logger"
	"time"
	"server/utils/uuid"
)

type Item struct {
	FlowId string `json:"flowid"`
	Phone string `json:"phone"`
	Channel string `json:"channel"`
	Signature string `json:"signature"`
	Message string `json:"message"`
	Arrived int `json:"arrived"`
	Reason string `json:"reason"`
	CreateDate   time.Time    `json:"createdate"`      //消息创建时间
}

type ItemResponse struct {
	FlowId string `json:"flowid"`
	Arrived int `json:"arrived"`
	Reason string `json:"reason"`
}

var mgoSession *mgo.Session

func SaveMessage(m *Item) (string,error){
	c := getConnection()
	m.FlowId = uuid.GetUUIDWithoutLine()
	err := c.Insert(m)
	if err != nil {
		logger.Error(nil, "insert message to mongo error %s", err)
		mgoSession = nil
		return "", nil
	}
	return m.FlowId,nil
}

func UpdateMessage(m *ItemResponse) error{
	c := getConnection()
	err := c.Update(bson.M{"flowid": m.FlowId}, bson.M{"$set": bson.M{"arrived": m.Arrived, "reason": m.Reason}})
	if err != nil {
		logger.Error(nil, "update message to mongo flowId:%s, error:%s",m.FlowId, err)
		mgoSession = nil
		return err
	}
	return nil
}

func FindMessage(flowId string) Item {
	c := getConnection()
	item := Item{}
	err := c.Find(bson.M{"flowid": flowId}).One(&item)
	if err != nil {
		logger.Error(nil, "get message from mongo error %s", err)
		mgoSession = nil
	}
	return item
}

func getConnection() *mgo.Collection {
	return getSession().DB("sms").C("sendlog")
}

func getSession() *mgo.Session {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(nil, "getSession error %s", err)
		}
	}()
	if mgoSession == nil {
		ms, err := mgo.Dial(config.CONFIG[config.ENV]["MONGO_ADDRESS"])
		if err != nil {
			panic(err) //直接终止程序运行
		}
		ms.SetMode(mgo.Monotonic, true)
		createIndex(ms) //初始化索引
		mgoSession = ms.Clone()
		return mgoSession
	}
	//最大连接池默认为4096
	return mgoSession.Clone()
}

func createIndex(session *mgo.Session) {
	c := session.DB("sms").C("sendlog")
	err := c.EnsureIndex(mgo.Index{
		Key:        []string{"flowid"},
		Unique:     true,
		Background: true,
		DropDups:   true,
	})
	if nil != err {
		logger.Error(nil, "crate index error %s", err)
	}
	err = c.EnsureIndex(mgo.Index{
		Key:        []string{"phone"},
		Unique:     false,
		Background: true,
		DropDups:   true,
	})
	if nil != err {
		logger.Error(nil, "crate index error %s", err)
	}
	err = c.EnsureIndex(mgo.Index{
		Key:        []string{"arrived"},
		Unique:     false,
		Background: true,
		DropDups:   true,
	})
	if nil != err {
		logger.Error(nil, "crate index error %s", err)
	}
	err = c.EnsureIndex(mgo.Index{
		Key:        []string{"signature"},
		Unique:     false,
		Background: true,
		DropDups:   true,
	})
	if nil != err {
		logger.Error(nil, "crate index error %s", err)
	}
	err = c.EnsureIndex(mgo.Index{
		Key:        []string{"channel"},
		Unique:     false,
		Background: true,
		DropDups:   true,
	})
	if nil != err {
		logger.Error(nil, "crate index error %s", err)
	}
	//err = c.EnsureIndex(mgo.Index{
	//	Key:        []string{"createdate"},
	//	Unique:     false,
	//	Background: true,
	//	//ExpireAfter: 60 * 60 * 24 * 14 * time.Second,  //缓存14天
	//	ExpireAfter: 60 * 60 * 24 * 14 * time.Second, //缓存14天
	//})
	//if nil != err {
	//	logger.Error(nil, "crate index error %s", err)
	//}
}
