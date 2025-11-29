package key_test

import (
	"os"
	"testing"

	"github.com/qiulaidongfeng/key"
)

func TestOneKey(t *testing.T) {
	os.Setenv("main_key", "tmp")
	key.Init()
	r1 := key.Aeskey

	os.Setenv("main_key", "tmp")
	key.Init()
	r2 := key.Aeskey

	if r1 != r2 {
		t.Fatalf("%v %v", r1, r2)
	}

	os.Setenv("main_key", "tmp2")
	key.Init()
	r3 := key.Aeskey
	if r1 == r3 {
		t.Fatalf("%v %v", r1, r3)
	}

	key.Decrypt("")
}

func TestTwoKey(t *testing.T) {
	os.Setenv("main_key", "tmp")
	key.Init()
	r1 := key.Aeskey
	e1 := key.Encrypt("tmp")

	os.Setenv("main_key2", "tmp")
	key.Init()
	r2 := key.Aeskey

	if r1 != r2 {
		t.Fatalf("%v %v", r1, r2)
	}

	os.Setenv("main_key2", "tmp2")
	key.Init()
	r3 := key.Aeskey
	e2 := key.Encrypt("tmp2")
	if r1 == r3 {
		t.Fatalf("%v %v", r1, r3)
	}

	if key.Decrypt(e1) != "tmp" {
		t.Fail()
	}
	if key.Decrypt(e2) != "tmp2" {
		t.Fail()
	}

	os.Setenv("main_key", "tm")
	key.Init()

	if key.Decrypt(e1) != "" {
		t.Fail()
	}
	if key.Decrypt(e2) != "tmp2" {
		t.Fail()
	}
}
