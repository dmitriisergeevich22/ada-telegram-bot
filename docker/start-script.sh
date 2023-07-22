#!/bin/bash

# Переходим в директорию docker
cd /home/dasy/myprojects/AdaTelegramBot/docker/

# Убиваем контейнеры ada_app и ada_db
sudo docker kill ada_app
sudo docker kill ada_db

# Удаляем контейнеры ada_app и ada_db
sudo docker rm ada_app
sudo docker rm ada_db

# Удаляем образ ada_telegram_bot-app
sudo docker rmi ada_telegram_bot-app

# Запускаем docker-compose в фоновом режиме и выводим логи в файл
sudo nohup docker compose -p ada_telegram_bot up > logs_ada_telegram_bot.txt &