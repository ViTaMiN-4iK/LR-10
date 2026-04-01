import uuid
import asyncio
from contextlib import asynccontextmanager
import httpx
from fastapi import FastAPI, HTTPException, status
from fastapi.responses import JSONResponse

from schemas import ItemCreate, Item
from storage import ItemStorage
from grpc_client import GoGrpcClient

# Глобальный объект хранилища
storage = ItemStorage()
grpc_client = None

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    print("Starting up FastAPI server...")
    # Создаем HTTP клиент для Go-сервиса (REST)
    app.state.go_client = httpx.AsyncClient(base_url="http://localhost:8080")
    
    # Создаем gRPC клиент для Go-сервиса
    global grpc_client
    grpc_client = GoGrpcClient(host="localhost", port=50051)
    
    yield
    
    # Shutdown
    print("Shutting down FastAPI server...")
    await app.state.go_client.aclose()
    if grpc_client:
        await grpc_client.close()

app = FastAPI(lifespan=lifespan)

@app.get("/health")
async def health():
    return {"status": "ok"}

@app.post("/items", status_code=status.HTTP_201_CREATED)
async def create_item(item: ItemCreate):
    # Генерируем UUID
    item_id = str(uuid.uuid4())
    
    # Создаем объект Item
    new_item = Item(
        id=item_id,
        name=item.name,
        price=item.price
    )
    
    # Сохраняем в хранилище
    await storage.create(new_item)
    
    return new_item

@app.get("/items/{item_id}")
async def get_item(item_id: str):
    item = await storage.get(item_id)
    
    if not item:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Item not found"
        )
    
    return item

# REST прокси к Go-сервису (через HTTP)
@app.get("/proxy-items/{item_id}")
async def proxy_get_item(item_id: str):
    """
    Прокси-запрос к Go-сервису для получения элемента по ID (через REST)
    """
    try:
        # Вызываем Go-сервис
        response = await app.state.go_client.get(f"/items/{item_id}")
        
        # Проверяем статус ответа
        if response.status_code == 404:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Item {item_id} not found in Go service"
            )
        
        response.raise_for_status()
        
        # Возвращаем данные от Go-сервиса
        return response.json()
        
    except httpx.ConnectError:
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail="Go service is not available"
        )
    except httpx.TimeoutException:
        raise HTTPException(
            status_code=status.HTTP_504_GATEWAY_TIMEOUT,
            detail="Go service timeout"
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Error communicating with Go service: {str(e)}"
        )

# Новый эндпоинт: gRPC прокси к Go-сервису
@app.get("/grpc-items/{item_id}")
async def grpc_get_item(item_id: str):
    """
    Получить элемент из Go-сервиса через gRPC
    """
    global grpc_client
    
    try:
        item = await grpc_client.get_item(item_id)
        
        if not item:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Item {item_id} not found in Go service (gRPC)"
            )
        
        return item
        
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"gRPC error: {str(e)}"
        )

@app.post("/grpc-items")
async def grpc_create_item(item: ItemCreate):
    """
    Создать элемент в Go-сервисе через gRPC
    """
    global grpc_client
    
    try:
        created_item = await grpc_client.create_item(item.name, item.price)
        return created_item
        
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"gRPC error: {str(e)}"
        )

# Для запуска через скрипт
if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)