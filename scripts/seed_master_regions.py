
import csv
import io
import os
import re
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

POSTAL_CSV_URL = (
    "https://raw.githubusercontent.com/"
    "mrayhanfadil/kodepos-id/main/kodepos_parsed.csv"
)

DB_CONFIG = {
    "host": os.getenv("DB_HOST", "localhost"),
    "port": os.getenv("DB_PORT", "5432"),
    "dbname": os.getenv("DB_NAME"),
    "user": os.getenv("DB_USER"),
    "password": os.getenv("DB_PASS"),
}

BATCH_SIZE = 1000

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


def normalize_code(value):
    """
    Normalize administrative codes.
    Example:
        35.07.32.2001 -> 3507322001
        3507322001    -> 3507322001
    """
    return re.sub(r"[^0-9]", "", str(value or "").strip())


def validate_unique(records, key, label):
    values = [record[key] for record in records]

    if len(values) != len(set(values)):
        raise RuntimeError(
            f"Duplicate {label} codes found in dataset"
        )


def fetch_release():
    print("[INFO] Fetching administrative dataset metadata...")

    response = requests.get(
        GITHUB_API_URL,
        timeout=30
    )
    response.raise_for_status()

    release = response.json()

    print(f"[INFO] Dataset version: {release['tag_name']}")

    return release


def fetch_release_csv(release, filename):
    asset = next(
        (
            item
            for item in release.get("assets", [])
            if item["name"].lower() == filename.lower()
        ),
        None
    )

    if asset is None:
        raise RuntimeError(
            f"{filename} not found in GitHub release"
        )

    response = requests.get(
        asset["browser_download_url"],
        timeout=60
    )
    response.raise_for_status()

    reader = csv.DictReader(
        io.StringIO(
            response.content.decode("utf-8-sig")
        )
    )

    if not reader.fieldnames:
        raise RuntimeError(
            f"Invalid CSV structure: {filename}"
        )

    return reader.fieldnames, list(reader)


def fetch_postal_data():
    print("[INFO] Fetching postal code dataset...")

    response = requests.get(
        POSTAL_CSV_URL,
        timeout=120
    )
    response.raise_for_status()

    reader = csv.DictReader(
        io.StringIO(
            response.content.decode("utf-8-sig")
        )
    )

    required_columns = {
        "kode_pos",
        "full_code",
    }

    if not reader.fieldnames or not required_columns.issubset(
        reader.fieldnames
    ):
        raise RuntimeError(
            f"Unexpected postal CSV structure: {reader.fieldnames}"
        )

    records = []

    for row in reader:
        village_code = normalize_code(
            row.get("full_code")
        )

        postal_code = (
            row.get("kode_pos") or ""
        ).strip()

        if not village_code or not postal_code:
            continue

        if not re.fullmatch(r"\d{5}", postal_code):
            continue

        records.append({
            "village_code": village_code,
            "postal_code": postal_code,
        })

    if not records:
        raise RuntimeError(
            "Postal dataset contains no valid records"
        )

    print(
        f"[INFO] Retrieved {len(records)} postal mappings"
    )

    return records


def fetch_all_data():
    release = fetch_release()

    district_headers, district_rows = fetch_release_csv(
        release,
        "districts.csv"
    )

    village_headers, village_rows = fetch_release_csv(
        release,
        "villages.csv"
    )

    district_parent_column = next(
        (
            column
            for column in (
                "regency_code",
                "city_code",
            )
            if column in district_headers
        ),
        None
    )

    village_parent_column = next(
        (
            column
            for column in (
                "district_code",
            )
            if column in village_headers
        ),
        None
    )

    if not {
        "code",
        "name",
    }.issubset(district_headers) or not district_parent_column:
        raise RuntimeError(
            f"Unexpected districts.csv structure: {district_headers}"
        )

    if not {
        "code",
        "name",
    }.issubset(village_headers) or not village_parent_column:
        raise RuntimeError(
            f"Unexpected villages.csv structure: {village_headers}"
        )

    districts = []

    for row in district_rows:
        districts.append({
            "code": normalize_code(row.get("code")),
            "city_code": normalize_code(
                row.get(district_parent_column)
            ),
            "name": (row.get("name") or "").strip(),
        })

    subdistricts = []

    for row in village_rows:
        subdistricts.append({
            "code": normalize_code(row.get("code")),
            "district_code": normalize_code(
                row.get(village_parent_column)
            ),
            "name": (row.get("name") or "").strip(),
        })

    for record in districts:
        if not all(record.values()):
            raise RuntimeError(
                f"Invalid district record: {record}"
            )

    for record in subdistricts:
        if not all(record.values()):
            raise RuntimeError(
                f"Invalid subdistrict record: {record}"
            )

    validate_unique(
        districts,
        "code",
        "district"
    )

    validate_unique(
        subdistricts,
        "code",
        "subdistrict"
    )

    print(
        f"[INFO] Retrieved {len(districts)} districts"
    )

    print(
        f"[INFO] Retrieved {len(subdistricts)} subdistricts"
    )

    postal_records = fetch_postal_data()

    return districts, subdistricts, postal_records

def get_code_map(cursor, table_name):
    """
    Retrieve administrative code -> internal database ID.
    Table names are fixed internal constants, not user input.
    """

    cursor.execute(
        f"SELECT id, code FROM {table_name}"
    )

    return {
        normalize_code(code): record_id
        for record_id, code in cursor.fetchall()
    }


def get_postal_map(cursor):
    cursor.execute(
        """
        SELECT id, postal_code
        FROM master_postal_codes
        """
    )

    return {
        str(postal_code): record_id
        for record_id, postal_code in cursor.fetchall()
    }


def get_count(cursor, table_name):
    cursor.execute(
        f"SELECT COUNT(*) FROM {table_name}"
    )

    return cursor.fetchone()[0]


def insert_records(cursor, query, records):
    if not records:
        return

    execute_values(
        cursor,
        query,
        records,
        page_size=BATCH_SIZE
    )


def seed_districts(cursor, districts, now):
    print("\n[INFO] Seeding master_districts...")

    city_map = get_code_map(
        cursor,
        "master_cities"
    )

    missing_cities = {
        district["city_code"]
        for district in districts
        if district["city_code"] not in city_map
    }

    if missing_cities:
        raise RuntimeError(
            "Missing city codes in database: "
            + ", ".join(sorted(missing_cities))
        )

    records = [
        (
            city_map[district["city_code"]],
            district["code"],
            district["name"],
            now,
            now,
        )
        for district in districts
    ]

    before = get_count(
        cursor,
        "master_districts"
    )

    query = """
        INSERT INTO master_districts
            (city_id, code, name, created_at, updated_at)
        VALUES %s
        ON CONFLICT (code) DO NOTHING;
    """

    insert_records(
        cursor,
        query,
        records
    )

    after = get_count(
        cursor,
        "master_districts"
    )

    print(
        f"[INFO] Newly inserted districts: {after - before}"
    )

    print(
        f"[INFO] Total districts: {after}"
    )


def seed_subdistricts(cursor, subdistricts, now):
    print("\n[INFO] Seeding master_subdistricts...")

    district_map = get_code_map(
        cursor,
        "master_districts"
    )

    missing_districts = {
        village["district_code"]
        for village in subdistricts
        if village["district_code"] not in district_map
    }

    if missing_districts:
        raise RuntimeError(
            "Missing district codes in database: "
            + ", ".join(sorted(missing_districts))
        )

    records = [
        (
            district_map[village["district_code"]],
            village["code"],
            village["name"],
            now,
            now,
        )
        for village in subdistricts
    ]

    before = get_count(
        cursor,
        "master_subdistricts"
    )

    query = """
        INSERT INTO master_subdistricts
            (district_id, code, name, created_at, updated_at)
        VALUES %s
        ON CONFLICT (code) DO NOTHING;
    """

    insert_records(
        cursor,
        query,
        records
    )

    after = get_count(
        cursor,
        "master_subdistricts"
    )

    print(
        f"[INFO] Newly inserted subdistricts: {after - before}"
    )

    print(
        f"[INFO] Total subdistricts: {after}"
    )


def seed_postal_codes(
    cursor,
    postal_records,
    now
):
    print("\n[INFO] Seeding master_postal_codes...")

    subdistrict_map = get_code_map(
        cursor,
        "master_subdistricts"
    )

    postal_codes = sorted({
        record["postal_code"]
        for record in postal_records
        if record["village_code"] in subdistrict_map
    })

    before = get_count(
        cursor,
        "master_postal_codes"
    )

    postal_values = [
        (
            postal_code,
            now,
            now,
        )
        for postal_code in postal_codes
    ]

    query = """
        INSERT INTO master_postal_codes
            (postal_code, created_at, updated_at)
        VALUES %s
        ON CONFLICT (postal_code) DO NOTHING;
    """

    insert_records(
        cursor,
        query,
        postal_values
    )

    after = get_count(
        cursor,
        "master_postal_codes"
    )

    print(
        f"[INFO] Newly inserted postal codes: {after - before}"
    )

    print(
        f"[INFO] Total postal codes: {after}"
    )

    postal_map = get_postal_map(cursor)

    relations = {
        (
            subdistrict_map[record["village_code"]],
            postal_map[record["postal_code"]],
        )
        for record in postal_records
        if (
            record["village_code"] in subdistrict_map
            and record["postal_code"] in postal_map
        )
    }

    cursor.execute(
        """
        CREATE UNIQUE INDEX IF NOT EXISTS
        idx_master_subdistrict_postalcode_unique
        ON master_subdistrict_postal_codes
            (subdistrict_id, postal_code_id);
        """
    )

    before = get_count(
        cursor,
        "master_subdistrict_postal_codes"
    )

    relation_values = [
        (
            subdistrict_id,
            postal_code_id,
            now,
            now,
        )
        for subdistrict_id, postal_code_id in relations
    ]

    relation_query = """
        INSERT INTO master_subdistrict_postal_codes
            (
                subdistrict_id,
                postal_code_id,
                created_at,
                updated_at
            )
        VALUES %s
        ON CONFLICT (subdistrict_id, postal_code_id)
        DO NOTHING;
    """

    insert_records(
        cursor,
        relation_query,
        relation_values
    )

    after = get_count(
        cursor,
        "master_subdistrict_postal_codes"
    )

    print(
        f"[INFO] Newly inserted relationships: {after - before}"
    )

    print(
        f"[INFO] Total relationships: {after}"
    )

    matched_villages = {
        record["village_code"]
        for record in postal_records
        if record["village_code"] in subdistrict_map
    }

    unmatched_villages = {
        record["village_code"]
        for record in postal_records
        if record["village_code"] not in subdistrict_map
    }

    print(
        f"[INFO] Villages matched with postal data: "
        f"{len(matched_villages)}"
    )

    print(
        f"[WARNING] Unmatched postal village codes: "
        f"{len(unmatched_villages)}"
    )

    if unmatched_villages:
        sample = sorted(unmatched_villages)[:10]

        print(
            f"[WARNING] Sample unmatched codes: {sample}"
        )

def seed_all(
    districts,
    subdistricts,
    postal_records
):
    connection = None

    try:
        print("\n[INFO] Connecting to PostgreSQL...")

        connection = psycopg2.connect(**DB_CONFIG)

        with connection:
            with connection.cursor() as cursor:

                now = datetime.now(timezone.utc)

                seed_districts(
                    cursor,
                    districts,
                    now
                )

                seed_subdistricts(
                    cursor,
                    subdistricts,
                    now
                )

                seed_postal_codes(
                    cursor,
                    postal_records,
                    now
                )

        print("\n[SUCCESS] Master region seeding completed!")

    except Exception as error:
        print(
            f"\n[ERROR] Master region seeding failed: {error}"
        )
        raise

    finally:
        if connection:
            connection.close()

def main():
    try:
        validate_config()

        districts, subdistricts, postal_records = (
            fetch_all_data()
        )

        seed_all(
            districts,
            subdistricts,
            postal_records
        )

    except Exception as error:
        print(f"\n[FATAL] {error}")
        sys.exit(1)


if __name__ == "__main__":
    main()


# Run:
# python scripts/seed_master_regions.py
