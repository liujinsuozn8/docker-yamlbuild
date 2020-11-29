package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"path/filepath"
)

type ImgBuildCmdBuilder struct {
	ImgDir   string
	BuildDir string
	YamlData []byte
	CmdList  []ImageBuildCmd
}

func CreateImgBuildCmdBuilder(yamlData []byte, imgDir string, buildDir string) *ImgBuildCmdBuilder {

	return &ImgBuildCmdBuilder{
		ImgDir:   imgDir,
		YamlData: yamlData,
		BuildDir: buildDir,
	}
}

func (this *ImgBuildCmdBuilder) Build() error {
	// 解析 yaml 文件数据
	optionsList, err := this.parseDataToOptionList()
	if err != nil {
		return err
	}

	// 创建 ImgBuildCmd
	for i := 0; i < len(optionsList); i++ {
		InfoLog(fmt.Sprintf("image index: %d, image id: %s, analyzing...", i+1, optionsList[i].Id))
		// check id
		if optionsList[i].Id == "" {
			errMsg := fmt.Sprintf("image index: %d, not set id in option\n", i+1)
			err := errors.New(errMsg)
			return err
		}

		// opt标准化
		err = this.normalizeOptions(i, &optionsList[i])
		if err != nil {
			return err
		}

		// 读取dockerfile
		df, err := CreateDockerfile(
			filepath.Join(this.ImgDir, optionsList[i].Id, "Dockerfile"),
			optionsList[i].BuildArg)
		if err != nil {
			return err
		}

		this.CmdList = append(this.CmdList, ImageBuildCmd{
			Dockerfile:   df,
			BuildOptions: &optionsList[i],
		})
	}

	InfoLog("yaml analysis has ended")
	return nil
}

// 将 byte 数据转换成 ImgBuildOptions 对象数组
func (this *ImgBuildCmdBuilder) parseDataToOptionList() ([]ImageBuildOptions, error) {
	var optionsList []ImageBuildOptions

	err := yaml.Unmarshal(this.YamlData, &optionsList)
	if err != nil {
		return nil, err
	}

	return optionsList, nil
}

// opt 标准化
func (this *ImgBuildCmdBuilder) normalizeOptions(index int, opt *ImageBuildOptions) error {
	err := this.normalizeTag(index, opt)
	if err != nil {
		return err
	}

	err = this.normalizeBuildDir(index, opt)
	if err != nil {
		return err
	}

	return nil
}

func (this *ImgBuildCmdBuilder) normalizeTag(index int, opt *ImageBuildOptions) error {
	// set default value for tag
	if len(opt.Tags) == 0 {
		opt.Tags = []string{opt.Id}
	}

	return nil
}

func (this *ImgBuildCmdBuilder) normalizeBuildDir(index int, opt *ImageBuildOptions) error {
	// 如果 BuildPath = ""，则使用指令中的 --build-path参数
	// 如果 BuildPath !=""，则使用yaml中的参数
	if opt.BuildDir == "" {
		opt.BuildDir = this.BuildDir
	} else {
		// 检查yml指定的目录是否存在
		if !IsDir(opt.BuildDir) {
			return errors.New(fmt.Sprintf("error image option index: %d.\n %s is not such directory", index, opt.BuildDir))
		}
	}
	return nil
}
