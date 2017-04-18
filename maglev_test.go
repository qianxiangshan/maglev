package maglev

import (
	"testing"
)

func TestInit(t *testing.T) {

	var mg Maglev

	var sernames = []string{"1", "2", "3"}

	t.Log("start\n")
	mg.Init(sernames)
	t.Log("stop\n")

	t.Log(mg.Get("1"))
	t.Log(mg.Get("2"))
	t.Log(mg.Get("3"))
	t.Log(mg.Get("5"))
	t.Log(mg.Get("10"))

}
