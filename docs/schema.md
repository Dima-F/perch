# Схема бази даних

```mermaid
erDiagram
    water_body_type {
        int id PK
        text name
    }
    water_body {
        int id PK
        text name
        int water_body_type_id FK
    }
    locations {
        int id PK
        text name
        text region
        text notes
        int waterbody_id FK
    }
    fishing_types {
        int id PK
        text name
    }
    fish_species {
        int id PK
        text name
    }
    brands {
        int id PK
        text name
        text notes
    }
    lure_types {
        int id PK
        text name
        text notes
    }
    models {
        int id PK
        int brand_id FK
        int luretype_id FK
        text name
        text notes
    }
    lures {
        int id PK
        int model_id FK
        text color
        text size
        real weight_g
        text notes
    }
    fishing_sessions {
        int id PK
        datetime start_time
        datetime end_time
        text notes
    }
    session_fishing_types {
        int session_id FK
        int fishing_type_id FK
    }
    session_locations {
        int session_id FK
        int location_id FK
    }
    catches {
        int id PK
        int session_id FK
        int fish_id FK
        int lure_id FK
        int count
        real avg_length_cm
        real max_length_cm
        int weight_g
        real jig_weight_g
        text jig_setup
        text notes
    }
    blank_types {
        int id PK
        text name
        text notes
    }
    blanks {
        int id PK
        int brand_id FK
        int blank_type_id FK
        text name
        text casting
        text length
        text line
        text bought_at
        text notes
    }
    reels {
        int id PK
        int brand_id FK
        text name
        int reel_size
        text bearing_count
        text gear_rate
        text bought_at
        text notes
    }
    braided_lines {
        int id PK
        int brand_id FK
        text name
        text line_width
        real max_load
        text color
        int length
        text notes
    }
    spools {
        int id PK
        int reel_id FK
        int spool_number
        int size
        text notes
    }
    spool_braid {
        int id PK
        int spool_id FK
        int braid_id FK
    }

    water_body_type ||--o{ water_body       : "тип"
    water_body      ||--o{ locations        : "водойма"
    brands        |o--o{ models            : "бренд"
    lure_types    ||--o{ models            : "тип"
    models        ||--o{ lures             : "модель"
    fishing_sessions ||--o{ session_fishing_types : "тип"
    fishing_types    ||--o{ session_fishing_types : "тип"
    fishing_sessions ||--o{ session_locations  : "локація"
    locations        ||--o{ session_locations  : "локація"
    fishing_sessions ||--o{ catches         : "сесія"
    fish_species     ||--o{ catches         : "вид"
    lures            |o--o{ catches         : "приманка"
    brands        ||--o{ blanks             : "бренд"
    blank_types   ||--o{ blanks             : "тип"
    brands        |o--o{ reels              : "бренд"
    brands        ||--o{ braided_lines      : "бренд"
    reels         ||--o{ spools             : "котушка"
    spools        ||--o| spool_braid        : "шпуля"
    braided_lines ||--o| spool_braid        : "плетінка"
```

## Примітки до обмежень

| Таблиця | Обмеження |
|---|---|
| `models` | `UNIQUE(brand_id, name)` — унікальність у межах бренду |
| `catches.fish_id` | `ON DELETE RESTRICT` — не можна видалити вид риби, поки є улови |
| `catches.session_id` | `ON DELETE CASCADE` — улови видаляються разом із сесією |
| `catches.lure_id` | `ON DELETE SET NULL` — улов зберігається при видаленні приманки |
