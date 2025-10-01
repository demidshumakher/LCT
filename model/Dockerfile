# Базовый образ с Python
FROM python:3.10-slim

# Установим зависимости для PyTorch
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Рабочая директория
WORKDIR /app

# Копируем файлы
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

EXPOSE 8081

# Запускаем сервер
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8081"]
