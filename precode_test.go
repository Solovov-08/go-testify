package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cafeList = map[string][]string{
	"moscow": {"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

// Тест при случае: Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое. Если в параметре count указано больше, чем есть всего кафе, должны вернуться все доступные кафе.
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {

	totalCount := 4               // ожидаемое количество
	testCount := strconv.Itoa(20) // кол-во для теста

	req, err := http.NewRequest("GET", "http://localhost/?count="+testCount+"&city=moscow", nil) // здесь нужно создать запрос к сервису
	if err != nil {
		log.Println("Error creating request:", err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code, "status code is not 200") // проверка кода ответа

	cafesByte, err := io.ReadAll(responseRecorder.Body) // читаем тело ответа
	if err != nil {
		log.Println("Error reading response:", err)
	}

	stringerCafes := string(cafesByte)         // конвертируем байты в строку
	cafes := strings.Split(stringerCafes, ",") // дробим строку на слайс

	assert.Equal(t, totalCount, len(cafes), "incorrect number of cities in response") // проверяем длину слайса и проверочное число
}

// Тест при случае: Город, который передаётся в параметре city, не поддерживается. Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
func TestUnsupportedCityHandling(t *testing.T) {

	respBody := "wrong city value" // ожидаемое тело ответа
	testCity := "lipetsk"          // тестовый город

	req, err := http.NewRequest("GET", "http://localhost/?count=4&city="+testCity, nil) // здесь нужно создать запрос к сервису
	if err != nil {
		log.Println("Error creating request:", err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "status code is not 400") // проверка кода ответа

	errorByte, err := io.ReadAll(responseRecorder.Body) // читаем тело ответа
	if err != nil {
		log.Println("Error reading response:", err)
	}

	strigerError := string(errorByte) // конвертируем в строку

	assert.Equal(t, respBody, strigerError, "response body is not as expected") // проверка тела ответа с ожидаемой строкой
}
