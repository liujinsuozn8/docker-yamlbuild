package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func main() {
	// 1. 解析命令行参数
	var yamlPath string
	var imgDir string
	var buildDir string
	var outStd bool

	flag.BoolVar(&outStd, "o", false, "write command to std")
	flag.StringVar(&yamlPath, "y", "", "yaml file path")
	flag.StringVar(&imgDir, "img-dir", "", "image directory")
	flag.StringVar(&buildDir, "build-dir", "", "directory of build image")
	flag.Parse()

	// 2. 读取 yaml 文件的数据
	if yamlPath == ""{
		ErrorFatalLog("-y is empty")
	}
	if imgDir == "" {
		ErrorFatalLog("--img-dir is empty")
	}
	if buildDir == "" {
		ErrorFatalLog("--build-dir is empty")
	}

	yamlData,err := ioutil.ReadFile(yamlPath)
	if err != nil {
		ErrorFatalLog(err)
	}

	// 3. 将路径转换为绝对路径
	imgDir, err = filepath.Abs(imgDir)
	if err != nil {
		ErrorFatalLog(err)
	}

	buildDir, err = filepath.Abs(buildDir)
	if err != nil {
		ErrorFatalLog(err)
	}

	// 4. 创建指令执行对象
	builder := CreateImgBuildCmdBuilder(yamlData, imgDir, buildDir)

	err= builder.Build()

	if err != nil {
		ErrorFatalLog(err)
	}

	// 5. 执行操作
	if outStd {
		// 输出指令到控制台
		for _, cmd := range builder.CmdList {
			fmt.Println(cmd.GetCmdString())
		}
	} else {
		// 编译镜像
		for _, cmd := range builder.CmdList {
			InfoLog(" CMD: ", cmd.GetCmdString())
			err := cmd.BuildImage()
			if err != nil {
				ErrorFatalLog(err)
			}
		}
	}
}
