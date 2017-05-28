package app

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
)

const workFolder = "./update_work"

type BinRepo interface {
	// LatestInfo gets information about the latest build.
	LatestInfo() (*ReleaseInfo, error)
	// Latest gets the actual binary.
	LatestBinary() (io.ReadCloser, error)
	// Delete any residual files post update.
	CleanUp()
}

type ReleaseInfo struct {
	Version string `json:"version"`
	Build   string `json:"build"`
	Time    int64  `json:"time"`
	Notes   string `json:"notes"`
}

type S3Repo struct {
}

func (r *S3Repo) LatestInfo() (*ReleaseInfo, error) {
	i, err := http.Get("https://s3.amazonaws.com/kwk-cli/release-info.json")
	if err != nil {
		return nil, err
	}
	ri := &ReleaseInfo{}
	info, err := ioutil.ReadAll(i.Body)
	json.Unmarshal(info, ri)
	if err != nil {
		out.Debug("Failed unmarshalling release info. %v", err)
		return nil, err
	}
	i.Body.Close()
	return ri, nil
}

func (r *S3Repo) LatestBinary() (io.ReadCloser, error) {
	fn := fmt.Sprintf("kwk-%s-%s", runtime.GOOS, runtime.GOARCH)
	fnt := fmt.Sprintf("%s.tar.gz", fn)
	url := fmt.Sprintf("https://s3.amazonaws.com/kwk-cli/latest/bin/%s", fnt)
	err := os.MkdirAll(workFolder, out.StandardFilePermission)
	if err != nil {
		return nil, err
	}
	fnt = path.Join(workFolder, fnt)

	out.Debug("Getting latest from: %s", url)
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
	tarFile := path.Join(workFolder, fn+".tar")
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

func (r *S3Repo) CleanUp() {
	out.Debug("Removing work folder.")
	err := os.RemoveAll(workFolder)
	if err != nil {
		out.LogErrM("Error cleaning up:", err)
	}
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
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}
