package processors_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var protoFields protoFieldsStruct

type protoFieldsStruct struct {
	plugin *protogen.Plugin

	// We've hardcoded the file name since we only use one for now
	protoFileName string
}

func (s *protoFieldsStruct) getMessage( //nolint:gocognit
	t *testing.T,
	messageName string,
) *protogen.Message {
	t.Helper()

	var file *protogen.File
	for _, f := range s.plugin.Files {
		if f.Proto.Name != nil && *f.Proto.Name == s.protoFileName {
			file = f
			break
		}
	}
	if file == nil {
		t.Fatalf("proto file not found: %s", s.protoFileName)
	}

	var message *protogen.Message
	for _, m := range file.Messages {
		if string(m.Desc.Name()) == messageName {
			message = m
			break
		}
	}
	if message == nil {
		t.Fatalf("message not found: %s#%s", s.protoFileName, messageName)
	}

	return message
}

func (s *protoFieldsStruct) getField(
	t *testing.T,
	messageName string,
	fieldName string,
) *protogen.Field {
	t.Helper()

	message := s.getMessage(t, messageName)

	var field *protogen.Field
	for _, f := range message.Fields {
		if string(f.Desc.Name()) == fieldName {
			field = f
			break
		}
	}
	if field == nil {
		t.Fatalf("field not found: %s#%s.%s", s.protoFileName, messageName, fieldName)
	}

	return field
}

func TestMain(m *testing.M) {
	if err := initProtoFields(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func initProtoFields() error {
	reqJson, err := ioutil.ReadFile("testdata/request.json")
	if err != nil {
		return err
	}

	var req pluginpb.CodeGeneratorRequest
	if err = json.Unmarshal(reqJson, &req); err != nil {
		return err
	}

	p, err := protogen.Options{}.New(&req)
	if err != nil {
		return err
	}

	protoFields = protoFieldsStruct{
		plugin:        p,
		protoFileName: "example/server.proto",
	}

	return nil
}

func readFixture(t *testing.T, fixtureName string) []byte {
	t.Helper()

	content, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s", fixtureName))
	if err != nil {
		t.Fatalf("error reading fixture %s: %v", fixtureName, err)
	}

	return content
}
