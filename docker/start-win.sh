#!/bin/bash

# Убиваем контейнеры
docker kill ada-app
docker kill ada-db

# Удаляем контейнеры
docker rm ada-app
docker rm ada-db

# Удаляем образ
docker rmi ada-telegram-bot-app

# Запускаем docker-compose в фоновом режиме и выводим логи в файл
docker compose -p ada-telegram-bot up > logs.log