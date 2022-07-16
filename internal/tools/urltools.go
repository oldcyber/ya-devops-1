package tools

import (
	"strings"
)

func SetURL(k, v, t string) string {
	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	endpoint := "http://127.0.0.1:8080/update/"
	return endpoint + t + "/" + k + "/" + v
}

func GetURL(url string, h string) []string {
	//	Парсим URL
	if h == "update" {
		urlPath := strings.Split(strings.TrimLeft(url, "update/"), "/")
		if len(urlPath) < 3 {
			return nil
		}
		return urlPath
	} else if h == "value" {
		urlPath := strings.Split(strings.TrimLeft(url, "value/"), "/")
		if len(urlPath) < 2 {
			return nil
		}
		return urlPath
	}
	return nil
}
