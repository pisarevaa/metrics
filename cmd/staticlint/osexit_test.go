package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// Тестирование анализатора os.Exit().
func TestOSExitInMainAnalyzer(t *testing.T) {
	// функция analysistest.Run применяет тестируемый анализатор OSExitInMainAnalyzer
	// к пакетам из папки testdata и проверяет ожидания
	// ./... — проверка всех поддиректорий в testdata
	// можно указать ./pkg1 для проверки только pkg1
	analysistest.Run(t, analysistest.TestData(), OSExitInMainAnalyzer, "./...")
}
