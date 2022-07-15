package tools

import (
	"strings"
)

func SetURL(k, v, t string) string {
	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	endpoint := "http://127.0.0.1:8080/update/"
	return endpoint + t + "/" + k + "/" + v
}

func GetURL(url string) []string {
	//	Парсим URL
	urlPath := strings.Split(strings.TrimLeft(url, "update/"), "/")
	return urlPath
}
