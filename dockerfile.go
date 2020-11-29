package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	CMD_ENV  = "ENV"
	CMD_ARG  = "ARG"
	CMD_COPY = "COPY"
	CMD_ADD  = "ADD"
)

type Dockerfile struct {
	Envs               map[string]string
	ArgMap             map[string]string
	ArgKeys            []string
	Path               string
	DFSteps            []*DFStep
	LocalResourcePaths []string
}

type DFStep struct {
	StepType   string
	Components []string
	StepStr    string
	RowNos     []int
}

func CreateDockerfile(path string, buildArgs map[string]string) (*Dockerfile, error) {
	df := Dockerfile{
		Path:               path,
		Envs:               make(map[string]string),
		ArgMap:             buildArgs,
		ArgKeys:            []string{},
		DFSteps:            []*DFStep{},
		LocalResourcePaths: []string{},
	}

	// 打开 Dockerfile
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = f.Close(); err != nil {
			ErrorFatalLog(err)
		}
	}()

	// 逐行读取 Dockerfile，并解析
	s := bufio.NewScanner(f)
	rowNo := 0
	for s.Scan() {
		rowNo = rowNo + 1
		err := analyzeStep(&df, rowNo, s.Text())
		if err != nil {
			return nil, err
		}
	}

	err = s.Err()
	if err != nil {
		return nil, err
	}

	return &df, err
}

func analyzeStep(df *Dockerfile, rowNo int, text string) error {
	// 获取指令类型
	stepType := GetFirstWord(text)

	// 解析指令
	if stepType == CMD_ENV {
		dfStep, err := parseENVStep(rowNo, text)
		if err != nil {
			return newDFError(err.Error(), df.Path, rowNo, text)
		}

		err = tryAddEnv(df, dfStep)
		if err != nil {
			return newDFError(err.Error(), df.Path, rowNo, text)
		}
	} else if stepType == CMD_ARG {
		dfStep, err := parseARGStep(rowNo, text)
		if err != nil {
			return newDFError(err.Error(), df.Path, rowNo, text)
		}

		err = tryAddArg(df, dfStep)
		if err != nil {
			return newDFError(err.Error(), df.Path, rowNo, text)
		}
	} else if stepType == CMD_ADD {
		dfStep, err := parseADDStep(rowNo, text)
		if err != nil {
			return newDFError(err.Error(), df.Path, rowNo, text)
		}
		df.DFSteps = append(df.DFSteps, dfStep)

		localRscPath, err := ReplaceParam(dfStep.Components[0], df.Envs, df.ArgMap)
		if err != nil {
			return newDFError(err.Error(), df.Path, rowNo, text)
		}

		// check path
		inDfdirRscPath := filepath.Join(filepath.Dir(df.Path), localRscPath)
		inResourceRscPath := filepath.Join(filepath.Dir(df.Path), localRscPath)
		if IsExist(inDfdirRscPath) || IsExist(inResourceRscPath) {
			df.LocalResourcePaths = append(df.LocalResourcePaths, localRscPath)
		} else {
			return newDFError(fmt.Sprintf("can not find:\n\t%s\nor\n\t%s\n", inDfdirRscPath, inResourceRscPath), df.Path, rowNo, text)
		}

	} else if stepType == CMD_COPY {
		dfStep, err := parseCOPYStep(rowNo, text)
		if err != nil {
			return err
		}
		df.DFSteps = append(df.DFSteps, dfStep)

		localRscPath, err := ReplaceParam(dfStep.Components[0], df.Envs, df.ArgMap)
		if err != nil {
			return newDFError(err.Error(), df.Path, rowNo, text)
		}

		// check path
		inDfdirRscPath := filepath.Join(filepath.Dir(df.Path), localRscPath)
		inResourceRscPath := filepath.Join(filepath.Dir(df.Path), localRscPath)
		if IsExist(inDfdirRscPath) || IsExist(inResourceRscPath) {
			df.LocalResourcePaths = append(df.LocalResourcePaths, localRscPath)
		} else {
			return newDFError(fmt.Sprintf("can not find:\n\t%s\nor\n\t%s\n", inDfdirRscPath, inResourceRscPath), df.Path, rowNo, text)
		}
	} else {
		df.DFSteps = append(df.DFSteps, &DFStep{
			StepType:   stepType,
			Components: []string{},
			StepStr:    text,
			RowNos:     []int{rowNo},
		})
	}

	return nil
}

func newDFError(baseErrorMsg, dfPath string, rowNo int, lineText string) error {
	return errors.New(fmt.Sprintf("[%s] has error\nerror row: \n[%d]:%s\n\nerror message:\n%s\n", dfPath, rowNo, lineText, baseErrorMsg))
}

func tryAddEnv(df *Dockerfile, dfStep *DFStep) error {
	// TODO multi env
	df.DFSteps = append(df.DFSteps, dfStep)

	// try add ENV
	value, err := ReplaceParam(dfStep.Components[1], df.Envs, df.ArgMap)
	if err != nil {
		return err
	}

	key, err := ReplaceParam(dfStep.Components[0], df.Envs, df.ArgMap)
	if err != nil {
		return err
	}

	df.Envs[key] = value

	return nil
}

func tryAddArg(df *Dockerfile, dfStep *DFStep) error {
	df.DFSteps = append(df.DFSteps, dfStep)

	// try add ARG

	key, err := ReplaceParam(dfStep.Components[0], df.Envs, df.ArgMap)
	if err != nil {
		return err
	}

	if len(dfStep.Components) == 2 {
		value, err := ReplaceParam(dfStep.Components[1], df.Envs, df.ArgMap)
		if err != nil {
			return err
		}

		df.ArgMap[key] = value
	} else {
		// 如果 option 中已经存在 ARG，则忽略; 否则添加
		if _, has := df.ArgMap[dfStep.Components[0]]; !has {
			df.ArgKeys = append(df.ArgKeys, dfStep.Components[0])
		}
	}

	return nil
}

func parseENVStep(rowNo int, text string) (*DFStep, error) {
	key, value, err := ExtractKeyValueFromStep(text)

	if err != nil {
		return nil, err
	}

	dfCmd := DFStep{
		StepType:   CMD_ENV,
		Components: []string{key, value},
		StepStr:    text,
		RowNos:     []int{rowNo},
	}

	return &dfCmd, nil
}

func parseARGStep(rowNo int, text string) (*DFStep, error) {
	//尝试从 key=value 的格式中获取 key、value
	key, value, err := ExtractKeyValueFromEqualityStep(text)

	if err == nil {
		return &DFStep{
			StepType:   CMD_ARG,
			Components: []string{key, value},
			StepStr:    text,
			RowNos:     []int{rowNo},
		}, nil
	}

	// 尝试只获取 key
	key, err = ExtractKeyFromStep(text)
	if err != nil {
		return nil, err
	}

	return &DFStep{
		StepType:   CMD_ARG,
		Components: []string{key},
		StepStr:    text,
		RowNos:     []int{rowNo},
	}, nil
}

func parseADDStep(rowNo int, text string) (*DFStep, error) {
	key, value, err := ExtractKeyValueFromSpaceStep(text)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s\nerror row: \n[%d]:%s\n", err.Error(), rowNo, text))
	}

	dfCmd := DFStep{
		StepType:   CMD_ADD,
		Components: []string{key, value},
		StepStr:    text,
		RowNos:     []int{rowNo},
	}

	return &dfCmd, nil
}

func parseCOPYStep(rowNo int, text string) (*DFStep, error) {
	// TODO multiBuild:  COPY --from=xxxx path1 path2
	key, value, err := ExtractKeyValueFromSpaceStep(text)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s\nerror row: \n[%d]:%s\n", err.Error(), rowNo, text))
	}

	dfCmd := DFStep{
		StepType:   CMD_COPY,
		Components: []string{key, value},
		StepStr:    text,
		RowNos:     []int{rowNo},
	}

	return &dfCmd, nil
}
