package main

import (
	"io"

	"github.com/pivotalservices/cfbackup"
	cfopsplugin "github.com/pivotalservices/cfops/plugin/cfopsplugin"
	"github.com/pivotalservices/gtils/command"
	"github.com/pivotalservices/gtils/persistence"
	"github.com/xchapter7x/lo"
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
	var reader io.ReadCloser
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

	lo.G.Info("restoring to %s using %s and %s", sshConfig.Host, sshConfig.Username, sshConfig.Password)
	if persistanceBackuper, err = s.GetPersistanceBackup(mysqlUserName, mysqlPassword, sshConfig); err == nil {
		if reader, err = s.PivotalCF.NewArchiveReader(outputFileName); err == nil {
			defer reader.Close()
			err = persistanceBackuper.Import(reader)
		}
	}
	return
}

const (
	pluginName                 = "mysql-tile"
	outputFileName             = pluginName + ".dmp"
	productName                = "p-mysql"
	jobName                    = "mysql"
	mysqlCredentialsName       = "mysql_admin_password"
	identityName               = "identity"
	passwordName               = "password"
	defaultSSHPort         int = 22
	mysqlRemoteArchivePath     = "/var/vcap/store/mysql/archive.backup"
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
	pb, err = persistence.NewRemoteMysqlDumpWithPath(user, pass, config,mysqlRemoteArchivePath)
	return
}
