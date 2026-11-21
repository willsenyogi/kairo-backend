
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

BATCH_SIZE = 100


def validate_config():
    required_vars = {
        "DB_NAME": DB_CONFIG["dbname"],
        "DB_USER": DB_CONFIG["user"],
        "DB_PASS": DB_CONFIG["password"],
    }

    missing = [
        key for key, value in required_vars.items()
        if not value
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

    required_columns = {
        "code",
        "province_code",
        "name"
    }

    if not reader.fieldnames or not required_columns.issubset(
        reader.fieldnames
    ):
        raise RuntimeError(
            f"Unexpected CSV structure: {reader.fieldnames}"
        )

    cities = []

    for row in reader:
        code = (row.get("code") or "").strip()
        province_code = (row.get("province_code") or "").strip()
        name = (row.get("name") or "").strip()

        if not code or not province_code or not name:
            raise RuntimeError(
                f"Invalid city record found: {row}"
            )

        cities.append({
            "code": code,
            "province_code": province_code,
            "name": name
        })

    if not cities:
        raise RuntimeError("Dataset contains no city records")

    city_codes = [city["code"] for city in cities]

    if len(city_codes) != len(set(city_codes)):
        raise RuntimeError("Duplicate city codes found in dataset")

    print(f"[INFO] Retrieved {len(cities)} records")

    return cities


def seed_cities(cities):
    connection = None

    try:
        print("[INFO] Connecting to PostgreSQL...")

        connection = psycopg2.connect(**DB_CONFIG)

        with connection:
            with connection.cursor() as cursor:
                province_codes = list({
                    city["province_code"]
                    for city in cities
                })

                cursor.execute(
                    """
                    SELECT id, code
                    FROM master_provinces
                    WHERE code = ANY(%s);
                    """,
                    (province_codes,)
                )

                province_map = {
                    str(code): province_id
                    for province_id, code in cursor.fetchall()
                }

                missing_provinces = (
                    set(province_codes) - set(province_map.keys())
                )

                if missing_provinces:
                    raise RuntimeError(
                        "Missing province codes in database: "
                        f"{', '.join(sorted(missing_provinces))}"
                    )

                now = datetime.now(timezone.utc)

                records = [
                    (
                        province_map[city["province_code"]],
                        city["code"],
                        city["name"],
                        now,
                        now
                    )
                    for city in cities
                ]

                insert_query = """
                    INSERT INTO master_cities
                        (province_id, code, name, created_at, updated_at)
                    VALUES %s
                    ON CONFLICT DO NOTHING;
                """

                execute_values(
                    cursor,
                    insert_query,
                    records,
                    page_size=BATCH_SIZE
                )

                inserted = cursor.rowcount

                cursor.execute(
                    """
                    SELECT COUNT(*)
                    FROM master_cities;
                    """
                )

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


# Run:
# python scripts/inject_master_city.py
