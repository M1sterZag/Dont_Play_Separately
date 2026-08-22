# Описание API

## Основные сущности:

### Users
- ID UUID
- version INT
- email TEXT
- hashed_password TEXT
- nickname TEXT
- description TEXT
- favorite_platform ID UUID
- avatar_link TEXT
- created_at TIME


### Teams
- ID UUID
- version INT
- platform_id ID UUID REF platform
- game_id ID UUID REF game
- is_rating BOOL
- desired_rating TEXT
- contact_link TEXT
- description TEXT
- amount_of_players INT
- created_at TIME
- owner_id ID UUID REF user

### Games
- ID UUID
- version INT
- title TEXT
- icon_link TEXT

