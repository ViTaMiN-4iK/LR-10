import pytest
import requests
import time
import subprocess
import os
import signal
import sys

class TestRealIntegration:
    """Интеграционные тесты с реальными запущенными сервисами"""
    
    @classmethod
    def setup_class(cls):
        """Проверяем, что сервисы запущены"""
        cls.go_running = False
        cls.python_running = False
        
        # Проверяем Go сервис
        try:
            response = requests.get("http://localhost:8080/health", timeout=2)
            if response.status_code == 200:
                cls.go_running = True
                print("✓ Go service is running")
        except:
            print("✗ Go service is NOT running")
        
        # Проверяем Python сервис
        try:
            response = requests.get("http://localhost:8000/health", timeout=2)
            if response.status_code == 200:
                cls.python_running = True
                print("✓ Python service is running")
        except:
            print("✗ Python service is NOT running")
    
    def test_real_grpc_create_and_get(self):
        """Реальный тест gRPC через Python API (требует запущенные сервисы)"""
        if not (self.go_running and self.python_running):
            pytest.skip("Services not running, skipping real integration test")
        
        # Создаем элемент через gRPC
        create_response = requests.post(
            "http://localhost:8000/grpc-items",
            json={"name": "Real gRPC Test", "price": 555.55}
        )
        assert create_response.status_code == 200
        created = create_response.json()
        assert "id" in created
        item_id = created["id"]
        
        # Получаем элемент через gRPC
        get_response = requests.get(f"http://localhost:8000/grpc-items/{item_id}")
        assert get_response.status_code == 200
        retrieved = get_response.json()
        assert retrieved["id"] == item_id
        assert retrieved["name"] == "Real gRPC Test"
        assert retrieved["price"] == 555.55
    
    def test_real_rest_proxy(self):
        """Реальный тест REST прокси (требует запущенные сервисы)"""
        if not (self.go_running and self.python_running):
            pytest.skip("Services not running, skipping real integration test")
        
        # Создаем элемент напрямую в Go
        create_response = requests.post(
            "http://localhost:8080/items",
            json={"name": "Direct Go Item", "price": 777.77}
        )
        assert create_response.status_code == 201
        go_item = create_response.json()
        item_id = go_item["id"]
        
        # Получаем через Python прокси
        proxy_response = requests.get(f"http://localhost:8000/proxy-items/{item_id}")
        assert proxy_response.status_code == 200
        proxy_item = proxy_response.json()
        assert proxy_item["id"] == item_id
        assert proxy_item["name"] == "Direct Go Item"
        assert proxy_item["price"] == 777.77
    
    def test_real_end_to_end_flow(self):
        """Сквозной тест: создание в Go, получение через gRPC и REST"""
        if not (self.go_running and self.python_running):
            pytest.skip("Services not running, skipping real integration test")
        
        # 1. Создаем в Go напрямую
        create_response = requests.post(
            "http://localhost:8080/items",
            json={"name": "E2E Test", "price": 999.99}
        )
        assert create_response.status_code == 201
        item_id = create_response.json()["id"]
        
        # 2. Получаем через gRPC
        grpc_response = requests.get(f"http://localhost:8000/grpc-items/{item_id}")
        assert grpc_response.status_code == 200
        assert grpc_response.json()["name"] == "E2E Test"
        
        # 3. Получаем через REST прокси
        proxy_response = requests.get(f"http://localhost:8000/proxy-items/{item_id}")
        assert proxy_response.status_code == 200
        assert proxy_response.json()["name"] == "E2E Test"
        
        # 4. Получаем напрямую из Go
        go_response = requests.get(f"http://localhost:8080/items/{item_id}")
        assert go_response.status_code == 200
        assert go_response.json()["name"] == "E2E Test"
    
    def test_real_grpc_not_found(self):
        """Реальный тест 404 через gRPC"""
        if not (self.go_running and self.python_running):
            pytest.skip("Services not running, skipping real integration test")
        
        response = requests.get("http://localhost:8000/grpc-items/real-non-existent-id")
        assert response.status_code == 404
        assert "not found" in response.json()["detail"]
    
    def test_real_proxy_not_found(self):
        """Реальный тест 404 через REST прокси"""
        if not (self.go_running and self.python_running):
            pytest.skip("Services not running, skipping real integration test")
        
        response = requests.get("http://localhost:8000/proxy-items/real-non-existent-id")
        assert response.status_code == 404
        assert "not found" in response.json()["detail"]