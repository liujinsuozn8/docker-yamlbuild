package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var WORD_REGEXP *regexp.Regexp
var KV_STEP *regexp.Regexp
var KV_EQUALITY_STEP *regexp.Regexp
var KV_SPACE_STEP *regexp.Regexp
var KEY_IN_STEP *regexp.Regexp
var PARAM *regexp.Regexp
var COMMENT *regexp.Regexp
var SKIP_PARMAS []string
var BuildIgnoreList []string

func init() {
	SKIP_PARMAS = []string{
		"PATH",
	}

	WORD_REGEXP = regexp.MustCompile(`\b\w+\b`)

	// 1. 由等号、空格分割 key、value
	// sample: ENV KEY=VALUE
	// sample: ENV KEY   VALUE
	// sample: ENV KEY VALUE
	// ^\s* 开头可以是 多个空格
	// ([^ =]+)\s+ 表示指令，指令后跟多个空格
	// ([^ =]+) 表示key
	// (?:\s+|=)，key后面跟多个空格个或1个 =
	// (\S.*) 表示必须以一个非空格的字符开头
	KV_STEP = regexp.MustCompile(`^\s*[^#][^ =]+\s+([^ =]+)(?:\s+|=)(\S.*)$`)

	// 2. 只处理由等号分割的等式形 CMD
	KV_EQUALITY_STEP = regexp.MustCompile(`^\s*[^ =]+\s+([^ =]+)=(\S.*)$`)

	// 3. 只处理由空格分割的 CMD
	KV_SPACE_STEP = regexp.MustCompile(`^\s*[^ =]+\s+([^ ]+)\s+(\S.*)$`)

	// 4. 只获取key
	KEY_IN_STEP = regexp.MustCompile(`^\s*[^ =]+\s+([^ =]+).*$`)

	// 5. 捕获key中的参数
	PARAM = regexp.MustCompile(`\${?([^ ${}:/]+)}?`)

	BuildIgnoreList = []string{"*", "!src/", "!Dockerfile"}

	COMMENT = regexp.MustCompile(`^\s*#.*`)
}

func InfoLog(v ...interface{}) {
	log.Print(" [Info]:", fmt.Sprint(v...))
}
func ErrorFatalLog(v ...interface{}) {
	log.Fatal("[Error]:", fmt.Sprint(v...))
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else {
		return os.IsExist(err)
	}
}

func IsFile(path string) bool {
	fi, e := os.Stat(path)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}

func IsDir(path string) bool {
	fi, e := os.Stat(path)
	if e != nil {
		return false
	}
	return fi.IsDir()
}

func GetFirstWord(text string) string {
	return GetNthWord(text, 1)
}

// 获取字符串中的第一个单词。index从1开始
func GetNthWord(text string, index int) string {
	if index <= 0 {
		panic("index must greater than 0")
	}

	matches := WORD_REGEXP.FindAllString(text, index)

	if len(matches) < index {
		return ""
	} else {
		return matches[index-1]
	}
}

func IsComment(text string) bool {
	return COMMENT.MatchString(text)
}

func ExtractKeyValueFromStep(text string) (string, string, error) {
	matches := KV_STEP.FindStringSubmatch(text)

	if len(matches) < 3 {
		return "", "", errors.New("can not match")
	}

	return matches[1], matches[2], nil
}

func ExtractKeyValueFromEqualityStep(text string) (string, string, error) {
	matches := KV_EQUALITY_STEP.FindStringSubmatch(text)

	if len(matches) < 3 {
		return "", "", errors.New("can not match")
	}

	return matches[1], matches[2], nil
}

func ExtractKeyValueFromSpaceStep(text string) (string, string, error) {
	matches := KV_SPACE_STEP.FindStringSubmatch(text)

	if len(matches) < 3 {
		return "", "", errors.New("can not match")
	}

	return matches[1], matches[2], nil
}

func ExtractKeyFromStep(text string) (string, error) {
	matches := KEY_IN_STEP.FindStringSubmatch(text)

	if len(matches) < 2 {
		return "", errors.New("can not match")
	}

	return matches[1], nil
}

func ExtractParam(text string) [][]string {
	return PARAM.FindAllStringSubmatch(text, -1)
}

func isSkipParam(text string) bool {
	for _, p := range SKIP_PARMAS {
		if text == p {
			return true
		}
	}
	return false
}

func ReplaceParam(text string, envs map[string]string, args map[string]string) (string, error) {
	matches := ExtractParam(text)
	// 防止字符串中的遍历全部都是 isSkipParam，如果没有发生过替换则结束循环

	hasReplace := true

	// 多次循环，防止多层嵌套的情况：${{xxx}}
	for ; len(matches) != 0 && hasReplace; {
		hasReplace = false
		for _, m := range matches {
			if isSkipParam(m[1]) {
				continue
			}
			if value, has := envs[m[1]]; has {
				text = strings.Replace(text, m[0], value, 1)
				hasReplace = true
			} else if value, has := args[m[1]]; has {
				text = strings.Replace(text, m[0], value, 1)
				hasReplace = true
			} else {
				return "", errors.New(fmt.Sprintf("can not find %s in evn or arg", m[0]))
			}
		}

		matches = ExtractParam(text)
	}
	return text, nil
}
