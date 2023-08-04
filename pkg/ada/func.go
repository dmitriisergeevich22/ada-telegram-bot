package ada

import "fmt"

// Обработчик функций
// TODO
func (a *AdaBot) handlerFunc(userId int64, data string) error {
	fmt.Println("START handlerFunc")
	_, err := a.db.GetLastSession(userId)
	if err != nil {
		return fmt.Errorf("user_id: %d; error db.GetLastSession: %w", userId, err)
	}

	return nil
}
