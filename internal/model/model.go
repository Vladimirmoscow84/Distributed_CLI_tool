package model

// Параметры поиска
type GrepConfig struct {
	Pattern    string //строка для поиска
	IgnoreCase bool   //флаг игнорирования регситра
	ShowNumber bool   //необходимость показывания номера строки
	Invert     bool   //необходимость инвертирования совпадения
}

// результат работы grep
// потом будет расширен!!!
type GrepResult struct {
	Lines []string
}
