import pytest
from fastapi.testclient import TestClient
from main import app

@pytest.fixture
def client():
    with TestClient(app) as test_client:
        yield test_client

def test_grpc_create_item(client):
    """Тест создания элемента через gRPC"""
    response = client.post("/grpc-items", json={"name": "Test Item", "price": 99.99})
    assert response.status_code == 200
    data = response.json()
    assert "id" in data
    assert data["name"] == "Test Item"
    assert data["price"] == 99.99
    # Не возвращаем значение, просто сохраняем в переменной
    return  # Убираем возврат

def test_grpc_get_item(client):
    """Тест получения элемента через gRPC"""
    # Сначала создаем
    create_response = client.post("/grpc-items", json={"name": "Get Test", "price": 49.99})
    assert create_response.status_code == 200
    item_id = create_response.json()["id"]
    
    # Получаем
    get_response = client.get(f"/grpc-items/{item_id}")
    assert get_response.status_code == 200
    data = get_response.json()
    assert data["id"] == item_id
    assert data["name"] == "Get Test"

def test_grpc_get_item_not_found(client):
    """Тест 404 через gRPC"""
    response = client.get("/grpc-items/non-existent-id")
    # Сначала проверим статус
    assert response.status_code == 404, f"Expected 404, got {response.status_code}: {response.text}"
    data = response.json()
    assert "not found" in data["detail"]