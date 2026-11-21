
import csv
import io
import os
import sys
from pathlib import Path
from datetime import datetime, timezone

import requests
import psycopg2

from dotenv import load_dotenv
from psycopg2.extras import execute_values


BASE_DIR = Path(__file__).resolve().parent.parent

load_dotenv(BASE_DIR / ".env")

GITHUB_API_URL = (
    "https://api.github.com/repos/"
    "lokabisa-oss/region-id/releases/latest"
)

DB_CONFIG = {
    "host": os.getenv("DB_HOST", "localhost"),
    "port": os.getenv("DB_PORT", "5432"),
    "dbname": os.getenv("DB_NAME"),
    "user": os.getenv("DB_USER"),
    "password": os.getenv("DB_PASS"),
}

TABLE_NAME = "master_cities"
BATCH_SIZE = 100

def validate_config():
    required_vars = [
        "DB_NAME",
        "DB_USER",
        "DB_PASSWORD",
    ]

    missing = [
        key for key in required_vars
        if not DB_CONFIG.get({
            "DB_NAME": "dbname",
            "DB_USER": "user",
            "DB_PASSWORD": "password",
        }[key])
    ]

    if missing:
        raise EnvironmentError(
            f"Missing database configuration: {', '.join(missing)}"
        )

def fetch_city_data():
    print("[INFO] Fetching Indonesia administrative data...")

    response = requests.get(
        GITHUB_API_URL,
        timeout=30
    )
    response.raise_for_status()

    release = response.json()

    asset = next(
        (
            item for item in release.get("assets", [])
            if item["name"].lower() == "regencies.csv"
        ),
        None
    )

    if asset is None:
        raise RuntimeError(
            "regencies.csv not found in the latest GitHub release"
        )

    print(f"[INFO] Dataset version: {release['tag_name']}")

    csv_response = requests.get(
        asset["browser_download_url"],
        timeout=60
    )
    csv_response.raise_for_status()

    reader = csv.DictReader(
        io.StringIO(
            csv_response.content.decode("utf-8-sig")
        )
    )

    if not reader.fieldnames or "name" not in reader.fieldnames:
        raise RuntimeError(
            f"Unexpected CSV structure: {reader.fieldnames}"
        )

    cities = []

    for row in reader:
        name = (row.get("name") or "").strip()

        if name:
            cities.append(name)

    if not cities:
        raise RuntimeError("Dataset contains no city records")

    if len(cities) != len(set(cities)):
        raise RuntimeError("Duplicate city names found in dataset")

    print(f"[INFO] Retrieved {len(cities)} records")

    return cities

def seed_cities(cities):
    connection = None

    try:
        print("[INFO] Connecting to PostgreSQL...")

        connection = psycopg2.connect(**DB_CONFIG)

        with connection:
            with connection.cursor() as cursor:

                # Ensure name uniqueness for ON CONFLICT
                cursor.execute("""
                    CREATE UNIQUE INDEX IF NOT EXISTS
                    idx_master_cities_name
                    ON master_cities (name);
                """)

                now = datetime.now(timezone.utc)

                records = [
                    (name, now, now)
                    for name in cities
                ]

                insert_query = """
                    INSERT INTO master_cities
                        (name, created_at, updated_at)
                    VALUES %s
                    ON CONFLICT (name) DO NOTHING;
                """

                execute_values(
                    cursor,
                    insert_query,
                    records,
                    page_size=BATCH_SIZE
                )

                inserted = cursor.rowcount

                cursor.execute("""
                    SELECT COUNT(*)
                    FROM master_cities;
                """)

                total = cursor.fetchone()[0]

        print("\n[SUCCESS] Master city seeding completed!")
        print(f"[INFO] Newly inserted : {inserted}")
        print(f"[INFO] Total records  : {total}")

    except Exception as error:
        print(f"[ERROR] Database seeding failed: {error}")
        raise

    finally:
        if connection:
            connection.close()

def main():
    try:
        validate_config()

        cities = fetch_city_data()

        seed_cities(cities)

    except Exception as error:
        print(f"\n[FATAL] {error}")
        sys.exit(1)


if __name__ == "__main__":
    main()


# run this with `python scripts/inject_master_city.py`