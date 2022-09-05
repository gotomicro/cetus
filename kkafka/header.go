package kkafka

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

func HeadersPrint(headers []kafka.Header) {
	fmt.Printf("headers print begin \n")
	for _, header := range headers {
		fmt.Printf("%s: %s \n", header.Key, string(header.Value))
	}
	fmt.Printf("headers print finish \n")
	return
}

func Headers2Map(headers []kafka.Header) map[string][]byte {
	res := make(map[string][]byte, 0)
	for _, header := range headers {
		res[header.Key] = header.Value
	}
	return res
}

func HeadersValue(headers []kafka.Header, key string) []byte {
	for _, header := range headers {
		if header.Key == key {
			return header.Value
		}
	}
	return nil
}

func HeadersAdd(headers []kafka.Header, key string, val []byte) []kafka.Header {
	isSet := false
	for index, header := range headers {
		if header.Key == key {
			headers[index].Value = val
			isSet = true
			break
		}
	}
	if !isSet {
		headers = append(headers, kafka.Header{
			Key:   key,
			Value: val,
		})
	}
	return headers
}

func HeadersBatchAdd(headers []kafka.Header, needAddHeaders []kafka.Header) []kafka.Header {
	newHeadersMap := Headers2Map(needAddHeaders)
	for _, h := range headers {
		if _, ok := newHeadersMap[h.Key]; !ok {
			needAddHeaders = append(needAddHeaders, h)
		}
	}
	return needAddHeaders
}
