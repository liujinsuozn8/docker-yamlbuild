package main
//
//import (
//	"fmt"
//	"reflect"
//	"testing"
//)
//
//func TestCreateDockerfileOpen(t *testing.T) {
//	path := "./test/img/test-img/xxxx"
//
//	_, err := CreateDockerfile(path)
//	if err == nil {
//		t.Errorf("TestCreateDockerfileOpen error: should be error, but opened")
//	}
//}
//
//func TestDockerfileEVNNotValue(t *testing.T) {
//	path := "./test/img/df-env-not-value/Dockerfile"
//
//	df, err := CreateDockerfile(path)
//	if err == nil {
//		t.Error("TestDockerfileEVNNotValue not error !!!. DockerFile struct = ", df)
//	} else {
//		rowNo := 4
//		text := "ENV KAFKA_HOME"
//		exceptErrMsg := fmt.Sprintf("can not match\n"+
//			"error row: \n"+
//			"[%d]:%s\n",
//			rowNo,
//			text)
//
//		if exceptErrMsg != err.Error() {
//			errMsg := fmt.Sprintf("TestDockerfileEVNNotValue error: not except error message\nactual = %s\nexcept = %s\n",
//				err.Error(),
//				exceptErrMsg)
//
//			t.Error(errMsg)
//		}
//	}
//}
//
//func TestDockerfileEVNNotKeyValue(t *testing.T) {
//	path := "./test/img/df-env-not-key-value/Dockerfile"
//
//	df, err := CreateDockerfile(path)
//	if err == nil {
//		t.Error("TestDockerfileEVNNotKeyValue not error !!!. DockerFile struct = ", df)
//	} else {
//		rowNo := 4
//		text := "ENV"
//		exceptErrMsg := fmt.Sprintf("can not match\n"+
//			"error row: \n"+
//			"[%d]:%s\n",
//			rowNo,
//			text)
//
//		if exceptErrMsg != err.Error() {
//			errMsg := fmt.Sprintf("TestDockerfileEVNNotValue error: not except error message\nactual = %s\nexcept = %s\n",
//				err.Error(),
//				exceptErrMsg)
//
//			t.Error(errMsg)
//		}
//	}
//}
//
//func TestDockerfileEVNNormal(t *testing.T) {
//	path := "./test/img/df-env-normal/Dockerfile"
//
//	df, err := CreateDockerfile(path)
//	if err != nil {
//		t.Error("TestDockerfileEVNNormal error: ", err.Error())
//		return
//	}
//
//	except := Dockerfile{
//		Path: path,
//		Envs: map[string]string{
//			"KAFKA_HOME": "/opt/module/kafka",
//			"VERSION":    "12345",
//			"PATH":       "$PATH:$KAFKA_HOME/bin",
//		},
//		ArgMap: make(map[string]string),
//		KeyArgs:     []string{},
//		DFSteps: []*DFStep{
//			&DFStep{StepType: "ENV",
//				Components: []string{"KAFKA_HOME", "/opt/module/kafka"},
//				StepStr:     "ENV KAFKA_HOME /opt/module/kafka",
//				RowNos:     []int{3},
//			},
//			&DFStep{StepType: "ENV",
//				Components: []string{"VERSION", "12345"},
//				StepStr:     "ENV VERSION=12345",
//				RowNos:     []int{4},
//			},
//			&DFStep{StepType: "ENV",
//				Components: []string{"PATH", "$PATH:$KAFKA_HOME/bin"},
//				StepStr:     "ENV PATH $PATH:$KAFKA_HOME/bin",
//				RowNos:     []int{5},
//			},
//		},
//	}
//
//	if !reflect.DeepEqual(df, &except) {
//		t.Error("TestDockerfileEVNNotValue compare error:\nactual = \n")
//		t.Error(df)
//		t.Error("except = \n")
//		t.Error(&except)
//	}
//
//}
