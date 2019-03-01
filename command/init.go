package command

import (
	"github.com/asdine/storm"
	"github.com/clakeboy/golib/utils"
	"gonotify/asset"
	"gonotify/models"
	"io/ioutil"
	"os"
	"time"
)

func InitSystem() {
	checkGenerate()
	generateConfigFile()
	generateDefault()
	generateAsset()
	Exit("初始化完成!")
}

//检查是否已经初始化生成
func checkGenerate() {
	if utils.PathExists("./db") {
		Exit("已经初始化过,请直接运行程序!")
	}
}

//生成默认配置文件
func generateConfigFile() {
	conf :=
		`
# 系统配置
system:
  port: "13380"
  ip: ""
  pid: notify_server.pid
# boltdb 本地缓存数据库配置
boltdb:
  path: ./db/storm.db
# cookie 配置
cookie:
  path: /
  domain: ""
  source: false
  http_only: false
`

	err := ioutil.WriteFile("./main.conf", []byte(conf), 0665)
	if err != nil {
		Exit("创建配置文件失败: 无法写入配置文件")
	}
}

//生成默认用户
func generateDefault() {
	if !utils.PathExists("./db") {
		err := os.MkdirAll("./db", 0775)
		if err != nil {
			Exit("创建数据失败: 无法创建文件")
		}
	}
	db, err := storm.Open("./db/storm.db")
	defer db.Close()
	if err != nil {
		Exit("创建数据失败: 无法写入文件")
	}
	model := models.NewAccountModel(db)
	account := &models.AccountData{
		Account:     "admin",
		Password:    utils.EncodeMD5("123123"),
		UserName:    "管理员",
		CreatedDate: int(time.Now().Unix()),
	}
	err = model.Save(account)
	if err != nil {
		Exit("创建用户数据失败")
	}
}

//生成静态资源
func generateAsset() {
	asset.RestoreAssets("./", "html")
}
