package main

import (
	"fmt"
	cfopsplugin "github.com/pivotalservices/cfops/plugin/cfopsplugin"
)

func main() {
	cfopsplugin.Start(new(MysqlPlugin))
}

const pluginName = "mysql-tile"

func NewMysqlPlugin() *MysqlPlugin {
	return &MysqlPlugin{
		Meta: cfopsplugin.Meta{
			Name: pluginName,
		},
	}
}

type MysqlPlugin struct {
	Meta cfopsplugin.Meta
	MysqlUserName string
	MysqlPassword string
	MysqlIP string
	VMUserName string
	VMKey string
	VMPassword string
}

func (s *MysqlPlugin) GetMeta() (meta cfopsplugin.Meta) {
	meta = s.Meta
	return
}

func (s *MysqlPlugin) Backup() (err error) {
	return
}
func (s *MysqlPlugin) Restore() (err error) {
	return
}

func (s *MysqlPlugin) Setup(pcf cfopsplugin.PivotalCF) (err error) {
	mySqlProduct := pcf.GetProducts()
	fmt.Println(len(mySqlProduct))
	for i, _ := range mySqlProduct{
	fmt.Println(i)
	}
	return
}
