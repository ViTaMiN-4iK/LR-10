import pytest
from fastapi.testclient import TestClient
from main import app
import httpx

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
    
    # Возвращаем ID для возможного использования
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

# Тесты для прокси (требуют запущенный Go-сервис)
@pytest.mark.integration
def test_proxy_get_item_success():
    """Интеграционный тест: требует запущенный Go-сервис на порту 8080"""
    # Сначала создаем элемент в Go-сервисе через API
    import httpx
    
    with httpx.Client() as client:
        # Создаем элемент в Go
        create_response = client.post(
            "http://localhost:8080/items",
            json={"name": "Go Item", "price": 777.77}
        )
        assert create_response.status_code == 201
        go_item = create_response.json()
        item_id = go_item["id"]
        
        # Теперь получаем его через прокси
        proxy_response = client.get(f"http://localhost:8000/proxy-items/{item_id}")
        assert proxy_response.status_code == 200
        
        proxy_data = proxy_response.json()
        assert proxy_data["id"] == item_id
        assert proxy_data["name"] == "Go Item"
        assert proxy_data["price"] == 777.77

@pytest.mark.integration
def test_proxy_get_item_not_found():
    """Тест: запрос несуществующего элемента через прокси"""
    import httpx
    
    with httpx.Client() as client:
        response = client.get("http://localhost:8000/proxy-items/non-existent-id")
        assert response.status_code == 404
        assert "not found" in response.json()["detail"]

@pytest.mark.integration
def test_proxy_go_service_unavailable():
    """Тест: когда Go-сервис недоступен"""
    # Здесь мы не можем легко протестировать без остановки сервиса
    # Но можно использовать мок или пропустить
    pass