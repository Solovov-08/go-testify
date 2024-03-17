package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тест при случае: Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое.
func TestCorrectRequestReturns200Status(t *testing.T) {

	testCount := strconv.Itoa(4)
	testCity := "moscow"

	req, err := http.NewRequest("GET", "http://localhost/?count="+testCount+"&city="+testCity, nil)
	require.NoError(t, err, "failed to create request")

	responseRecoder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecoder, req)

	assert.Equal(t, http.StatusOK, responseRecoder.Code, "status code is not 200")

}

// Тест при случае: Если в параметре count указано больше, чем есть всего кафе, должны вернуться все доступные кафе.
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {

	totalCount := 4               // ожидаемое количество
	testCount := strconv.Itoa(20) // кол-во для теста

	req, err := http.NewRequest("GET", "http://localhost/?count="+testCount+"&city=moscow", nil) // здесь нужно создать запрос к сервису
	require.NoError(t, err, "failed to create request")

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	cafesByte, err := io.ReadAll(responseRecorder.Body) // читаем тело ответа
	if err != nil {
		t.Fatalf("Error reading response: %v", err)
	}

	stringerCafes := string(cafesByte)         // конвертируем байты в строку
	cafes := strings.Split(stringerCafes, ",") // дробим строку на слайс

	assert.Equal(t, totalCount, len(cafes), "incorrect number of cities in response") // проверяем длину слайса и проверочное число
}

// Тест при случае: Город, который передаётся в параметре city, не поддерживается. Сервис возвращает код ответа 400 и ошибку 'wrong city value' в теле ответа.
func TestUnsupportedCityHandling(t *testing.T) {

	respBody := "wrong city value" // ожидаемое тело ответа
	testCity := "lipetsk"          // тестовый город

	req, err := http.NewRequest("GET", "http://localhost/?count=4&city="+testCity, nil) // здесь нужно создать запрос к сервису
	require.NoError(t, err, "failed to create request")

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "status code is not 400") // проверка кода ответа

	errorByte, err := io.ReadAll(responseRecorder.Body) // читаем тело ответа
	if err != nil {
		t.Fatalf("Error reading response: %v", err)
	}

	strigerError := string(errorByte) // конвертируем в строку

	assert.Equal(t, respBody, strigerError, "response body is not as expected") // проверка тела ответа с ожидаемой строкой
}
