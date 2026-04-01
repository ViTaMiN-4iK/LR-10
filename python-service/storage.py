import asyncio
from typing import Dict, Optional
from schemas import Item  # Правильный импорт

class ItemStorage:
    def __init__(self):
        self._items: Dict[str, Item] = {}
        self._lock = asyncio.Lock()
    
    async def create(self, item: Item) -> Item:
        async with self._lock:
            self._items[item.id] = item
            return item
    
    async def get(self, item_id: str) -> Optional[Item]:
        async with self._lock:
            return self._items.get(item_id)