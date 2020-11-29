package main

import (
	"fmt"
	"regexp"
	"testing"
)

func TestIsExist(t *testing.T){
	result := IsExist("./xxx")
	if result {
		t.Error("./xxx not exist")
	}

	result = IsExist("./test")
	if !result {
		t.Error("./test exists")
	}
}

func TestIsFile(t *testing.T) {
	result := IsFile("./xxxx")
	if result {
		t.Error("error, not this file: xxxx")
	}

	result = IsFile("./utils_test.go")
	if !result {
		t.Error("error, utils_test.go is not file")
	}
}

func TestIsDir(t *testing.T) {
	result := IsDir("./xxx")
	if result {
		t.Error("error, not this file: xxxx")
	}

	result = IsDir("./test")
	if !result {
		t.Error("error, test is not dir")
	}
}

func TestGetFirstWord(t *testing.T) {
	text := "  aaa  bbb dddd"
	word := GetFirstWord(text)
	if word != "aaa" {
		t.Errorf("error. actual = %s, except = aaa", word)
	}

	text = " # aaa  bbb dddd"
	word = GetFirstWord(text)
	fmt.Println(word)

}

func TestGetNthWord(t *testing.T) {
	text := "  xx  bbb dddd"

	word := GetNthWord(text, 2)
	if word != "bbb" {
		t.Errorf("error. actual = %s, except = bbb", word)
	}

	word02 := GetNthWord(text, 3)
	if word02 != "dddd" {
		t.Errorf("error. actual = %s, except = dddd", word02)
	}

	word03 := GetNthWord(text, 4)
	if word03 != "" {
		t.Errorf("error. actual = '%s', except = ''", word03)
	}
}

func TestExtractKeyValueFromStep(t *testing.T) {
	text := " env  key    value"
	key, value, err := ExtractKeyValueFromStep(text)

	if err != nil {
		t.Error(err.Error())
	}

	if key != "key" {
		t.Errorf("key: actual = %s, except = key", key)
	}

	if value != "value" {
		t.Errorf("value: actual = %s, except = value", value)
	}

	text = " env  key    "
	key, value, err = ExtractKeyValueFromStep(text)

	if err == nil {
		t.Errorf("error should be error, actual key = %s, actual value = %s", key, value)
	}
}

func TestExtractKeyValueFromEqualityStep(t *testing.T) {
	// case 1
	text := " env  key    value"
	key, value, err := ExtractKeyValueFromEqualityStep(text)

	if err == nil {
		t.Errorf("should be error，but match successful.\n key = %s, value = %s", key, value)
	}

	if "can not match" != err.Error() {
		t.Errorf("erro info is wrong.\n actual = %s, except = can not match", err.Error())
	}

	// case 2
	text = " env  key= value"
	key, value, err = ExtractKeyValueFromEqualityStep(text)

	if err == nil {
		t.Errorf("should be error，but match successful.\n key = %s, value = %s", key, value)
	}

	if "can not match" != err.Error() {
		t.Errorf("erro info is wrong.\n actual = %s, except = can not match", err.Error())
	}

	// case 3
	text = " env  key=value"
	key, value, err = ExtractKeyValueFromEqualityStep(text)

	if err != nil {
		t.Errorf("error.\n" + err.Error())
	}

	if key != "key" {
		t.Errorf("error.\nactual = %s, except = key", key)
	}

	if value != "value" {
		t.Errorf("error.\nactual = %s, except = value", value)
	}
}

func TestExtractKeyValueFromSpaceStep(t *testing.T) {
	// case 1
	text := " env  key=value"
	key, value, err := ExtractKeyValueFromSpaceStep(text)

	if err == nil {
		t.Errorf("should be error，but match successful.\n key = %s, value = %s", key, value)
	}

	if "can not match" != err.Error() {
		t.Errorf("erro info is wrong.\n actual = %s, except = can not match", err.Error())
	}

	// case 2
	text = " env  key= value"
	key, value, err = ExtractKeyValueFromSpaceStep(text)

	if err == nil {
		t.Errorf("should be error，but match successful.\n key = %s, value = %s", key, value)
	}

	if "can not match" != err.Error() {
		t.Errorf("erro info is wrong.\n actual = %s, except = can not match", err.Error())
	}

	// case 3
	text = " env  key value"
	key, value, err = ExtractKeyValueFromSpaceStep(text)

	if err != nil {
		t.Errorf("error.\n" + err.Error())
	}

	if key != "key" {
		t.Errorf("error.\nactual = %s, except = key", key)
	}

	if value != "value" {
		t.Errorf("error.\nactual = %s, except = value", value)
	}
}

func TestExtractKeyFromStep(t *testing.T) {
	text := " env  key    value"
	key, err := ExtractKeyFromStep(text)

	if err != nil {
		t.Error(err.Error())
	}

	if key != "key" {
		t.Errorf("error. key: actual = %s, except = key", key)
	}

	text = " env    "
	key, err = ExtractKeyFromStep(text)

	if err == nil {
		t.Errorf("error should be error, actual key = %s", key)
	}

}

func TestExtractMultiKeyValueFromStep(t *testing.T) {
	KV_EQUALITY := regexp.MustCompile(`^\s*[^ =]+\s+([^ =]+=\S[^ =]+\s+)+`)
	//KV_EQUALITY := regexp.MustCompile(`([^ =]+)=(\S[^ =]+)`)

	//text:="env aaa=bbb ccc=ddd eee=fff"
	text := "env aaa=DDD aa=CCC xxx=www "
	//result := KV_EQUALITY.FindAllStringSubmatch(text, -1)
	result := KV_EQUALITY.FindStringSubmatch(text)

	fmt.Println(len(result))
	for _, r := range result {
		fmt.Println(r)
	}
}

func TestExtractParam(t *testing.T) {
	PARAM = regexp.MustCompile(`\${?([^ ${}:/]+)}`)
	//PARAM = regexp.MustCompile(`\${(.+)}`)
	text := "qqq${aaaa}--${${bbb}}--$$ccc/;--${aaaa}"
	//text := "ddd"
	result := PARAM.FindAllStringSubmatch(text, -1)
	fmt.Println(result)
	fmt.Println(len(result))
}

func TestIsComment(t *testing.T){
	text := "  # aa bb cc"
	if !IsComment(text) {
		t.Error(text + ": should be error")
	}

	text = "   aa bb cc"
	if IsComment(text) {
		t.Error(text + ": not comment")
	}
}