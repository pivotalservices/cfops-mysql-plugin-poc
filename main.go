package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pivotalservices/cfbackup"
	cfopsplugin "github.com/pivotalservices/cfops/plugin/cfopsplugin"
	"github.com/pivotalservices/gtils/command"
	"github.com/pivotalservices/gtils/persistence"
	"github.com/xchapter7x/lo"
)

var (
	//NewRemoteExecuter -
	NewRemoteExecuter = command.NewRemoteExecutor
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
	lo.G.Debug("Starting setup of mysql-tile")
	s.PivotalCF = pcf
	s.InstallationSettings = pcf.GetInstallationSettings()
	lo.G.Debug("Finished setup of mysql-tile", err)
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
	lo.G.Debug("Starting backup of mysql-tile")
	var writer io.WriteCloser
	var persistanceBackuper cfbackup.PersistanceBackup
	var mysqlUserName, mysqlPassword string
	var sshConfigs []command.SshConfig

	if sshConfigs, err = s.getSSHConfig(); err == nil {
		//take first node to execute backup on
		sshConfig := sshConfigs[0]
		if mysqlUserName, mysqlPassword, err = s.getMysqlCredentials(); err == nil {
			lo.G.Debug("Successfully got mysqlCredentials")
			if persistanceBackuper, err = s.GetPersistanceBackup(mysqlUserName, mysqlPassword, sshConfig); err == nil {
				if writer, err = s.PivotalCF.NewArchiveWriter(outputFileName); err == nil {
					defer writer.Close()
					lo.G.Debug("Starting mysql dump")
					err = persistanceBackuper.Dump(writer)
					lo.G.Debug("Dump finished", err)
				}
			}
		}
	}
	lo.G.Debug("Finished backup of mysql-tile", err)
	return
}

//Restore - method to execute restore
func (s *MysqlPlugin) Restore() (err error) {
	lo.G.Debug("Starting restore of mysql-tile")
	var reader io.ReadCloser
	var persistanceBackuper cfbackup.PersistanceBackup
	var mysqlUserName, mysqlPassword string

	var sshConfigs []command.SshConfig

	if sshConfigs, err = s.getSSHConfig(); err == nil {
		//take first node to execute restore on
		sshConfig := sshConfigs[0]

		if mysqlUserName, mysqlPassword, err = s.getMysqlCredentials(); err == nil {
			if persistanceBackuper, err = s.GetPersistanceBackup(mysqlUserName, mysqlPassword, sshConfig); err == nil {
				if reader, err = s.PivotalCF.NewArchiveReader(outputFileName); err == nil {
					defer reader.Close()
					if err = persistanceBackuper.Import(reader); err == nil {
						err = s.GetPrivilegeFlusher(sshConfig, mysqlPassword)
					}
				}
			}
		}
	}
	lo.G.Debug("Finished restore of mysql-tile", err)
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
		GetPrivilegeFlusher:  flushPrivileges,
	}
}

//MysqlPlugin - structure
type MysqlPlugin struct {
	PivotalCF            cfopsplugin.PivotalCF
	InstallationSettings cfbackup.InstallationSettings
	Meta                 cfopsplugin.Meta
	GetPersistanceBackup func(string, string, command.SshConfig) (cfbackup.PersistanceBackup, error)
	GetPrivilegeFlusher  func(command.SshConfig, string) error
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

func flushPrivileges(sshConfig command.SshConfig, mysqlAdminPwd string) (err error) {
	var remoteExecuter command.Executer
	var writer io.WriteCloser

	if remoteExecuter, err = NewRemoteExecuter(sshConfig); err == nil {
		writer = os.Stdout
		lo.G.Info("flushing priviledges after restore on ip ->", sshConfig.Host)
		var commandToRun = fmt.Sprintf("/var/vcap/packages/mariadb/bin/mysql -u root -h localhost --password=%s -e 'FLUSH PRIVILEGES'", mysqlAdminPwd)
		err = remoteExecuter.Execute(writer, commandToRun)
		lo.G.Info("Done running flush priviledges on ip ->", sshConfig.Host, err)
	}
	return
}
