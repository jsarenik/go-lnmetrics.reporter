package db

import (
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"

	log "github.com/OpenLNMetrics/go-metrics-reported/pkg/log"
)

type database struct {
	instance *leveldb.DB
}

var instance database

func GetInstance() *database {
	return &instance
}

func (this *database) InitDB(homedir string) error {
	path := homedir + "/db"
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.GetInstance().Error(err)
		return err
	}
	log.GetInstance().Info("Created database at " + path)
	this.instance = db
	return nil
}

func (this *database) PutValue(key string, value interface{}) error {
	json, err := json.Marshal(&value)
	if err != nil {
		log.GetInstance().Error(err)
		return err
	}
	log.GetInstance().Debug(
		fmt.Sprintf("Storing value with key %s and value %s", key, json))
	return this.instance.Put([]byte(key), []byte(json), nil)
}

func (this *database) GetValue(key string) (string, error) {
	value, err := this.instance.Get([]byte(key), nil)
	if err != nil {
		log.GetInstance().Error(err)
		return "", err
	}
	return string(value), nil
}

func (this *database) DeleteValue(key string) error {
	return this.instance.Delete([]byte(key), nil)
}

// TODO Add method to iterate over a method.
