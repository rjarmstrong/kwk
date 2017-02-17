package update

import (
	"io"
	"fmt"
	"runtime"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"os"
)


type Remoter interface {
	Latest() (io.ReadCloser, error)
	ReleaseInfo() (*ReleaseInfo, error)
}

type S3Remoter struct {
}

func (r *S3Remoter) Latest() (io.ReadCloser, error) {
	fn := fmt.Sprintf("kwk-%s-%s", runtime.GOOS, runtime.GOARCH)
	fnt := fmt.Sprintf("%s.tar.gz", fn)
	url := fmt.Sprintf("https://s3.amazonaws.com/kwk-cli/latest/bin/%s", fnt)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	out, err := os.Create(fnt)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	out.Close()
	exe(true,"tar", "-xvf", fnt) //TODO: tar prob no supported everywhere
	return os.Open(fn)
}


func (r *S3Remoter) ReleaseInfo() (*ReleaseInfo, error) {
	i, err := http.Get("https://s3.amazonaws.com/kwk-cli/release-info.json")
	if err != nil {
		return nil, err
	}
	ri := &ReleaseInfo{}
	info, err := ioutil.ReadAll(i.Body)
	json.Unmarshal(info, ri)
	if err != nil {
		return nil, err
	}
	i.Body.Close()
	return ri, nil
}