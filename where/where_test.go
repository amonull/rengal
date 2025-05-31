package where

import (
	"testing"

	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/amonull/rengal/filesystem"
)

func init() {
	filesystem.SetMemMapFs()
}

func TestConfig(t *testing.T) {
	Convey("When gettings config path", t, func() {
		path := Config()
		Convey("It should exist", func() {
			exists := lo.Must(filesystem.Api().Exists(path))
			So(exists, ShouldBeTrue)

			Convey("And it should be a directory", func() {
				isDir := lo.Must(filesystem.Api().IsDir(path))
				So(isDir, ShouldBeTrue)
			})
		})
	})
}

func TestSources(t *testing.T) {
	Convey("When gettings sources path", t, func() {
		path := Sources()
		Convey("It should exist", func() {
			exists := lo.Must(filesystem.Api().Exists(path))
			So(exists, ShouldBeTrue)

			Convey("And it should be a directory", func() {
				isDir := lo.Must(filesystem.Api().IsDir(path))
				So(isDir, ShouldBeTrue)
			})
		})
	})
}

func TestLogs(t *testing.T) {
	Convey("When gettings logs path", t, func() {
		path := Logs()
		Convey("It should exist", func() {
			exists := lo.Must(filesystem.Api().Exists(path))
			So(exists, ShouldBeTrue)

			Convey("And it should be a directory", func() {
				isDir := lo.Must(filesystem.Api().IsDir(path))
				So(isDir, ShouldBeTrue)
			})
		})
	})
}
