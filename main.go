package main

import (
	"fmt"
	"strings"

	"github.com/pivotalservices/cfbackup"
	cfopsplugin "github.com/pivotalservices/cfops/plugin/cfopsplugin"
)

func main() {
	cfopsplugin.Start(new(MysqlPlugin))
}

const (
	pluginName           = "mysql-tile"
	productName          = "p-mysql"
	jobName              = "mysql"
	vmCredentialsName    = "vm_credentials"
	mysqlCredentialsName = "mysql_admin_password"
	mysqlPrefixName      = "mysql-"
	identityName         = "identity"
	passwordName         = "password"
)

func NewMysqlPlugin() *MysqlPlugin {
	return &MysqlPlugin{
		Meta: cfopsplugin.Meta{
			Name: pluginName,
		},
	}
}

type MysqlPlugin struct {
	Meta          cfopsplugin.Meta
	MysqlUserName string
	MysqlPassword string
	MysqlIP       string
	VMUserName    string
	VMKey         string
	VMPassword    string
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
	fmt.Println(s)
	return
}
