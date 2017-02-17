package update

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/inconshreveable/go-update"
	"io"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"io/ioutil"
	"strings"
	"errors"
	"os"
)

func Test_Update(t *testing.T) {
	Convey("Runner test", t, func() {
		am := &ApplierMock{}
		rm := &RemoterMock{}
		r := Runner{
			Applier:    am.Apply,
			Rollbacker: am.RollbackError,
			Remoter:    rm,
		}

		Convey(`Given the current version is  equal to the latest version should NOT run update`, func() {
			sys.Version = "v0.0.2"
			rm.RI = ReleaseInfo{Current: "v0.0.2"}
			err := r.Run()
			So(err, ShouldBeNil)
			So(rm.LatestCalled, ShouldBeFalse)
			So(am.ApplyCalledWith, ShouldBeNil)
		})

		Convey(`Given the current version is not equal to the latest version should run update`, func() {
			sys.Version = "v0.0.1"
			rm.RI = ReleaseInfo{Current: "v0.0.2"}
			err := r.Run()
			So(err, ShouldBeNil)
			So(rm.LatestCalled, ShouldBeTrue)
			So(am.ApplyCalledWith[0], ShouldResemble, ioutil.NopCloser(strings.NewReader("This is the binary")))
			So(am.ApplyCalledWith[1], ShouldResemble, update.Options{})
		})

		Convey(`Given the applier returns an error should rollback`, func() {
			sys.Version = "v0.0.1"
			rm.RI = ReleaseInfo{Current: "v0.0.2"}
			m := "Couldn't apply."
			am.ApplyErr = errors.New(m)
			err := r.Run()
			So(err.Error(), ShouldEqual, m)
			So(am.RollbackCalledWith.Error(), ShouldEqual, m)
		})

		// TODO: Run on ad-hoc basis
		// Convey(`Test remoter info and bin downloader`, func() {
		//	r := S3Remoter{}
		//	ri, err := r.ReleaseInfo()
		//	So(err, ShouldBeNil)
		//	So(ri.Current, ShouldEqual, "v1.2.4")
		//	rdr, err := r.Latest()
		//	So(err, ShouldBeNil)
		//	out, err := os.Create("kwk")
		//	So(err, ShouldBeNil)
		//	io.Copy(out, rdr)
		//})

	})
}

type ApplierMock struct {
	ApplyCalledWith []interface{}
	RollbackCalledWith error
	ApplyErr        error
}

func (am *ApplierMock) Apply(update io.Reader, opts update.Options) error {
	am.ApplyCalledWith = []interface{}{update, opts}
	if am.ApplyErr != nil {
		return am.ApplyErr
	}
	return nil
}

func (am *ApplierMock) RollbackError(err error) error {
	am.RollbackCalledWith = err
	return err
}

type RemoterMock struct {
	RI                ReleaseInfo
	LatestCalled      bool
	ReleaseInfoCalled bool
}

func (rm *RemoterMock) Latest() (io.ReadCloser, error) {
	rm.LatestCalled = true
	r := strings.NewReader("This is the binary")
	return ioutil.NopCloser(r), nil
}

func (rm *RemoterMock) ReleaseInfo() (*ReleaseInfo, error) {
	rm.ReleaseInfoCalled = true
	return &rm.RI, nil
}
