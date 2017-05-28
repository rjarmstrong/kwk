package updater

import (
	"testing"
	"github.com/kwk-super-snippets/types/errs"
	gu "github.com/inconshreveable/go-update"
	"strings"
	"io/ioutil"
	"io"
	"github.com/stretchr/testify/assert"
	"errors"
)


func TestRunner_Run(t *testing.T) {
	am := &ApplierMock{}
	repo := &RepoMock{}
	doc := &DocMock{}
	v := "v0.0.1"
	r := New(v, repo, am.Apply, am.RollbackError, doc)

	t.Log("Given there is an update record with NO error should NOT run update.")
	repo.ReleaseInfo = ReleaseInfo{Version: "v0.0.2"}
	doc.GetErr = nil
	err := r.Run()
	assert.Nil(t, err)
	assert.True(t, repo.LatestInfoCalled)
	assert.False(t, repo.LatestBinCalled)
	assert.False(t, am.ApplyCalled)

	t.Log(`Given the current version is equal to the latest version should NOT run update.`)
	repo.ReleaseInfo = ReleaseInfo{Version: "v0.0.1"}
	doc.GetHydrates = &Record{}
	err = r.Run()
	assert.Nil(t, err)
	assert.False(t, repo.LatestBinCalled)
	assert.False(t, am.ApplyCalled)

	t.Log("Given there is NOT an update record and the remote version is newer should update.")
	repo.ReleaseInfo = ReleaseInfo{Version: "v0.0.2"}
	doc.GetErr = errs.FileNotFound
	err = r.Run()
	assert.Nil(t, err)
	assert.True(t, repo.LatestBinCalled)
	assert.True(t, am.ApplyCalled)

	t.Log(`When updating Given the applier returns an error should rollback.`)
	repo.ReleaseInfo = ReleaseInfo{Version: "v0.0.2"}
	doc.GetErr = errs.FileNotFound
	m := "Couldn't apply."
	am.ApplyErr = errors.New(m)
	err = r.Run()
	assert.IsType(t, am.ApplyErr, err)
	assert.IsType(t, am.ApplyErr, am.RollbackCalledWith)

	// //TODO: Run on ad-hoc basis
	//Convey(`Test remoter info and bin downloader`, func() {
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
	//})
}


type DocMock struct {
	GetErr      error
	GetHydrates interface{}
}

func (*DocMock) Upsert(fullKey string, data interface{}) error {
	return nil
}

func (dm *DocMock) Get(fullKey string, value interface{}, fresherThan int64) error {
	return dm.GetErr
}

func (*DocMock) Delete(fullKey string) error {
	panic("implement me")
}

func (*DocMock) DeleteAll() error {
	panic("implement me")
}

type ApplierMock struct {
	ApplyCalled        bool
	RollbackCalledWith error
	ApplyErr           error
}

func (am *ApplierMock) Apply(update io.Reader, opts gu.Options) error {
	am.ApplyCalled = true
	if am.ApplyErr != nil {
		return am.ApplyErr
	}
	return nil
}

func (am *ApplierMock) RollbackError(err error) error {
	am.RollbackCalledWith = err
	return err
}

type RepoMock struct {
	ReleaseInfo
	LatestBinCalled  bool
	LatestInfoCalled bool
}

func (rm *RepoMock) GetLatestBinary() (io.ReadCloser, error) {
	rm.LatestBinCalled = true
	r := strings.NewReader("This is the binary")
	return ioutil.NopCloser(r), nil
}

func (rm *RepoMock) GetLatestInfo() (*ReleaseInfo, error) {
	rm.LatestInfoCalled = true
	return &rm.ReleaseInfo, nil
}

func (rm *RepoMock) CleanUp() {

}
