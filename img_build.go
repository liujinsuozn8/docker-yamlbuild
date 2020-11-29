package main

import (
	"bufio"
	"fmt"
	"github.com/otiai10/copy"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ImageBuildOptions struct {
	BuildArg       map[string]string `yaml:"build-arg"`
	Id             string            `yaml:"id"`
	Tags           []string          `yaml:"tag"`
	DockerfilePath string            `yaml:"file"`
	BuildDir       string            `yaml:"build-dir"`
}

type ImageBuildCmd struct {
	Dockerfile   *Dockerfile
	BuildOptions *ImageBuildOptions
}

func (this *ImageBuildCmd) GetCmdArgs() []string {
	cmdArgs := []string{"docker","build"}

	// -f
	// 如果不是在 Resource 中编译，则可以附加 -f 选项
	//if !buildInResource {
	//	if this.BuildOptions.DockerfilePath != "" {
	//		cmd.WriteString(" -f " + this.BuildOptions.DockerfilePath)
	//	}
	//}

	// -t
	for _, tag := range this.BuildOptions.Tags {
		cmdArgs = append(cmdArgs, "-t", tag)
	}

	// --build-arg
	for k, v := range this.BuildOptions.BuildArg {
		cmdArgs = append(cmdArgs, "--build-arg", fmt.Sprintf("%s=%s", k, v))

	}

	// 编译的执行位置
	cmdArgs = append(cmdArgs, this.BuildOptions.BuildDir)

	return cmdArgs
}

func  (this *ImageBuildCmd)GetCmdString() string {
	cmdArgs := this.GetCmdArgs()
	return strings.Join(cmdArgs, " ")
}

func (this *ImageBuildCmd) BuildImage() error {
	InfoLog(this.BuildOptions.Id + " is building ...")

	// 1. 获取 dockerfile 所在的目录
	dfDir, dfName := filepath.Split(this.Dockerfile.Path)

	// 2. 如果 dockerfile 所在的目录 与 编译目录相同，则是在 dockerfile 下编译，只创建 ./dockerignore
	// 否则需要将 dockerfile、src 拷贝到 编译目录
	dfCopied := false
	srcCopied := false

	defer func() {
		if dfCopied {
			err := os.Remove(filepath.Join(this.BuildOptions.BuildDir, dfName))
			if err != nil {
				ErrorFatalLog(err)
			}
		}

		if srcCopied {
			err := os.RemoveAll(filepath.Join(this.BuildOptions.BuildDir, "src"))
			if err != nil {
				ErrorFatalLog(err)
			}
		}

		ignorePath := filepath.Join(this.BuildOptions.BuildDir, ".dockeringore")
		if IsFile(ignorePath) {
			err := os.Remove(ignorePath)
			if err != nil {
				ErrorFatalLog(err)
			}
		}
	}()

	if dfDir != this.BuildOptions.BuildDir {
		// 3. 拷贝 dockerfile
		dfCopied = true
		err := copy.Copy(this.Dockerfile.Path, filepath.Join(this.BuildOptions.BuildDir, dfName))
		if err != nil {
			return err
		}

		// 4. 拷贝src
		srcCopied = true
		imgSrc := filepath.Join(dfDir, "src")
		if IsDir(imgSrc) {
			err := copy.Copy(imgSrc, filepath.Join(this.BuildOptions.BuildDir, "src"))
			if err != nil {
				return err
			}
		}
	}

	// 5. 创建 ./dockerignore
	ignore, err := os.OpenFile(filepath.Join(this.BuildOptions.BuildDir, ".dockeringore"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer ignore.Close()
	// 写入忽略内容
	for _, lrp := range this.Dockerfile.LocalResourcePaths {
		_, err := fmt.Fprintln(ignore, "!"+lrp)
		if err != nil {
			return err
		}
	}

	for _, ignoreTarget := range BuildIgnoreList {
		_, err := fmt.Fprintln(ignore, ignoreTarget)
		if err != nil {
			return err
		}
	}

	// 6. 准备指令
	cmdArgs := this.GetCmdArgs()
	command := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	stdoutPipe, err := command.StdoutPipe()
	if err != nil {
		ErrorFatalLog(err)
	}
	defer stdoutPipe.Close()

	stderrPipe, err := command.StderrPipe()
	if err != nil {
		ErrorFatalLog(err)
	}
	defer stderrPipe.Close()

	// 实时输出执行信息
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			fmt.Println(string(scanner.Bytes()))
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			fmt.Println(string(scanner.Bytes()))
		}
	}()

	// 7. 执行指令
	if err := command.Run(); err != nil {
		return err
	}

	return nil
}
