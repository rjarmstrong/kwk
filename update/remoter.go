package update

import (
	"io"
	"fmt"
	"runtime"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"bitbucket.com/sharingmachine/kwkcli/log"
	"bitbucket.com/sharingmachine/kwkcli/cache"
	"compress/gzip"
	"archive/tar"
)

const workFolder = "./update_work"


type Remoter interface {
	Latest() (io.ReadCloser, error)
	CleanUp()
	ReleaseInfo() (*ReleaseInfo, error)
}

type S3Remoter struct {
}

func (r *S3Remoter) Latest() (io.ReadCloser, error) {
	fn := fmt.Sprintf("kwk-%s-%s", runtime.GOOS, runtime.GOARCH)
	fnt := fmt.Sprintf("%s.tar.gz", fn)
	url := fmt.Sprintf("https://s3.amazonaws.com/kwk-cli/latest/bin/%s", fnt)
	err := os.MkdirAll(workFolder, cache.StandardFilePermission)
	if err != nil {
		return nil, err
	}
	fnt = path.Join(workFolder, fnt)

	log.Debug("Getting latest from: %s", url)
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
	tarFile := path.Join(workFolder, fn + ".tar")
	target := path.Join(workFolder, fn)
	err = unGZip(fnt, tarFile)
	if err != nil {
		return nil, err
	}
	err = unTar(tarFile, target)
	if err != nil {
		return nil, err
	}
	return os.Open(target)
}

func unGZip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return err
}

func unTar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		info := header.FileInfo()
		if info.IsDir() {
			continue
		}

		file, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *S3Remoter) CleanUp() {
	log.Debug("Removing work folder.")
	err := os.RemoveAll(workFolder)
	if err != nil {
		log.Error("Error cleaning up.", err)
	}
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