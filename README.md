# Ada Telegram Bot
**Developer:** `https://t.me/DmitriySergeevich22`
**Bot:** `https://t.me/Ada_Telegram_Bot`

The Advertising Assistant is designed to help you run and manage your advertisements. It allows you to create ad buying events. Create ad sales events. We will notify you about previously created events. Allows you to see all your sales and purchases, create averages.
In the future it is planned:
- Automatic check of the arrival of subscribers (collection of statistics from advertising).
- Implementation of the average cost of a subscriber in a tegram.
- Collecting general data and providing averages among all clients.


# Wiki
## Понятия
### Меню (Menu)
Данный тип может вызывать только другое меню (menu) или цепочки (chain). 

### Цепочки (Chain)
Данный тип может вызывать только функции (func). Это связанный список функций имеющих атомарное изменение данных.
Вано: Цепочки сперва собирают все данные - а затем выполняют атомарно операцию изменения данных.

**Пример**: 
```
// Изменение цены рекламной интеграции
func cAdEditPrice() error {
    // Получение новой цены
    price := getPrice()
    // Изменение цены
    setAdPrice(price)
}   

```

### Функции (Func)
Данный тип это функции которые выполняют одно действие.

**Пример**:
```
// Получение новой цены:
func getNewPrice() {
    // Отправка сообщения о том что необходимо прислать новую цену
    setMessage('Требуется отправить новую стоимость:')
    // Получение нувой цены от пользователя и запись в session
    session.Set('newPrice') = dataPrice
}
```