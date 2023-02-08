package helper

import (
	"encoding/json"
	"math/rand"
	"time"
)

func ContainsString(l []string, s string) bool {
	for _, a := range l {
		if a == s {
			return true
		}
	}
	return false
}

func StringToStruct(s string, outputStruct interface{}) error {
	return BytesToStruct([]byte(s), outputStruct)
}

func BytesToStruct(bytes []byte, outputStruct interface{}) error {
	return json.Unmarshal(bytes, outputStruct)
}

func RandomInList(list []string) string {
	rand.Seed(time.Now().Unix())
	return list[rand.Intn(len(list))]
}

func Struct2Struct(inputStruct interface{}, outputStruct interface{}) error {
	bytes, _ := json.Marshal(inputStruct)
	return json.Unmarshal(bytes, outputStruct)
}

func Struct2String(inputStruct interface{}) string {
	bytes, _ := json.Marshal(inputStruct)
	return string(bytes)
}

func Struct2Map(inputStruct interface{}) map[string]string {
	var inInterface map[string]string
	inrec, _ := json.Marshal(inputStruct)
	_ = json.Unmarshal(inrec, &inInterface)
	return inInterface
}
