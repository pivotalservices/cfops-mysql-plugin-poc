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
	cfopsplugin.Start(new(MysqlPlugin))
}

const (
	pluginName               = "mysql-tile"
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
	}
}

type MysqlPlugin struct {
	Meta              cfopsplugin.Meta
	MysqlUserName     string
	MysqlPassword     string
	MysqlIP           string
	VMUserName        string
	VMKey             string
	VMPassword        string
	PersistanceBackup cfbackup.PersistanceBackup
}

func (s *MysqlPlugin) GetMeta() (meta cfopsplugin.Meta) {
	meta = s.Meta
	return
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

func (s *MysqlPlugin) setPersistanceBackup() (err error) {
	s.PersistanceBackup, err = persistence.NewRemoteMysqlDump(s.MysqlUserName, s.MysqlPassword, s.getSshConfig())
	return
}

func (s *MysqlPlugin) Backup() (err error) {
	//TODO complete this method but need a io.writer
	var writer io.Writer
	s.PersistanceBackup.Dump(writer)
	return
}
func (s *MysqlPlugin) Restore() (err error) {
	return
}

func (s *MysqlPlugin) setIP(ips map[string][]string) {
	for vmName, ipList := range ips {
		if strings.HasPrefix(vmName, mysqlPrefixName) {
			s.MysqlIP = ipList[0]
		}
	}
}

func (s *MysqlPlugin) getMysqlProperties(jobsList []cfbackup.Jobs) (mysqlProperties []cfbackup.Properties) {
	for _, job := range jobsList {
		if job.Identifier == jobName {
			mysqlProperties = job.Properties
		}
	}
	return
}
func (s *MysqlPlugin) setMysqlCredentials(jobsList []cfbackup.Jobs) {
	mysqlProperties := s.getMysqlProperties(jobsList)

	for _, property := range mysqlProperties {
		if property.Identifier == mysqlCredentialsName {

			s.MysqlUserName = property.Value.(map[string]interface{})[identityName].(string)
			s.MysqlPassword = property.Value.(map[string]interface{})[passwordName].(string)

		}

	}

}

func (s *MysqlPlugin) setVMCredentials(jobsList []cfbackup.Jobs) {
	mysqlProperties := s.getMysqlProperties(jobsList)

	for _, property := range mysqlProperties {
		if property.Identifier == vmCredentialsName {

			s.VMUserName = property.Value.(map[string]interface{})[identityName].(string)
			s.VMPassword = property.Value.(map[string]interface{})[passwordName].(string)

		}

	}

}
func (s *MysqlPlugin) Setup(pcf cfopsplugin.PivotalCF) (err error) {
	mySqlProduct := pcf.GetProducts()[productName]
	s.setIP(mySqlProduct.IPS)
	s.setMysqlCredentials(mySqlProduct.Jobs)
	s.setVMCredentials(mySqlProduct.Jobs)
	s.setPersistanceBackup()
	return
}
