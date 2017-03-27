package models

import (
	"testing"
	"models"
)

func TestConn(t *testing.T) {
	defer models.Conn.Close()

	ret := models.Publishs().JoinPubsByIds([]int{22956})
	t.Log(ret)
}
