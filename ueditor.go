package controllers

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

type UeditorController struct {
	beego.Controller
}
type FileItem struct {
	Url string `json:"url"` // 文件url
}
var filenames=make([]FileItem,0)

func Listdir(path,def string)[]FileItem{
	fileInfoList,err := ioutil.ReadDir(path+def)
	//fmt.Println("fileInfoList",len(fileInfoList))
	if err != nil {
		fmt.Println("err1",err)
	}
	var  name FileItem
	for _,fileInfo:=range fileInfoList{
		if fileInfo.IsDir(){
			//fmt.Println("IsDir",path+fileInfo.Name())
			Listdir(path,fileInfo.Name())
			continue
		}else {
			//fmt.Println("def",def+fileInfo.Name())
			if strings.Split(fileInfo.Name(),".")[1] != "mp4"{
				name.Url=def+"/"+fileInfo.Name()
				filenames=append(filenames,name)
			}

		}

	}
	return filenames

}

func (u *UeditorController)Action()  {
	action := u.GetString("action")
	//fmt.Println(action)

	data:=make(map[string]interface{})

	switch action {
	//自动读入配置文件，只要初始化UEditor即会发生
	case "config":

		jsonByte, _ := ioutil.ReadFile("static/ueditor/conf/config.json")
		re, _ := regexp.Compile("\\/\\*[\\S\\s]+?\\*\\/")
		jsonByte = re.ReplaceAll(jsonByte, []byte(""))
		u.Ctx.WriteString(string(jsonByte))

	case "uploadimage":
		{
			f,h,err:= u.GetFile("upfile")
			defer f.Close()

			rand.Seed(time.Now().UnixNano())

			randNum := fmt.Sprintf("%d", rand.Intn(9999)+1000)

			hashName := md5.Sum([]byte( time.Now().Format("200601_02_15_04_05") + randNum ))

			filename := fmt.Sprintf("%x", hashName) + h.Filename

			uploadDir := time.Now().Format("20060102/")

			err = os.MkdirAll(beego.AppConfig.String("Path")+uploadDir, 777)
			//fmt.Println("mkdir",err)

			url:= beego.AppConfig.String("Path") + uploadDir + filename
			//fmt.Println("url",url)

			err=u.SaveToFile("upfile",url)
			//fmt.Println("err",err)

			if err != nil {
				data["state"]="FAIL"
				data["url"]=url
				data["title"]=filename
				data["original"]=filename
				u.Data["json"]=data
				u.ServeJSON()

			} else {
				data["state"]="SUCCESS"
				data["url"]=url
				data["title"]=filename
				data["original"]=filename
				u.Data["json"]=data
				//fmt.Println(u.Data["json"])
				u.ServeJSON()
			}

		}

	case "uploadvideo":
		{
			f,h,err:= u.GetFile("upfile")
			defer f.Close()

			rand.Seed(time.Now().UnixNano())

			randNum := fmt.Sprintf("%d", rand.Intn(9999)+1000)

			hashName := md5.Sum([]byte( time.Now().Format("200601_02_15_04_05") + randNum ))

			filename := fmt.Sprintf("%x", hashName) + h.Filename

			uploadDir := time.Now().Format("20060102/")

			err = os.MkdirAll(beego.AppConfig.String("Path")+uploadDir, 777)
			//fmt.Println("mkdir",err)

			url:= beego.AppConfig.String("Path") + uploadDir + filename
			//fmt.Println("url",url)

			err=u.SaveToFile("upfile",url)
			//fmt.Println("err",err)

			if err != nil {
				data["state"]="FAIL"
				data["url"]=url
				data["title"]=filename
				data["original"]=filename
				u.Data["json"]=data
				u.ServeJSON()

			} else {
				data["state"]="SUCCESS"
				data["url"]=url
				data["title"]=filename
				data["original"]=filename
				u.Data["json"]=data
				u.ServeJSON()
			}

		}

	case "listimage":
		{
			
			//获取文件或目录相关信息
			filenames:=Listdir(beego.AppConfig.String("Path"),"")
			//fmt.Println(filenames)
			data["state"]="SUCCESS"
			data["list"]=filenames
			data["total"]=len(filenames)
			data["start"]=0
			u.Data["json"]=data
			u.ServeJSON()
		}

	}

}
