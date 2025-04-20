package rpc

import (
	"encoding/json"
	"fmt"
	"errors"
	"strconv"
	"bytes"
)

func EncodeMessage(msg any) (string, error) {
	content, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content), nil
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r','\n','\r','\n'})
	if !found {
		return "", nil, errors.New("Couldn't find seperator between header and content")
	}

	contentLengthBytes := header[len([]byte("Content-Length: ")):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, err
	}

	var baseMsg BaseMessage
	if err = json.Unmarshal(content[:contentLength], &baseMsg); err != nil {
		return "", nil, err
	}

	return baseMsg.Method, content[:contentLength], nil
}

func Split(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r','\n','\r','\n'})
	if !found {
		return 0, nil, nil
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, err
	}
	if len(content) < contentLength {
		return 0, nil, nil
	}

	totalLength := len(header) + 4 + contentLength // +4 because of the \r\n\r\n
	return totalLength, data[:totalLength], nil
}
