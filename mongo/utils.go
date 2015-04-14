package mongo

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

func InsertOne(collection string, value interface{}) error {
	coll, err := GetCollection(collection)
	if err != nil {
		return fmt.Errorf("failed get mongo collection: %s", err)
	}

	err = coll.Insert(value)
	if err != nil {
		return fmt.Errorf("insert mongo failed: %s", err)
	}

	return nil
}

func QueryOne(collection string, query interface{}, result interface{}) error {
	coll, err := GetCollection(collection)
	if err != nil {
		return fmt.Errorf("failed get mongo collection: %s", err)
	}

	err = coll.Find(query).One(result)
	if err != nil {
		return err
	}

	return nil
}

type ConnType struct {
	Type string `bson:"type"`
}

func GetConnTypeByUUID(uuid string) (string, error) {
	var conn_type ConnType

	err := QueryOne("connection", bson.M{"uuid": uuid}, &conn_type)
	if err != nil {
		return "", err
	}
	return conn_type.Type, nil
}
