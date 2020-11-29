package main
//
//import (
//	"fmt"
//	"reflect"
//	"testing"
//)
//
//func TestNormalizeID(t *testing.T) {
//	data := `
//- id: id01
//- tag:
//  - xxx
//`
//	builder := CreateImgBuildCmdBuilder("xxx", []byte(data))
//	err := builder.Build()
//	if err != nil {
//		exceptErrMsg := "image index: 2, not set id in option\n"
//		if err.Error() != exceptErrMsg {
//			errMsg := fmt.Sprintf("TestNormalizeID error \nnot except error message.\n actual = %s\n except = %s\n",
//				err.Error(),
//				exceptErrMsg)
//
//			t.Error(errMsg)
//		}
//	} else {
//		t.Error("should be error")
//	}
//}
//
//func TestNormalizeTag(t *testing.T) {
//	data := `
//- id: id01
//- id: id02
//`
//	builder := CreateImgBuildCmdBuilder("xxx", []byte(data))
//	err := builder.Build()
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	optList := builder.GetOptionsList()
//
//	testList := []ImageBuildOptions{{
//		Id:   "id01",
//		Tags: []string{"id01"},
//	}, {
//		Id:   "id02",
//		Tags: []string{"id02"},
//	}}
//
//	if !reflect.DeepEqual(optList, testList) {
//		t.Error("options:", optList)
//		t.Error("testList:", testList)
//	}
//}
