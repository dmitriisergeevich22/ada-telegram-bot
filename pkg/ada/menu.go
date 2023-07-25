package ada

import "fmt"

// Обработчик меню
// TODO
func (a *AdaBot) handlerMenu(userId int64, data string) error {
	fmt.Println("START handlerMenu")
	_, err := a.db.GetLastSession(userId)
	if err != nil {
		return fmt.Errorf("user_id: %d; error db.GetLastSession: %w", userId, err)
	}

	return nil
}
