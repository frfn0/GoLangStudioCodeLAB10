Сурков Всеволод Сергеевич, группа 220032-11, вариант №2, лабораторная работа №10
Для задания №2:
Проверка эндпоинта /ping :
  curl http://localhost:8080/ping
Проверка эндпоинта /hello/:name :
  curl http://localhost:8080/hello/World
Для задания №4:
Тест:
  Invoke-RestMethod -Uri "http://localhost:8000/process-user" `
                    -Method Post `
                    -ContentType "application/json" `
                    -Body '{"user_id": 1, "name": "Alice", "age": 30}'
Для задания №6:
Запуск тестирования с помощью hey для Go-сервиса:
  hey -n 10000 -c 100 http://localhost:8080/ping > test_go.txt
Запуск тестирования с помощью hey для FastAPI-сервиса:
  hey -n 10000 -c 100 http://localhost:8000/ping > test_fastapi.txt
Результаты записываются в соответствующие файлы, их сравнение произведено в файле comparison.txt
