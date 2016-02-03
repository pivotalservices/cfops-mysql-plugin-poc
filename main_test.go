package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
})
