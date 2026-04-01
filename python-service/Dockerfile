FROM python:3.11-slim

WORKDIR /app

# Install dependencies
COPY python-service/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy source code (включая сгенерированные pb файлы)
COPY python-service/ .

EXPOSE 8000

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]