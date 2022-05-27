package consistenthash

import (
	"fmt"
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	consistenthash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// We intentionally create a hash function, and the replica will be:
	// 2, 4, 6, 12, 14, 16, 22, 24, 26
	consistenthash.Add("6", "2", "4")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if consistenthash.Get(k) != v {
			t.Errorf("Asking for key %s, expect a value %s\n", k, v)
		}
	}

	// now add 8, 18, 28
	consistenthash.Add("8")

	// 27 should now map to 8
	testCases["27"] = "8"

	for k, v := range testCases {
		if consistenthash.Get(k) != v {
			fmt.Println("hello:" + consistenthash.Get(k))
			t.Errorf("Asking for key %s, expect a value %s\n", k, v)
		}
	}
}
