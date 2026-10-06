"""Задание 6. FastAPI — второй сервис для сравнения скорости под нагрузкой."""

from fastapi import FastAPI

app = FastAPI(title="Lab 10 Task 6: FastAPI", version="1.0.0")


@app.get("/ping")
async def ping() -> dict[str, str]:
    """Ответ, полностью аналогичный ответу Gin-сервиса."""
    return {"message": "pong"}