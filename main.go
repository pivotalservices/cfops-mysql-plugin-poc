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

//GetMeta - method to provide metadata
func (s *MysqlPlugin) GetMeta() (meta cfopsplugin.Meta) {
	meta = s.Meta
	return
}

//Setup - on setup method
func (s *MysqlPlugin) Setup(pcf cfopsplugin.PivotalCF) (err error) {
	s.PivotalCF = pcf
	s.InstallationSettings = pcf.GetInstallationSettings()
	return
}

func (s *MysqlPlugin) getSSHConfig() (sshConfig []command.SshConfig, err error) {
	var IPs []string
	var vmCredentials cfbackup.VMCredentials

	if IPs, err = s.InstallationSettings.FindIPsByProductAndJob(productName, jobName); err == nil {
		if vmCredentials, err = s.InstallationSettings.FindVMCredentialsByProductAndJob(productName, jobName); err == nil {
			for _, ip := range IPs {
				sshConfig = append(sshConfig, command.SshConfig{
					Username: vmCredentials.UserID,
					Password: vmCredentials.Password,
					Host:     ip,
					Port:     defaultSSHPort,
					SSLKey:   vmCredentials.SSLKey,
				})
			}
		}
	}
	return
}

//Backup - method to execute backup
func (s *MysqlPlugin) Backup() (err error) {
	var writer io.WriteCloser
	var persistanceBackuper cfbackup.PersistanceBackup
	var mysqlUserName, mysqlPassword string
	var sshConfigs []command.SshConfig

	if sshConfigs, err = s.getSSHConfig(); err == nil {
		//take first node to execute backup on
		sshConfig := sshConfigs[0]
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

	}

	return
}

//Restore - method to execute restore
func (s *MysqlPlugin) Restore() (err error) {
	var reader io.ReadCloser
	var persistanceBackuper cfbackup.PersistanceBackup
	var mysqlUserName, mysqlPassword string

	var sshConfigs []command.SshConfig

	if sshConfigs, err = s.getSSHConfig(); err == nil {
		//take first node to execute restore on
		sshConfig := sshConfigs[0]

		mysqlUserName, mysqlPassword, err = s.getMysqlCredentials()
		if err != nil {
			return
		}
		if persistanceBackuper, err = s.GetPersistanceBackup(mysqlUserName, mysqlPassword, sshConfig); err == nil {
			if reader, err = s.PivotalCF.NewArchiveReader(outputFileName); err == nil {
				defer reader.Close()
				err = persistanceBackuper.Import(reader)
			}
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

//NewMysqlPlugin - Contructor helper
func NewMysqlPlugin() *MysqlPlugin {
	return &MysqlPlugin{
		Meta: cfopsplugin.Meta{
			Name: pluginName,
		},
		GetPersistanceBackup: newMysqlDumper,
	}
}

//MysqlPlugin - structure
type MysqlPlugin struct {
	PivotalCF            cfopsplugin.PivotalCF
	InstallationSettings cfbackup.InstallationSettings
	Meta                 cfopsplugin.Meta
	GetPersistanceBackup func(string, string, command.SshConfig) (cfbackup.PersistanceBackup, error)
}

func (s *MysqlPlugin) getMysqlCredentials() (userName, pwd string, err error) {
	var props map[string]string
	if props, err = s.InstallationSettings.FindPropertyValues(productName, jobName, mysqlCredentialsName); err == nil {
		userName = props[identityName]
		pwd = props[passwordName]
	}
	return
}

func newMysqlDumper(user string, pass string, config command.SshConfig) (pb cfbackup.PersistanceBackup, err error) {
	pb, err = persistence.NewRemoteMysqlDumpWithPath(user, pass, config, mysqlRemoteArchivePath)
	return
}
