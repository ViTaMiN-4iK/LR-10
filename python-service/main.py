import uuid
import asyncio
from contextlib import asynccontextmanager

from fastapi import FastAPI, HTTPException, status
from fastapi.responses import JSONResponse

from schemas import ItemCreate, Item
from storage import ItemStorage

# Глобальный объект хранилища
storage = ItemStorage()

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    print("Starting up FastAPI server...")
    yield
    # Shutdown
    print("Shutting down FastAPI server...")
    # Здесь можно добавить очистку ресурсов

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

# Для запуска через скрипт
if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)