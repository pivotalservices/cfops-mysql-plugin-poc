package main

import (
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

func (s *MysqlPlugin) Setup(cfopsplugin.PivotalCF) (err error) {
	return
}
