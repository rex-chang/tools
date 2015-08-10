/**
 * 本地文件和CDN上的文件对比,检查是否更新成功
 * @author    rex chang
 */
package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/issue9/term/colors"
)

var (
	remotePath  string
	rootPath    string
	checkEndMsg = "All is ok"
)

func main() {
	remoteUrl := flag.String("t", "", "远程对比文件的地址,如[http://simple.bxds.com/ss/v1/],\n\r\t 默认为resource下[version.manifest]中的[packageUrl]")
	resource := flag.String("s", "./resource", "本地路径,默认为[./resource]")
	flag.Parse()

	_, err := os.Stat(*resource)
	if err != nil { //判断文件目录是否合法
		// fmt.Print(*resource + "文件目录不存在!请检查输入")
		printErr(err.Error())
		return
	}

	rootPath = *resource
	remotePath = *remoteUrl
	if *remoteUrl == "" {
		_, err := os.Stat("./resource/version.manifest")

		if err != nil {
			// printErr(err.Error())
			printErr(err.Error())
		}
		fi, err := os.Open("./resource/version.manifest")
		defer fi.Close()
		data, err := ioutil.ReadAll(fi)

		jsonData, err := simplejson.NewJson(data)
		if err != nil {
			printErr(err.Error())
		}
		ds, err := jsonData.Get("packageUrl").String()
		if err != nil {
			printErr(err.Error())
		}
		remotePath = ds
	}
	fmt.Println("start")
	fmt.Print("Remote Url:")
	colors.Printf(colors.Stdout, colors.Cyan, colors.Black, "%s\n", remotePath)

	visitFile(*resource)
	fmt.Println(checkEndMsg)
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// stdIn := os.Stdin
	// for {
	// 	c, _ := f.Read(stdIn)
	// 	ioutil.R
	// }

	// var command string
	// for {

	// 	fmt.Scan(&command)
	// 	print(command)
	// 	if command == nil {
	// 		break
	// 	}

	// }
}

func md5File(fp string) string {
	fp = filepath.Join(rootPath, fp)
	fi, err := os.Stat(fp)
	if err != nil || fi.IsDir() {
		return ""
	}
	fc, err := os.Open(fp)
	if err != nil {
		panic(err)
	} else {
		defer fc.Close()
	}
	defer fc.Close()
	md5Ctx := md5.New()
	fcb, _ := ioutil.ReadAll(fc)
	md5Ctx.Write(fcb)
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

func md5RemoteFile(remoteFile string) string {
	rtF, err := http.Get(remoteFile)

	if err != nil {
		return ""
	} else {
		defer rtF.Body.Close()
	}

	md5Ctx := md5.New()

	body, err := ioutil.ReadAll(rtF.Body)
	md5Ctx.Write(body)
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

/**
 *  遍历文件
 * @param  {[type]} path string)       (fp string [description]
 * @return {[type]}      [description]
 */
func visitFile(path string) (fp string) {
	fileInfo, err := os.Stat(path)
	if err != nil { //判断文件目录是否合法
		printErr(err.Error())
		return
	}
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		fileDir, _ := filepath.Rel(rootPath, path)

		if info.IsDir() {
			return nil
		}
		// colors.Print(colors.Stdout, colors.Red, colors.Blue, "colors")
		uri := strings.Replace(fileDir, "\\", "/", -1)

		fgColor := colors.Green
		msg := "[success]"
		if md5File(fileDir) != md5RemoteFile(remotePath+uri) {
			fgColor = colors.Red
			msg = "[failed]"
			checkEndMsg = "\ncheckFailed, Please waiting for a moment\n"
		}
		colors.Printf(colors.Stdout, colors.White, colors.Black, "%-60s", remotePath+uri)
		colors.Printf(colors.Stdout, fgColor, colors.Black, "\t%s\n", msg)
		return nil
	})

	return path + "\\" + fileInfo.Name()
}

func printErr(msg string) {
	colors.Print(colors.Stdout, colors.Red, colors.Black, msg+"\n")
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(1)
}
