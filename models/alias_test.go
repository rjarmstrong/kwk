package models

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Alias(t *testing.T) {

	Convey("Alias", t, func() {

		Convey(`If distinctName is preceeded with a slash should be an absolute alias`, func() {
			fk := "/richard/fun.txt"
			a, err := ParseAlias(fk)
			So(err, ShouldBeEmpty)
			So(a.Ext, ShouldEqual, "txt")
			So(a.Name, ShouldEqual, "fun")
			So(a.Pouch, ShouldBeEmpty)
			So(a.Username, ShouldEqual, "richard")

			fk = "/richard/dong/fun.txt"
			a, err = ParseAlias(fk)
			So(err, ShouldBeEmpty)
			So(a.Ext, ShouldEqual, "txt")
			So(a.Name, ShouldEqual, "fun")
			So(a.Pouch, ShouldEqual, "dong")
			So(a.Username, ShouldEqual, "richard")
		})

		Convey(`Given an absolute distinctName which has only one segment should fail`, func() {
			fk := "/richard"
			_, err := ParseAlias(fk)
			So(err, ShouldNotBeNil)
			So(err, ShouldHaveSameTypeAs, ClientErr{})
		})

		Convey(`Given an empty last segment should return an error`, func() {
			fk := "richard/examples/"
			_, err := ParseAlias(fk)
			So(err, ShouldNotBeNil)
			So(err, ShouldHaveSameTypeAs, ClientErr{})
		})


		Convey(`If distinctName is preceeded WITHOUT a slash should be an absolute alias`, func() {
			fk := "richard/fun.txt"
			a, err := ParseAlias(fk)
			So(err, ShouldBeEmpty)
			So(a.Ext, ShouldEqual, "txt")
			So(a.Name, ShouldEqual, "fun")
			So(a.Pouch, ShouldEqual, "richard")
			So(a.Username, ShouldEqual, "")

			fk = "richard/dong/fun.txt"
			a, err = ParseAlias(fk)
			So(err, ShouldBeEmpty)
			So(a.Ext, ShouldEqual, "txt")
			So(a.Name, ShouldEqual, "fun")
			So(a.Pouch, ShouldEqual, "dong")
			So(a.Username, ShouldEqual, "richard")
		})

		Convey(`If distinctName is empty should have empty values`, func() {
			fk := ""
			a, err := ParseAlias(fk)
			So(err, ShouldBeEmpty)
			So(a.Ext, ShouldBeEmpty)
			So(a.Name, ShouldBeEmpty)
			So(a.Pouch, ShouldBeEmpty)
		})

		Convey(`Should parse distinctName with extension to Alias`, func() {
			fk := "myscript.js"
			rk, err := ParseAlias(fk)
			So(err, ShouldBeEmpty)
			So(rk.Ext, ShouldEqual, "js")
			So(rk.Name, ShouldEqual, "myscript")
		})

		Convey(`Should parse distinctName without extension to Alias`, func() {
			fk := "myscript"
			rk, err := ParseAlias(fk)
			So(err, ShouldBeEmpty)
			So(rk.Ext, ShouldBeEmpty)
			So(rk.Name, ShouldEqual, "myscript")
		})

		Convey(`Should parse distinctName with multiple dots to extension to Alias`, func() {
			fk := "myscript.test.js"
			rk, err := ParseAlias(fk)
			So(err, ShouldBeEmpty)
			So(rk.Ext, ShouldEqual, "js")
			So(rk.Name, ShouldEqual, "myscript.test")
			So(rk.SnipName.String(), ShouldEqual, fk)
		})

		Convey(`Given a prefix of . should ignore it`, func() {
			fk := ".mykey"
			a, err := ParseAlias(fk)
			So(err, ShouldBeNil)
			So(a.Name, ShouldEqual, "mykey")
		})

		Convey(`Given a prefix which has more than 3 segments should fail`, func() {
			fk := "ding/dong/bing/bong.txt"
			_, err := ParseAlias(fk)
			So(err, ShouldNotBeNil)
			So(err, ShouldHaveSameTypeAs, ClientErr{})
		})

		Convey(`Given a list of distinctNames should parse them and return a single pouch`, func() {
			dn := []string{"ding/help/apple.txt", "ding/help/pear.txt", "ding/help/banana.txt"}
			sns, pouch, err := ParseMany(dn)
			So(err, ShouldBeNil)
			So(pouch, ShouldEqual, "help")
			So(sns[0].Name, ShouldResemble, "apple")
			So(sns[0].Ext, ShouldResemble, "txt")
			So(sns[1].Name, ShouldResemble, "pear")
			So(sns[2].Name, ShouldResemble, "banana")
		})

		Convey(`Given a list of distinctNames if one is invalid should return error`, func() {
			dn := []string{"ding/help/apple.txt", "ding/help/", "ding/help/banana.txt"}
			sns, pouch, err := ParseMany(dn)
			So(err, ShouldNotBeNil)
			So(pouch, ShouldBeEmpty)
			So(sns, ShouldBeNil)
		})

		Convey(`Given a list of distinctNames if there are mulitple pouches detected should return error`, func() {
			dn := []string{"ding/settings/apple.txt", "ding/help/pear", "ding/help/banana.txt"}
			sns, pouch, err := ParseMany(dn)
			So(err, ShouldNotBeNil)
			So(err.(ClientErr).Msgs[0].Code, ShouldEqual, Code_MultiplePouches)
			So(pouch, ShouldBeEmpty)
			So(sns, ShouldBeNil)
		})

		Convey(`Given a 3 segment alias should print as absolute alias`, func() {
			fk := "richard/examples/zog.sh"
			a, err := ParseAlias(fk)
			So(err, ShouldBeNil)
			So(a.String(), ShouldEqual, "/richard/examples/zog.sh")
		})

		Convey(`Given a 2 segment absolute alias should print as absolute alias without a pouch`, func() {
			fk := "/richard/zog.sh"
			a, err := ParseAlias(fk)
			So(err, ShouldBeNil)
			So(a.String(), ShouldEqual, "/richard/zog.sh")
		})

		Convey(`Given a 2 segment relative alias should print as relative alias without a pouch`, func() {
			fk := "richard/zog.sh"
			a, err := ParseAlias(fk)
			So(err, ShouldBeNil)
			So(a.String(), ShouldEqual, "richard/zog.sh")
		})

		Convey(`Given a 1 segment relative alias should print as relative alias without a pouch`, func() {
			fk := "zog.sh"
			a, err := ParseAlias(fk)
			So(err, ShouldBeNil)
			So(a.String(), ShouldEqual, "zog.sh")
		})

	})
}
