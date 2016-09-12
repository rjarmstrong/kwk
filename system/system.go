package system

import (
	"os/exec"
	"fmt"
	"bytes"
	"io"
	"encoding/json"
	"github.com/kennygrant/sanitize"
	"io/ioutil"
	"path"
	"os"
)

func ExecSafe(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	var stderr bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
}

func PrettyPrint(obj interface{}) {
	fmt.Println("")
	p, _ := json.MarshalIndent(obj, "", "  ")
	fmt.Print(string(p))
	fmt.Print("\n\n")
}

func UpsertDirectory(dir string) error {
	ok, err := Exists(dir);
	if ok {
		return nil
	}
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}
	return nil
}

func WriteToFile(directoryName string, fullKey string, uri string) (string, error) {
	dirPath, err := GetDirPath(directoryName)
	filePath := path.Join(dirPath, sanitize.Name(fullKey))
	err = ioutil.WriteFile(filePath, []byte(uri), 0666)
	return filePath, err
}

func ReadFromFile(directoryName string, fullKey string) (string, error) {
	dirPath, err := GetDirPath(directoryName)
	if err != nil { return "", err }
	fp := path.Join(dirPath, fullKey)
	bts, err := ioutil.ReadFile(fp)
	return string(bts), err
}

func GetDirPath(directoryName string) (string, error) {
	dir := path.Join(GetCachePath(), directoryName)
	err:= UpsertDirectory(dir)
	return dir, err
}

func GetCachePath() string {
	p := fmt.Sprintf("/Users/%s/Library/Caches/kwk", os.Getenv("USER"))
	if err := os.Mkdir(p, os.ModeDir); err != nil {
		if os.IsExist(err) {
			return p
		}
		panic(err)
	}
	return p
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}

func Upgrade(){
	distributionUri := "/Volumes/development/go/src/github.com/kwk-links/kwk-cli/kwk-cli"
	installPath := "/usr/local/bin/kwk"
	//download in future
	CopyFile(distributionUri, installPath)
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	fmt.Println("      Successfully upgraded!")
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}