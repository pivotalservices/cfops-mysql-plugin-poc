package main_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cfbackup"
	. "github.com/pivotalservices/cfops-mysql-plugin"
	"github.com/pivotalservices/cfops/plugin/cfopsplugin"
	"github.com/pivotalservices/cfops/tileregistry"
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
	testInstallationSettings("./fixtures/installation-settings-1-6-aws.json")
})

func testInstallationSettings(installationSettingsPath string) {
	var mysqlplugin *MysqlPlugin
	Describe(fmt.Sprintf("given a installationSettingsFile %s", installationSettingsPath), func() {
		Describe("given a Backup() method", func() {
			Context("when called on a properly setup mysqlplugin object", func() {
				var err error
				fakePersistenceBackup := new(FakePersistenceBackup)
				var controlTmpDir string
				BeforeEach(func() {
					controlTmpDir, _ = ioutil.TempDir("", "unit-test")
					mysqlplugin = &MysqlPlugin{
						Meta: cfopsplugin.Meta{
							Name: "mysql-tile",
						},
						GetPersistanceBackup: func(user, pass string, config command.SshConfig) (pb cfbackup.PersistanceBackup, err error) {
							return fakePersistenceBackup, nil
						},
						GetPrivilegeFlusher: func(config command.SshConfig, pwd string) (err error) {
							return
						},
					}
					configParser := cfbackup.NewConfigurationParser(installationSettingsPath)
					pivotalCF := cfopsplugin.NewPivotalCF(configParser.InstallationSettings, tileregistry.TileSpec{
						ArchiveDirectory: controlTmpDir,
					})
					mysqlplugin.Setup(pivotalCF)
					err = mysqlplugin.Backup()
				})

				AfterEach(func() {
					os.RemoveAll(controlTmpDir)
				})

				It("then it should dump the target mysql contents", func() {
					Ω(err).ShouldNot(HaveOccurred())
					Ω(fakePersistenceBackup.DumpCallCount).Should(Equal(1))
				})

				It("then it should create an archive file", func() {
					Ω(err).ShouldNot(HaveOccurred())
					Ω(IsEmpty(controlTmpDir)).ShouldNot(BeTrue())
				})
			})
		})

		Describe("given a Setup() method", func() {
			Context("when called with a PivotalCF containing a MySQL tile", func() {
				var pivotalCF cfopsplugin.PivotalCF
				BeforeEach(func() {
					configParser := cfbackup.NewConfigurationParser(installationSettingsPath)
					pivotalCF = cfopsplugin.NewPivotalCF(configParser.InstallationSettings, tileregistry.TileSpec{})
					mysqlplugin.Setup(pivotalCF)
				})

				It("then it should extract Mysql username required for backup/restore", func() {
					Ω(mysqlplugin.PivotalCF).ShouldNot(BeNil())
				})

			})
		})
	})
}

func IsEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false // Either not empty or error, suits both cases
}
