package serializer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
)

func WriteProtobufToBinaryFile(
	message proto.Message, fileName string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal message to binary: %w", err)
	}

	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write binary to file: %w", err)
	}
	return nil
}

func ReadProtobufFromBinaryFile(
	filename string, message proto.Message) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("cannot unmarshal binary to message: %w", err)
	}
	return nil
}

func WriteProtobufToJSONFile(
	message proto.Message, fileName string) error {
	data, err := ProtobufToJSON(message)
	if err != nil {
		return fmt.Errorf("cannot marshal message to JSON: %w", err)
	}

	err = ioutil.WriteFile(fileName, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("cannot write binary to file: %w", err)
	}
	return nil
}

func ReadProtobufFromJSONFile(
	filename string, message proto.Message) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	err = json.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("cannot unmarshal binary to message: %w", err)
	}
	return nil
}
