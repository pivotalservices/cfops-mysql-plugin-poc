package main

import (
	"io"
	"strings"

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
	mySqlProduct := pcf.GetProducts()[productName]
	s.PivotalCF = pcf
	s.setIP(mySqlProduct.IPS)
	s.setMysqlCredentials()
	s.setVMCredentials()
	return
}

func (s *MysqlPlugin) Backup() (err error) {
	var writer io.WriteCloser
	var persistanceBackuper cfbackup.PersistanceBackup
	if persistanceBackuper, err = s.GetPersistanceBackup(s.MysqlUserName, s.MysqlPassword, s.getSshConfig()); err == nil {
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
	vmCredentialsName        = "vm_credentials"
	mysqlCredentialsName     = "mysql_admin_password"
	mysqlPrefixName          = "mysql-"
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
	MysqlUserName        string
	MysqlPassword        string
	MysqlIP              string
	VMUserName           string
	VMKey                string
	VMPassword           string
	GetPersistanceBackup func(string, string, command.SshConfig) (cfbackup.PersistanceBackup, error)
}

func (s *MysqlPlugin) getSshConfig() (sshConfig command.SshConfig) {
	sshConfig = command.SshConfig{
		Username: s.VMUserName,
		Password: s.VMPassword,
		Host:     s.MysqlIP,
		Port:     defaultSSHPort,
		SSLKey:   s.VMKey,
	}
	return
}

func (s *MysqlPlugin) setIP(ips map[string][]string) {
	for vmName, ipList := range ips {

		if strings.HasPrefix(vmName, mysqlPrefixName) {
			s.MysqlIP = ipList[0]
		}
	}
}

func (s *MysqlPlugin) setMysqlCredentials() (err error) {
	var props map[string]string
	props, err = s.PivotalCF.GetPropertyValues(productName, jobName, mysqlCredentialsName)
	if err == nil {
		s.MysqlUserName = props[identityName]
		s.MysqlPassword = props[passwordName]
	}
	return
}

func (s *MysqlPlugin) setVMCredentials() (err error) {
	var props map[string]string
	props, err = s.PivotalCF.GetPropertyValues(productName, jobName, vmCredentialsName)
	if err == nil {
		s.VMUserName = props[identityName]
		s.VMPassword = props[passwordName]
	}

	return
}

func newMysqlDumper(user string, pass string, config command.SshConfig) (pb cfbackup.PersistanceBackup, err error) {
	pb, err = persistence.NewRemoteMysqlDump(user, pass, config)
	return
}
