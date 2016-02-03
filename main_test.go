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
	XDescribe("given a Setup() method", func() {
		Context("given a PivotalCF contain a MySQL tile", func() {
			var pivotalCF cfopsplugin.PivotalCF
			BeforeEach(func() {
				pivotalCF = cfopsplugin.NewPivotalCF(cfbackup.NewConfigurationParser("./fixtures/installation-settings-1-6-aws.json"))
				mysqlplugin.Setup(pivotalCF)
			})

			It("then it should extract information about the MySQL deployment required  for backup/restore", func() {
				Ω(mysqlplugin.MysqlUserName).ShouldNot(BeEmpty())
				Ω(mysqlplugin.MysqlPassword).ShouldNot(BeEmpty())
				Ω(mysqlplugin.MysqlIP).ShouldNot(BeEmpty())
				Ω(mysqlplugin.VMUserName).ShouldNot(BeEmpty())
				Ω(func() bool {
					return mysqlplugin.VMKey == "" && mysqlplugin.VMPassword == ""
				}()).ShouldNot(BeTrue())

			})
		})
	})
})
