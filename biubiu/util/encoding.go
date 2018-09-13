package util

import (
	"encoding/gob"
	"os"
)

func Serilization(filePath string, instance interface{}) (err error){
	file, err := os.Create(filePath)
	if err != nil{
		return err
	}
	enc := gob.NewEncoder(file)
	err = enc.Encode(instance)
	return
}
func DeSerilization(filePath string, structConstruct interface{}) (err error){
	file, err := os.Open(filePath)
	if err != nil{
		return err
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(structConstruct)
	return
}
