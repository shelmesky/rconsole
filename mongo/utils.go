package mongo

import (
	"fmt"
)

func InsertOne(value interface{}) error {
	coll, err := GetCollection("connection")
	if err != nil {
		return fmt.Errorf("failed get mongo collection: %s", err)
	}

	err = coll.Insert(value)
	if err != nil {
		return fmt.Errorf("insert mongo failed: %s", err)
	}

	return nil
}
