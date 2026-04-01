import pytest
from fastapi.testclient import TestClient
from main import app

client = TestClient(app)

def test_health():
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json() == {"status": "ok"}

def test_create_item():
    item_data = {
        "name": "Test Item",
        "price": 99.99
    }
    
    response = client.post("/items", json=item_data)
    assert response.status_code == 201
    
    data = response.json()
    assert "id" in data
    assert data["name"] == "Test Item"
    assert data["price"] == 99.99
    assert isinstance(data["id"], str)
    
    # Сохраняем ID для следующего теста
    return data["id"]

def test_get_item_not_found():
    response = client.get("/items/non-existent-id")
    assert response.status_code == 404
    assert response.json()["detail"] == "Item not found"

def test_get_item_success():
    # Сначала создаем элемент
    item_data = {
        "name": "Get Test",
        "price": 49.99
    }
    
    create_response = client.post("/items", json=item_data)
    assert create_response.status_code == 201
    created_item = create_response.json()
    item_id = created_item["id"]
    
    # Теперь получаем его
    get_response = client.get(f"/items/{item_id}")
    assert get_response.status_code == 200
    
    retrieved_item = get_response.json()
    assert retrieved_item["id"] == item_id
    assert retrieved_item["name"] == "Get Test"
    assert retrieved_item["price"] == 49.99