import pytest
import subprocess
import time
import requests

class TestIntegration:
    @classmethod
    def setup_class(cls):
        """Запускаем Go сервис перед тестами"""
        cls.go_process = subprocess.Popen(
            ["go", "run", "main.go"],
            cwd="../go-service",
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )
        # Ждем запуска
        time.sleep(2)
    
    @classmethod
    def teardown_class(cls):
        """Останавливаем Go сервис после тестов"""
        cls.go_process.terminate()
        cls.go_process.wait()
    
    def test_rest_proxy_integration(self):
        """Интеграционный тест REST прокси"""
        # Создаем элемент в Go через прямой запрос
        response = requests.post("http://localhost:8080/items", json={
            "name": "Integration Test",
            "price": 123.45
        })
        assert response.status_code == 201
        item_id = response.json()["id"]
        
        # Получаем через Python прокси
        proxy_response = requests.get(f"http://localhost:8000/proxy-items/{item_id}")
        assert proxy_response.status_code == 200
        assert proxy_response.json()["name"] == "Integration Test"