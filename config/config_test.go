package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetInt(t *testing.T) {
	bytes, err := json.Marshal(map[string]int{"key": 0})
	if err != nil {
		t.Fatal(err)
	}
	m := &map[string]interface{}{}
	err = json.Unmarshal(bytes, m)
	if err != nil {
		t.Fatal(err)
	}
	cfg = (*C)(m)

	// cfg = &C{}
	// cfg.Load("test.json")

	fmt.Println(cfg)
	// fmt.Println(cfg.GetInt("key"))
	// Fuck!
	v, ok := (*cfg)["key"].(int)
	fmt.Println(v, ok)
	fmt.Println(cfg.GetNum("key"))
	fmt.Println(cfg.GetInt("key"))
}
