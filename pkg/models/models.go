package models

type Menu string

const (
	MenuStart   Menu = "start"
	MenuSupport Menu = "support"
)

// Запущенная цепочка
type Chain struct {
	Name string `json:"name"`
	Step string `json:"step"`
}

// Запущенная функция
type Func struct {
	Name string `json:"name"`
	Step string `json:"step"`
}

// Сессия
type Session struct {
	M    []Menu                 `json:"menu,omitempty"`  // Открытые меню
	C    *Chain                 `json:"chain,omitempty"` // Запущенная цепочка
	F    *Func                  `json:"func,omitempty"`  // Запущенная функция
	Data map[string]interface{} `json:"data,omitempty"`  // Данные сессии
}
