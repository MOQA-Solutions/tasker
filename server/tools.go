package server

import(
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "hash/fnv"
)

func GenerateAPIKey() (string, error) {
    bytes := make([]byte, 32) 
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return "sk_live_" + base64.RawURLEncoding.EncodeToString(bytes), nil
}

func HashAPIKey(key string) string {
    hash := sha256.Sum256([]byte(key))
    return hex.EncodeToString(hash[:])
}

func Phash(s string, n int) int {
    h := fnv.New32a()
    h.Write([]byte(s))
    return int(h.Sum32()) % n
}