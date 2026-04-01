import grpc
import asyncio
from typing import Optional
import items_pb2
import items_pb2_grpc

class GoGrpcClient:
    def __init__(self, host: str = "localhost", port: int = 50051):
        self.channel = grpc.aio.insecure_channel(f"{host}:{port}")
        self.stub = items_pb2_grpc.ItemServiceStub(self.channel)
    
    async def get_item(self, item_id: str):
        """Получить элемент по ID через gRPC"""
        try:
            request = items_pb2.GetItemRequest(id=item_id)
            response = await self.stub.GetItem(request, timeout=5.0)
            return {
                "id": response.id,
                "name": response.name,
                "price": response.price
            }
        except grpc.RpcError as e:
            if e.code() == grpc.StatusCode.NOT_FOUND:
                return None
            raise
    
    async def create_item(self, name: str, price: float):
        """Создать элемент через gRPC"""
        try:
            request = items_pb2.CreateItemRequest(name=name, price=price)
            response = await self.stub.CreateItem(request, timeout=5.0)
            return {
                "id": response.id,
                "name": response.name,
                "price": response.price
            }
        except grpc.RpcError as e:
            raise
    
    async def close(self):
        await self.channel.close()