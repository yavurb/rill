package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

// Encode encodes the input in base64
// It can optionally zip the input before encoding
func Encode(obj interface{}, compress bool) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	if compress {
		b = zip(b)
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

// Decode decodes the input from base64
// It can optionally unzip the input after decoding
func Decode(in string, obj interface{}, compress bool) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		fmt.Println("Error decoding base64")
		return err
	}

	if compress {
		fmt.Println("Unzipping")
		b = unzip(b)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		fmt.Println("Error unmarshalling")
		return err
	}

	return nil
}

func zip(in []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(in)
	if err != nil {
		panic(err)
	}
	err = gz.Flush()
	if err != nil {
		panic(err)
	}
	err = gz.Close()
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}

func unzip(in []byte) []byte {
	var b bytes.Buffer
	_, err := b.Write(in)
	if err != nil {
		panic(err)
	}
	r, err := gzip.NewReader(&b)
	if err != nil {
		panic(err)
	}
	res, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return res
}
