# Key

一个go密钥管理库

## 开发理论

假设密钥轮换过程为

1. 使用旧密钥（有一个最大有效期，到期所有被加密数据失效）
2. 新旧密钥共存（优先新密钥加密，解密两个密钥都失败才算解密失败）
3. 仅使用新密密钥

## 介绍

基于两个环境变量实现可密钥轮换的加解密

环境变量main_key和main_key2分两种情况

只有main_key,表示当前主密钥

有main_key和main_key2,表示旧主密钥和当前主密钥

主密钥是随机字符组成的字符串，经Argon2id（time=2,memory=64*1024,thread=4,keyLen=32）生成aes密钥。

API

```go
// Encrypt 使用 AES-256-GCM 加密数据
func Encrypt(v string) string

// Decrypt 使用 AES-256-GCM 解密数据
// 如果解密失败，返回空字符串
func Decrypt(v string) string 
```

API会自动处理在不同情况加解密，确保密钥轮换时，新旧密钥可以共存。