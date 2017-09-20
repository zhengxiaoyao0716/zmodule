package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetInt(t *testing.T) {
	bytes, err := json.Marshal(map[string]int{"key": 1, "int64": 1 << 31}) // 2147483648
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
	f, ok := (*cfg)["key"].(float64)
	fmt.Println(f, ok)
	fmt.Println(cfg.GetNum("key"))

	i, ok := (*cfg)["key"].(int)
	fmt.Println(i, ok)
	fmt.Println(cfg.GetInt("key"))

	l, ok := (*cfg)["int64"].(int64)
	fmt.Println(l, ok)
	fmt.Println(cfg.GetI64("int64"))

	i, ok = (*cfg)["int64"].(int)
	fmt.Println(i, ok)
	fmt.Println(cfg.GetInt("int64"))
}
