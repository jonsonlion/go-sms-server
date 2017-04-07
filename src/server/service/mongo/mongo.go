package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"server/config"
	"server/utils/logger"
	"time"
)

type Item struct {
	Topic        string                 //消息topic
	ApnsId       string                 //消息编号
	Token        string                 //消息token
	Body         string                 //内容
	Sound        string                 //声音
	CustomFields map[string]interface{} //自定义消息
	Arrived      int32                  //是否到达，0否：1：是
	StatusCode   string                 //返回码
	Reason       string                 //原因
	CreateDate   time.Time              //消息创建时间
}

type ItemResponse struct {
	ApnsId     string //消息编号
	Arrived    int32  //是否到达，0否：1：是
	StatusCode string //返回码
	Reason     string //原因
}

var mgoSession *mgo.Session

func SaveMessage(m Item) {
	c := getConnection()
	err := c.Insert(&m)
	if err != nil {
		logger.Error(nil, "insert message to mongo error %s", err)
		mgoSession = nil
	}
}

func updateMessage(m ItemResponse) {
	c := getConnection()
	err := c.Update(bson.M{"apnsid": m.ApnsId}, bson.M{"$set": bson.M{"arrived": m.Arrived, "reason": m.Reason, "statuscode": m.StatusCode}})
	if err != nil {
		logger.Error(nil, "update message to mongo error %s", err)
		mgoSession = nil
	}
}

func FindMessage(apnsid string) Item {
	c := getConnection()
	item := Item{}
	err := c.Find(bson.M{"apnsid": apnsid}).One(&item)
	if err != nil {
		logger.Error(nil, "get message from mongo error %s", err)
		mgoSession = nil
	}
	return item
}

func getConnection() *mgo.Collection {
	return getSession().DB("push").C("messages")
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
	c := session.DB("push").C("messages")
	err := c.EnsureIndex(mgo.Index{
		Key:        []string{"apnsid"},
		Unique:     true,
		Background: true,
		DropDups:   true,
	})
	if nil != err {
		logger.Error(nil, "crate index error %s", err)
	}
	err = c.EnsureIndex(mgo.Index{
		Key:        []string{"statuscode"},
		Unique:     false,
		Background: true,
		DropDups:   true,
	})
	if nil != err {
		logger.Error(nil, "crate index error %s", err)
	}
	err = c.EnsureIndex(mgo.Index{
		Key:        []string{"createdate"},
		Unique:     false,
		Background: true,
		//ExpireAfter: 60 * 60 * 24 * 14 * time.Second,  //缓存14天
		ExpireAfter: 60 * 60 * 24 * 14 * time.Second, //缓存14天
	})
	if nil != err {
		logger.Error(nil, "crate index error %s", err)
	}
}
