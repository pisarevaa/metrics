package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestErrCheckAnalyzer(t *testing.T) {
	// функция analysistest.Run применяет тестируемый анализатор ErrCheckAnalyzer
	// к пакетам из папки testdata и проверяет ожидания
	// ./... — проверка всех поддиректорий в testdata
	// можно указать ./pkg1 для проверки только pkg1
	analysistest.Run(t, analysistest.TestData(), ErrCheckAnalyzer, "./...")
}
