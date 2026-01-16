package models

import (
	"time"
)

// ==========================================================
// DEFINIÇÕES E TIPOS
// ==========================================================

// KeywordCategory define a natureza da palavra-chave.
// Serve para o algoritmo distinguir um Personagem de um Sentimento ou Evento.
type KeywordCategory string

const (
	CategoryGeneral   KeywordCategory = "general"   // Ex: engraçado, meme, reação
	CategoryCharacter KeywordCategory = "character" // Ex: naruto, goku, luffy
	CategoryWork      KeywordCategory = "work"      // Ex: one_piece (Nome da Obra/Franquia)
	CategoryArtist    KeywordCategory = "artist"    // Ex: pixel_art, watercolor (Estilo visual)
	CategoryEmotion   KeywordCategory = "emotion"   // Ex: sad, happy, angry
	CategoryEvent     KeywordCategory = "event"     // Ex: christmas, halloween
)

// ==========================================================
// 1. ENTIDADE PRINCIPAL (A Inteligência)
// ==========================================================

type Keyword struct {
	ID int64 `json:"id" db:"id" gorm:"primaryKey"`

	// Slug Universal (Inglês).
	// É a chave canônica. Ex: "naruto_uzumaki".
	// Todas as buscas em qualquer idioma convergem para este slug.
	Slug string `json:"slug" db:"slug" gorm:"uniqueIndex"`

	// Nome de Exibição (Padrão).
	Name string `json:"name" db:"name"`

	// Categoria (Vital para Filtros).
	// Permite buscas do tipo: "Mostre-me todos os PERSONAGENS (CategoryCharacter)".
	Category KeywordCategory `json:"category" db:"category"`

	// Sinônimos / Aliases (JSONB).
	// Lista de termos que o usuário pode digitar e que significam a mesma coisa.
	// Ex: ["naruto", "ninja loiro", "uzumaki", "hokage"].
	Aliases []string `json:"aliases" db:"aliases" gorm:"type:text[];serializer:json"`

	// Metadados Extras (JSONB).
	// Flexibilidade para guardar cor da tag na UI, ícone, link wiki, etc.
	Meta map[string]string `json:"meta" db:"meta" gorm:"serializer:json"`

	// Métricas de Relevância
	// UsageCount: Popularidade de uso por Makers.
	// SearchCount: Tendência de busca por Usuários.
	UsageCount  uint64 `json:"usage_count" db:"usage_count"`
	SearchCount uint64 `json:"search_count" db:"search_count"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// =================================================================
// 2. TABELAS DE CONEXÃO (O Grafo de Conhecimento)
// =================================================================
// Estas tabelas formam a inteligência do sistema, permitindo relacionar
// entidades diferentes através de interesses comuns.

// MakerKeyword: Perfil de Especialidade do Criador.
// Ex: Este Maker posta muito sobre "PixelArt" (Weight alto).
type MakerKeyword struct {
	IDMaker   int64 `json:"id_maker" db:"id_maker" gorm:"primaryKey"`
	IDKeyword int64 `json:"id_keyword" db:"id_keyword" gorm:"primaryKey"`
	
	// Peso (0-100). Relevância do maker neste tema.
	Weight    int       `json:"weight" db:"weight"` 
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// AnimeKeyword: Temas do Anime.
// Ex: "Cyberpunk", "Dystopian" (Além dos gêneros fixos).
type AnimeKeyword struct {
	IDAnime   int64     `json:"id_anime" db:"id_anime" gorm:"primaryKey"`
	IDKeyword int64     `json:"id_keyword" db:"id_keyword" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// PackKeyword: Temas do Pacote.
// Ex: "Natal", "Halloween". Agrupa pacotes de diferentes animes.
type PackKeyword struct {
	IDPack    int64     `json:"id_pack" db:"id_pack" gorm:"primaryKey"`
	IDKeyword int64     `json:"id_keyword" db:"id_keyword" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// StickerKeyword: Analytics de Stickers.
// Usado para processos em background. A leitura em tempo real usa o cache no Sticker.
type StickerKeyword struct {
	IDSticker int64     `json:"id_sticker" db:"id_sticker" gorm:"primaryKey"`
	IDKeyword int64     `json:"id_keyword" db:"id_keyword" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}