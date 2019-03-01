package models

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/clakeboy/golib/utils"
	"gonotify/common"
	"log"
	"testing"
	"time"
)

func init() {
	var err error
	common.BDB, err = storm.Open("../db/storm.db")

	if err != nil {
		fmt.Println("open database error:", err)
	}
}

func TestAccountModel_GetByAccount(t *testing.T) {
	model := NewAccountModel(nil)
	data := &AccountData{
		Account:     "Clake",
		Password:    utils.EncodeMD5("123123"),
		UserName:    "管理员",
		Disable:     false,
		CreatedDate: int(time.Now().Unix()),
	}
	err := model.Save(data)
	if err != nil {
		fmt.Println(err)
	}
}

func TestNewAccountModel(t *testing.T) {
	model := NewAccountModel(nil)

	update := &AccountData{
		Id:          1,
		Password:    utils.EncodeMD5("123123"),
		CreatedDate: int(time.Now().Unix()),
	}

	err := model.Update(update)
	fmt.Println(err)
	//err := model.Save(update)
	//fmt.Println(err)

}

func TestNewAccountModel2(t *testing.T) {
	model := NewAccountModel(nil)
	//model.Init(new(AccountData))
	data, err := model.GetById(1)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(data)
}

func TestNewAccountModel3(t *testing.T) {
	model := NewTaskLogModel(common.BDB)
	model.Drop(new(TaskLogData))
}

func TestAccountModel_List(t *testing.T) {
	cipher := utils.NewAes("123456")
	out, err := cipher.DecryptString("QbUle6LWmDsTeL5iLlXE4A==")
	if err != nil {
		fmt.Println("decrypt error", err)
		return
	}

	fmt.Println(out)

	out, err = cipher.EncryptString("cipher text for php")
	if err != nil {
		fmt.Println("encrypt error", err)
		return
	}

	fmt.Println(out)

	fmt.Println(utils.RandStr(16, nil))
}
