package update

import (
	"errors"
	gu "github.com/inconshreveable/go-update"
	"github.com/kwk-super-snippets/cli/src/models"
	"github.com/kwk-super-snippets/cli/src/persist"
	"github.com/kwk-super-snippets/types/errs"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func Test_Update(t *testing.T) {
	Convey("Runner test", t, func() {
		am := &ApplierMock{}
		rm := &RemoterMock{}
		pm := &persist.PersisterMock{}
		models.SetPrefs(models.DefaultPrefs())

		r := Updater{
			Applier:        am.Apply,
			Rollbacker:     am.RollbackError,
			BinRepo:        rm,
			Persister:      pm,
			currentVersion: "v0.0.1",
		}

		Convey("Given there is an update record should NOT run update.", func() {
			rm.BI = ReleaseInfo{Version: "v0.0.2"}
			pm.GetHydrates = &Record{}
			err := r.Run()
			So(err, ShouldBeNil)
			So(rm.LatestCalled, ShouldBeFalse)
			So(am.ApplyCalledWith, ShouldBeNil)
		})

		Convey(`Given the current version is equal to the latest version should NOT run update.`, func() {
			rm.BI = ReleaseInfo{Version: "v0.0.1"}
			pm.GetHydrates = &Record{}
			err := r.Run()
			So(err, ShouldBeNil)
			So(rm.LatestCalled, ShouldBeFalse)
			So(am.ApplyCalledWith, ShouldBeNil)
		})

		Convey("Given there is NOT an update record and the remote version is newer should update.", func() {
			rm.BI = ReleaseInfo{Version: "v0.0.2"}
			pm.GetReturns = errs.FileNotFound
			err := r.Run()
			So(err, ShouldBeNil)
			So(rm.LatestCalled, ShouldBeTrue)
			So(am.ApplyCalledWith, ShouldNotBeNil)
		})

		Convey(`When updating Given the applier returns an error should rollback.`, func() {
			rm.BI = ReleaseInfo{Version: "v0.0.2"}
			pm.GetReturns = errs.FileNotFound
			m := "Couldn't apply."
			am.ApplyErr = errors.New(m)
			err := r.Run()
			So(err.Error(), ShouldEqual, m)
			So(am.RollbackCalledWith.Error(), ShouldEqual, m)
		})

		// //TODO: Run on ad-hoc basis
		Convey(`Test remoter info and bin downloader`, func() {
			//r := S3Remoter{}
			//ri, err := r.LatestInfo()
			//So(err, ShouldBeNil)
			//So(ri.Version, ShouldEqual, "1.2.3")
			//So(ri.Build, ShouldEqual, "12")
			//So(ri.Time, ShouldEqual, 233423423)
			//So(ri.Notes, ShouldResemble, "Feature A\nFeature B\n")
			//rdr, err := r.LatestBinary()
			//So(err, ShouldBeNil)
			//out, err := os.Create("kwk")
			//So(err, ShouldBeNil)
			//io.Copy(out, rdr)
		})

	})
}

type ApplierMock struct {
	ApplyCalledWith    []interface{}
	RollbackCalledWith error
	ApplyErr           error
}

func (am *ApplierMock) Apply(update io.Reader, opts gu.Options) error {
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
	BI                ReleaseInfo
	LatestCalled      bool
	ReleaseInfoCalled bool
}

func (rm *RemoterMock) LatestBinary() (io.ReadCloser, error) {
	rm.LatestCalled = true
	r := strings.NewReader("This is the binary")
	return ioutil.NopCloser(r), nil
}

func (rm *RemoterMock) LatestInfo() (*ReleaseInfo, error) {
	rm.ReleaseInfoCalled = true
	return &rm.BI, nil
}

func (rm *RemoterMock) CleanUp() {

}
