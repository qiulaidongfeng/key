// Package key 实现一种密钥轮换方案
package key

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"os"
	"unsafe"

	"golang.org/x/crypto/argon2"
)

var gcm cipher.AEAD
var gcm2 cipher.AEAD
var oldAeskey [32]byte

// AesKey 是当前的最新密钥
var Aeskey [32]byte

func init() {
	Test := os.Getenv("TEST") != ""
	if !Test {
		Init()
	}
}

func Init() {
	main_key := os.Getenv("main_key")
	if main_key == "" {
		panic("环境变量main_key应该提供主密钥")
	}
	main_key2 := os.Getenv("main_key2")
	s := salt()

	if main_key2 != "" {
		genKey(main_key, s, &oldAeskey, &gcm)
		genKey(main_key2, s, &Aeskey, &gcm2)
	} else {
		genKey(main_key, s, &Aeskey, &gcm)
	}
}

func genKey(key string, salt []byte, result *[32]byte, gcm *cipher.AEAD) {
	aes_key := argon2.IDKey([]byte(key), salt, 2, 64*1024, 4, 32)
	block, err := aes.NewCipher(aes_key)
	if err != nil {
		panic(err)
	}
	*result = [32]byte(aes_key)
	tmp, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		panic(err)
	}
	*gcm = tmp
}

func salt() []byte {
	salt, err := os.ReadFile("./salt")
	if err != nil {
		if os.IsNotExist(err) {
			var salt [32]byte
			rand.Read(salt[:])
			err := os.WriteFile("./salt", salt[:], 0600)
			if err != nil {
				panic(err)
			}
			return salt[:]
		}
		panic(err)
	}
	return salt
}

// Encrypt 使用 AES-256-GCM 加密数据
func Encrypt(v string) string {
	if gcm2 != nil {
		ev := gcm2.Seal(nil, nil, unsafe.Slice(unsafe.StringData(v), len(v)), nil)
		return unsafe.String(unsafe.SliceData(ev), len(ev))
	}
	ev := gcm.Seal(nil, nil, unsafe.Slice(unsafe.StringData(v), len(v)), nil)
	return unsafe.String(unsafe.SliceData(ev), len(ev))
}

// Decrypt 使用 AES-256-GCM 解密数据
// 如果解密失败，返回空字符串
func Decrypt(v string) string {
	ev, err := gcm.Open(nil, nil, unsafe.Slice(unsafe.StringData(v), len(v)), nil)
	if err != nil {
		ev, _ := gcm2.Open(nil, nil, unsafe.Slice(unsafe.StringData(v), len(v)), nil)
		return unsafe.String(unsafe.SliceData(ev), len(ev))
	}
	return unsafe.String(unsafe.SliceData(ev), len(ev))
}
