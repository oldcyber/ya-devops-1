package tools

func SetURL(k, v, t string) string {
	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	endpoint := "http://127.0.0.1:8080/update/"
	return endpoint + t + "/" + k + "/" + v
}
