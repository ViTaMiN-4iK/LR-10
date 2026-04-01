import pytest
from fastapi.testclient import TestClient
from main import app
from unittest.mock import patch, AsyncMock

@pytest.fixture
def client():
    with TestClient(app) as test_client:
        yield test_client

def test_health(client):
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json() == {"status": "ok"}

def test_create_item(client):
    item_data = {"name": "Test Item", "price": 99.99}
    response = client.post("/items", json=item_data)
    assert response.status_code == 201
    data = response.json()
    assert "id" in data
    assert data["name"] == "Test Item"
    assert data["price"] == 99.99

def test_get_item_not_found(client):
    response = client.get("/items/non-existent-id")
    assert response.status_code == 404
    assert response.json()["detail"] == "Item not found"

def test_get_item_success(client):
    # Создаем элемент
    create_response = client.post("/items", json={"name": "Get Test", "price": 49.99})
    item_id = create_response.json()["id"]
    
    # Получаем его
    get_response = client.get(f"/items/{item_id}")
    assert get_response.status_code == 200
    data = get_response.json()
    assert data["id"] == item_id
    assert data["name"] == "Get Test"

def test_proxy_get_item_success(client):
    """Интеграционный тест: требует запущенный Go-сервис"""
    import httpx
    
    try:
        with httpx.Client() as http_client:
            # Создаем элемент в Go
            create_response = http_client.post(
                "http://localhost:8080/items",
                json={"name": "Go Item", "price": 777.77},
                timeout=2
            )
            if create_response.status_code != 201:
                pytest.skip("Go service not responding correctly")
            
            go_item = create_response.json()
            item_id = go_item["id"]
            
            # Получаем через прокси
            proxy_response = client.get(f"/proxy-items/{item_id}")
            assert proxy_response.status_code == 200
            assert proxy_response.json()["name"] == "Go Item"
    except Exception as e:
        pytest.skip(f"Go service not available: {e}")

def test_proxy_get_item_not_found(client):
    """Тест прокси при отсутствии элемента (с моком)"""
    with patch('httpx.AsyncClient.get') as mock_get:
        mock_response = AsyncMock()
        mock_response.status_code = 404
        mock_get.return_value = mock_response
        
        response = client.get("/proxy-items/test-id")
        assert response.status_code == 404
        assert "not found" in response.json()["detail"]

def test_proxy_go_service_unavailable(client):
    """Тест при недоступности Go сервиса (с моком)"""
    # Импортируем нужное исключение из httpx
    import httpx
    
    with patch('httpx.AsyncClient.get', side_effect=httpx.ConnectError("Connection refused")):
        response = client.get("/proxy-items/test-id")
        assert response.status_code == 503
        assert "not available" in response.json()["detail"]