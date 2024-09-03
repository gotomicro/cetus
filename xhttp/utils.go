package xhttp

import (
	"net/url"
	"regexp"
)

const (
	// ProtocolRegexp HTTP 协议正则
	ProtocolRegexp = "http(s)?://"
	// DefaultPortStr HTTP 默认端口
	DefaultPortStr = "80"
)

func GetURLHost(sourceURL string) (string, error) {
	endpointHostPort := sourceURL
	containsProtocol, err := regexp.MatchString(ProtocolRegexp, sourceURL)
	if err != nil {
		return "", err
	}
	if containsProtocol {
		parsedURL, err := url.Parse(sourceURL)
		if err != nil {
			return "", err
		}
		// parsedURL.Host will be host or host:port
		endpointHostPort = parsedURL.Host
	}
	return endpointHostPort, nil
}
