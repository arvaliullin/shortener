"""Создает тестовые короткие ссылки через POST /api/shorten."""

import argparse
import os
import sys
import time
from pathlib import Path

import requests
from dotenv import load_dotenv


def main() -> None:
    load_dotenv(Path(__file__).resolve().parents[2] / ".env")
    count = read_count()
    token = access_token()
    created = create_links("http://localhost:8080", token, count)
    print(f"создано {created}")


def read_count() -> int:
    parser = argparse.ArgumentParser(description="Создать тестовые короткие ссылки")
    parser.add_argument("--count", type=int, default=100)
    count = parser.parse_args().count
    if count < 1:
        sys.exit("count must be positive")
    return count


def access_token() -> str:
    response = requests.post(
        env("KEYCLOAK_TOKEN_URL"),
        data={
            "grant_type": "client_credentials",
            "client_id": env("KEYCLOAK_DEBUG_CLIENT_ID"),
            "client_secret": env("KEYCLOAK_DEBUG_CLIENT_SECRET"),
        },
        timeout=5,
    )
    response.raise_for_status()
    return response.json()["access_token"]


def create_links(base: str, token: str, count: int) -> int:
    stamp = int(time.time())
    for index in range(count):
        original = f"https://example.com/seed/{stamp}/{index}"
        shorten(base, token, original)
    return count


def shorten(base: str, token: str, original: str) -> None:
    response = requests.post(
        f"{base}/api/shorten",
        json={"url": original},
        headers={"Authorization": f"Bearer {token}"},
        timeout=5,
    )
    if response.status_code != 201:
        sys.exit(f"shorten failed: {response.status_code} {response.text}")


def env(name: str) -> str:
    value = os.environ.get(name)
    if not value:
        sys.exit(f"missing {name}")
    return value


if __name__ == "__main__":
    main()
