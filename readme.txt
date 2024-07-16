curl запросы:
1. Создание нового аккаунта:
    curl -X POST http://localhost:8080/accounts -H "Content-Type: application/json" -d "{\"id\":1}"   
2. Пополнение баланса:
    curl -X POST http://localhost:8080/accounts/1/deposit -H "Content-Type: application/json" -d "{\"amount\":111}"
3. Снятие средств:
    curl -X POST http://localhost:8080/accounts/1/withdraw -H "Content-Type: application/json" -d "{\"amount\":200}"
4. Проверка баланса:
    curl -X GET http://localhost:8080/accounts/1/balance