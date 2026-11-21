
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


def fetch_province_data():
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
            if item["name"].lower() == "provinces.csv"
        ),
        None
    )

    if asset is None:
        raise RuntimeError(
            "provinces.csv not found in the latest GitHub release"
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
        "name"
    }

    if not reader.fieldnames or not required_columns.issubset(
        reader.fieldnames
    ):
        raise RuntimeError(
            f"Unexpected CSV structure: {reader.fieldnames}"
        )

    provinces = []

    for row in reader:
        code = (row.get("code") or "").strip()
        name = (row.get("name") or "").strip()

        if not code or not name:
            raise RuntimeError(
                f"Invalid province record found: {row}"
            )

        provinces.append({
            "code": code,
            "name": name
        })

    if not provinces:
        raise RuntimeError("Dataset contains no province records")

    province_codes = [
        province["code"]
        for province in provinces
    ]

    if len(province_codes) != len(set(province_codes)):
        raise RuntimeError("Duplicate province codes found")

    print(f"[INFO] Retrieved {len(provinces)} records")

    return provinces


def seed_provinces(provinces):
    connection = None

    try:
        print("[INFO] Connecting to PostgreSQL...")

        connection = psycopg2.connect(**DB_CONFIG)

        with connection:
            with connection.cursor() as cursor:

                now = datetime.now(timezone.utc)

                records = [
                    (
                        province["code"],
                        province["name"],
                        now,
                        now
                    )
                    for province in provinces
                ]

                insert_query = """
                    INSERT INTO master_provinces
                        (code, name, created_at, updated_at)
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
                    FROM master_provinces;
                    """
                )

                total = cursor.fetchone()[0]

        print("\n[SUCCESS] Master province seeding completed!")
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

        provinces = fetch_province_data()

        seed_provinces(provinces)

    except Exception as error:
        print(f"\n[FATAL] {error}")
        sys.exit(1)


if __name__ == "__main__":
    main()


# Run:
# python scripts/inject_master_province.py
