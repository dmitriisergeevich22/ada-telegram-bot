package ada

import "fmt"

// Обработчик цепочек
// TODO
func (a *AdaBot) handlerChain(userId int64, data string) error {
	fmt.Println("START handlerChain")
	_, err := a.db.GetLastSession(userId)
	if err != nil {
		return fmt.Errorf("user_id: %d; error db.GetLastSession: %w", userId, err)
	}

	return nil
}
