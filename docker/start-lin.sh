#!/bin/bash

# Убиваем контейнеры
sudo docker kill ada-app
sudo docker kill ada-db

# Удаляем контейнеры
sudo docker rm ada-app
sudo docker rm ada-db

# Удаляем образ
sudo docker rmi ada-telegram-bot-app

# Запускаем docker-compose в фоновом режиме и выводим логи в файл
sudo nohup docker compose -p ada-telegram-bot up > logs-ada-telegram-bot.txt &