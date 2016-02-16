package main

import (
	"io"

	"github.com/pivotalservices/cfbackup"
	cfopsplugin "github.com/pivotalservices/cfops/plugin/cfopsplugin"
	"github.com/pivotalservices/gtils/command"
	"github.com/pivotalservices/gtils/persistence"
)

func main() {
	cfopsplugin.Start(NewMysqlPlugin())
}

func (s *MysqlPlugin) GetMeta() (meta cfopsplugin.Meta) {
	meta = s.Meta
	return
}

func (s *MysqlPlugin) Setup(pcf cfopsplugin.PivotalCF) (err error) {
	s.PivotalCF = pcf
	return
}

func (s *MysqlPlugin) Backup() (err error) {
	var writer io.WriteCloser
	var persistanceBackuper cfbackup.PersistanceBackup
	var sshConfig command.SshConfig
	var mysqlUserName, mysqlPassword string

	sshConfig, err = s.PivotalCF.GetSSHConfig(productName, jobName)
	if err != nil {
		return
	}

	mysqlUserName, mysqlPassword, err = s.getMysqlCredentials()
	if err != nil {
		return
	}
	if persistanceBackuper, err = s.GetPersistanceBackup(mysqlUserName, mysqlPassword, sshConfig); err == nil {
		if writer, err = s.PivotalCF.NewArchiveWriter(outputFileName); err == nil {
			defer writer.Close()
			err = persistanceBackuper.Dump(writer)
		}
	}
	return
}

func (s *MysqlPlugin) Restore() (err error) {
	return
}

const (
	pluginName               = "mysql-tile"
	outputFileName           = pluginName + ".dmp"
	productName              = "p-mysql"
	jobName                  = "mysql"
	mysqlCredentialsName     = "mysql_admin_password"
	identityName             = "identity"
	passwordName             = "password"
	defaultSSHPort       int = 22
)

func NewMysqlPlugin() *MysqlPlugin {
	return &MysqlPlugin{
		Meta: cfopsplugin.Meta{
			Name: pluginName,
		},
		GetPersistanceBackup: newMysqlDumper,
	}
}

type MysqlPlugin struct {
	PivotalCF            cfopsplugin.PivotalCF
	Meta                 cfopsplugin.Meta
	GetPersistanceBackup func(string, string, command.SshConfig) (cfbackup.PersistanceBackup, error)
}

func (s *MysqlPlugin) getMysqlCredentials() (userName, pwd string, err error) {
	var props map[string]string
	props, err = s.PivotalCF.GetPropertyValues(productName, jobName, mysqlCredentialsName)
	if err == nil {
		userName = props[identityName]
		pwd = props[passwordName]
	}
	return
}

func newMysqlDumper(user string, pass string, config command.SshConfig) (pb cfbackup.PersistanceBackup, err error) {
	pb, err = persistence.NewRemoteMysqlDump(user, pass, config)
	return
}
