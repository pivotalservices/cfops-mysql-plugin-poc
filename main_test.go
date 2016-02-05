package main_test

import (
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cfbackup"
	. "github.com/pivotalservices/cfops-mysql-plugin-poc"
	"github.com/pivotalservices/cfops/plugin/cfopsplugin"
	"github.com/pivotalservices/gtils/command"
)

var _ = Describe("Given MysqlPlugin", func() {
	var mysqlplugin *MysqlPlugin
	Describe("given a Meta() method", func() {
		Context("called on a plugin with valid meta data", func() {
			var meta cfopsplugin.Meta
			BeforeEach(func() {
				mysqlplugin = NewMysqlPlugin()
				meta = mysqlplugin.GetMeta()
			})

			It("then it should return a meta data object with all required fields", func() {
				Ω(meta.Name).ShouldNot(BeEmpty())
			})
		})
	})

	Describe("given a Backup() method", func() {
		Context("when called on a properly setup mysqlplugin object", func() {
			var err error
			backupPath := path.Join(os.TempDir(), "mysql-backup")
			fakePersistenceBackup := new(FakePersistenceBackup)
			BeforeEach(func() {
				mysqlplugin = &MysqlPlugin{
					DestPath: backupPath,
					Meta: cfopsplugin.Meta{
						Name: "mysql-tile",
					},
					GetPersistanceBackup: func(user, pass string, config command.SshConfig) (pb cfbackup.PersistanceBackup, err error) {
						return fakePersistenceBackup, nil
					},
				}
				err = mysqlplugin.Backup()
			})
			It("then it should dump the target mysql contents", func() {
				Ω(fakePersistenceBackup.DumpCallCount).Should(Equal(1))
			})
		})
	})

	Describe("given a Setup() method", func() {
		Context("when called with a PivotalCF containing a MySQL tile", func() {
			var pivotalCF cfopsplugin.PivotalCF
			BeforeEach(func() {
				pivotalCF = cfopsplugin.NewPivotalCF(cfbackup.NewConfigurationParser("./fixtures/installation-settings-1-6-aws.json"))
				mysqlplugin.Setup(pivotalCF)
			})

			It("then it should extract Mysql username required for backup/restore", func() {
				Ω(mysqlplugin.MysqlUserName).ShouldNot(BeEmpty())
			})
			It("then it should extract Mysql password required for backup/restore", func() {
				Ω(mysqlplugin.MysqlPassword).ShouldNot(BeEmpty())
			})
			It("then it should extract Mysql VM IP required for backup/restore", func() {
				Ω(mysqlplugin.MysqlIP).ShouldNot(BeEmpty())
			})
			It("then it should extract Mysql VM username required for backup/restore", func() {
				Ω(mysqlplugin.VMUserName).ShouldNot(BeEmpty())
			})
			It("then it should extract Mysql VM Key or password required for backup/restore", func() {
				Ω(func() bool {
					return mysqlplugin.VMKey == "" && mysqlplugin.VMPassword == ""
				}()).ShouldNot(BeTrue())

			})

		})
	})

})
