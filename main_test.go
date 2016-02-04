package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cfbackup"
	. "github.com/pivotalservices/cfops-mysql-plugin-poc"
	"github.com/pivotalservices/cfops/plugin/cfopsplugin"
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
	Describe("given a Setup() method", func() {
		Context("when called with a PivotalCF contain a MySQL tile", func() {
			var pivotalCF cfopsplugin.PivotalCF
			BeforeEach(func() {
				pivotalCF = cfopsplugin.NewPivotalCF(cfbackup.NewConfigurationParser("./fixtures/installation-settings-1-6-aws.json"))
				mysqlplugin.Setup(pivotalCF)
			})

			XIt("then it should extract Mysql username required for backup/restore", func() {
				Ω(mysqlplugin.MysqlUserName).ShouldNot(BeEmpty())
			})
			XIt("then it should extract Mysql password required for backup/restore", func() {
				Ω(mysqlplugin.MysqlPassword).ShouldNot(BeEmpty())
			})
			It("then it should extract Mysql VM IP required for backup/restore", func() {
				Ω(mysqlplugin.MysqlIP).ShouldNot(BeEmpty())
			})
			XIt("then it should extract Mysql VM username required for backup/restore", func() {
				Ω(mysqlplugin.VMUserName).ShouldNot(BeEmpty())
			})
			XIt("then it should extract Mysql VM Key or password required for backup/restore", func() {
				Ω(func() bool {
					return mysqlplugin.VMKey == "" && mysqlplugin.VMPassword == ""
				}()).ShouldNot(BeTrue())

			})
		})
	})
})
