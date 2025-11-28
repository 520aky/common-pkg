package wuid

import "testing"

func TestGenUid(t *testing.T) {
	dsn := "admin:123456@tcp(127.0.0.1:3306)/lotterygame?charset=utf8mb4&parseTime=true"

	uid := GenUid(dsn)

	t.Log(uid)
}
