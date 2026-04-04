# main.py
from fastapi import FastAPI, HTTPException
import httpx
from pydantic import BaseModel

app = FastAPI()

# Модели данных
class UserData(BaseModel):
    user_id: int
    name: str
    age: int

class GoResponse(BaseModel):
    status: str
    message: str

# Адрес Go‑сервиса (при запуске в одной сети можно использовать localhost:8080)
GO_SERVICE_URL = "http://localhost:8080"

@app.post("/process-user", response_model=GoResponse)
async def process_user(user: UserData):
    """
    Эндпоинт FastAPI принимает данные пользователя и передаёт их в Go‑сервис.
    """
    async with httpx.AsyncClient() as client:
        try:
            response = await client.post(
                f"{GO_SERVICE_URL}/data",
                json=user.dict(),
                timeout=5.0
            )
            response.raise_for_status()
            return response.json()
        except httpx.RequestError as e:
            raise HTTPException(status_code=503, detail=f"Go service unavailable: {str(e)}")
        except httpx.HTTPStatusError as e:
            raise HTTPException(status_code=e.response.status_code, detail=e.response.text)

@app.get("/ping-go")
async def ping_go():
    """Проверка доступности Go‑сервиса."""
    async with httpx.AsyncClient() as client:
        try:
            resp = await client.get(f"{GO_SERVICE_URL}/ping")
            return {"go_status": resp.json()}
        except Exception as e:
            raise HTTPException(status_code=503, detail=str(e))