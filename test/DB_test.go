package service_test

import (
	"fmt"
	"testing"

	"github.com/solenovex/web-tutor/model"
	"github.com/solenovex/web-tutor/service"
)

func TestFindByTitle(t *testing.T) {
	i := service.FindNewByTitle("河南")
	fmt.Println(i)
}
func TestFindByGid(t *testing.T) {
	i := service.CheckGid("qwer")
	fmt.Println(i)
}
func TestInsert(t *testing.T) {
	data := model.News{Gid: "qwer", Title: "check", Text: "test", Date: "20220125"}
	i := service.InsertNews(data)
	fmt.Println(i)
}
