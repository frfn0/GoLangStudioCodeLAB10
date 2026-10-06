"""Задание 4. FastAPI-сервис, который вызывает Go-сервис по HTTP."""

from __future__ import annotations

import os

import httpx
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field

app = FastAPI(
    title="Lab 10 Task 4: FastAPI -> Go",
    description="Python-сервис принимает данные и передаёт их Go-сервису по HTTP.",
    version="1.0.0",
)

GO_SERVICE_URL = os.getenv("GO_SERVICE_URL", "http://localhost:8080")
REQUEST_TIMEOUT_SECONDS = 5.0


class UserData(BaseModel):
    """Запрос, который приходит в FastAPI."""

    user_id: int = Field(gt=0, examples=[1])
    name: str = Field(min_length=1, max_length=100, examples=["Alice"])
    age: int = Field(ge=0, le=150, examples=[30])


class GoResponse(BaseModel):
    """Ответ Go-сервиса."""

    status: str
    message: str


@app.post("/process-user", response_model=GoResponse)
async def process_user(user: UserData) -> GoResponse:
    """Передать пользователя в Go-сервис и вернуть его ответ."""
    try:
        async with httpx.AsyncClient(timeout=REQUEST_TIMEOUT_SECONDS) as client:
            response = await client.post(
                f"{GO_SERVICE_URL}/data",
                json=user.model_dump(),
            )
            response.raise_for_status()
            return GoResponse(**response.json())
    except httpx.RequestError as exc:
        raise HTTPException(
            status_code=503, detail=f"Go service unavailable: {exc}"
        ) from exc
    except httpx.HTTPStatusError as exc:
        raise HTTPException(
            status_code=exc.response.status_code,
            detail=exc.response.text,
        ) from exc


@app.get("/ping-go")
async def ping_go() -> dict[str, str]:
    """Проверить доступность Go-сервиса."""
    try:
        async with httpx.AsyncClient(timeout=REQUEST_TIMEOUT_SECONDS) as client:
            response = await client.get(f"{GO_SERVICE_URL}/ping")
            response.raise_for_status()
            return {"go_status": response.json()["message"]}
    except httpx.HTTPError as exc:
        raise HTTPException(status_code=503, detail=str(exc)) from exc