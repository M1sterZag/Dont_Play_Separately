package core_igdb_provider

type Game struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	GameType  int    `json:"game_type"`
	Checksum  string `json:"checksum"`
	UpdatedAt int    `json:"updated_at"`
	Cover     *Image `json:"cover"`
}

type Image struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type Platform struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Abbreviation string `json:"abbreviation"`
	Checksum     string `json:"checksum"`
	UpdatedAt    int    `json:"updated_at"`
	PlatformLogo *Image `json:"platform_logo"`
}
