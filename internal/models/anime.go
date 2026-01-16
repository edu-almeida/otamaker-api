package models

import (
	"otamaker-api/internal/constants"
	"time"
)

// ==========================================================
// 1. MODELO DE BANCO DE DADOS (Entidade)
// ==========================================================

type Anime struct {
	ID int64 `json:"id" db:"id" gorm:"primaryKey"`

	// FLAGS DE ESTADO E VISIBILIDADE
	IsAired     bool `json:"is_airing" db:"is_airing"`       // Se está em exibição atualmente.
	IsVisible   bool `json:"is_visible" db:"is_visible"`     // Controle soft de exibição.
	IsModerated bool `json:"is_moderated" db:"is_moderated"` // Bloqueio administrativo (Ban).
	IsFeatured  bool `json:"is_featured" db:"is_featured"`   // Destaque na Home.

	// IMAGENS
	ImageCoverURL        string `json:"image_cover_url" db:"image_cover_url"`
	ImageCoverPreviewURL string `json:"image_cover_preview_url" db:"image_cover_preview_url"`

	// INTERNACIONALIZAÇÃO (I18N)
	// Nomes e Sinopses em múltiplos idiomas (JSONB).
	Name     map[string]string `json:"name" db:"name" gorm:"serializer:json"`
	Synopsis map[string]string `json:"synopsis" db:"synopsis" gorm:"serializer:json"`

	// CLASSIFICAÇÃO E TAXONOMIA
	// Genres: Enum fixo (Ação, Comédia).
	Genres []constants.Genre `json:"id_genres" db:"id_genres" gorm:"type:text[];serializer:json"`
	// Season: Temporada de lançamento.
	Season constants.Season `json:"season" db:"season"`
	// Studios: Lista de estúdios produtores.
	Studios []string `json:"studios" db:"studios" gorm:"type:text[];serializer:json"`
	// Keywords: Cache de tags dinâmicas (ex: "cyberpunk", "time travel").
	Keywords []string `json:"keywords" db:"keywords" gorm:"type:text[];serializer:json"`

	// DADOS EXTERNOS E DATAS
	SourceScore *float32   `json:"source_score" db:"source_score"` // Nota do MAL/Anilist.
	FirstAired  *time.Time `json:"first_aired" db:"first_aired"`
	LastAired   *time.Time `json:"last_aired" db:"last_aired"`

	// DENORMALIZAÇÃO (Contadores de Cache)
	MakersCount         uint64 `json:"makers_count" db:"makers_count"`
	PacksCount          uint64 `json:"packs_count" db:"packs_count"`
	StickersCount       uint64 `json:"stickers_count" db:"stickers_count"`
	PacksDownloadsCount uint64 `json:"packs_downloads_count" db:"packs_downloads_count"`
	PacksLikesCount     uint64 `json:"packs_likes_count" db:"packs_likes_count"`

	// AUDITORIA
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at" db:"updated_at"`
	LastUpdateContext *string    `json:"last_update_context" db:"last_update_context"`
	IDModeration      *int64     `json:"id_moderation" db:"id_moderation"`
}

// ==========================================================
// 2. INPUTS (Payloads)
// ==========================================================

type CreateAnimeInput struct {
	Name                 map[string]string `json:"name" binding:"required"`
	Synopsis             map[string]string `json:"synopsis" binding:"required"`
	Genres               []string          `json:"genres" binding:"required,min=1"`
	Keywords             []string          `json:"keywords"`
	FirstAired           string            `json:"first_aired" binding:"required,datetime=2006-01-02"`
	LastAired            string            `json:"last_aired" binding:"omitempty,datetime=2006-01-02"`
	ImageCoverURL        string            `json:"image_cover_url" binding:"required,url"`
	ImageCoverPreviewURL string            `json:"image_cover_preview_url" binding:"required,url"`
	IsAired              bool              `json:"is_airing"`
	IsVisible            bool              `json:"is_visible"`
	IsFeatured           bool              `json:"is_featured"`
}

type UpdateAnimeInput struct {
	Name                 map[string]string `json:"name"`
	Synopsis             map[string]string `json:"synopsis"`
	Genres               []string          `json:"genres"`
	Keywords             []string          `json:"keywords"`
	Season               constants.Season  `json:"season"`
	SourceScore          *float32          `json:"source_score"`
	Studios              *[]string         `json:"studios"`
	FirstAired           *string           `json:"first_aired" binding:"omitempty,datetime=2006-01-02"`
	LastAired            *string           `json:"last_aired" binding:"omitempty,datetime=2006-01-02"`
	ImageCoverURL        *string           `json:"image_cover_url" binding:"omitempty,url"`
	ImageCoverPreviewURL *string           `json:"image_cover_preview_url" binding:"omitempty,url"`
	IsAired              *bool             `json:"is_airing"`
	IsFeatured           *bool             `json:"is_featured"`
	IsVisible            *bool             `json:"is_visible"`
	IsModerated          *bool             `json:"is_moderated"`
	IDModeration         *int64            `json:"id_moderation"`
}

/*
REGRAS DO ANIME:
1.  Idioma e Nomes: O anime deve ter seu nome e sinopse suportando múltiplos idiomas (Map), permitindo adição fácil de novas traduções.
2.  Gêneros: Obrigatório ter pelo menos 1 gênero (Genre) associado.
3.  Conteúdo Associado: Deve rastrear o número de criadores, pacotes e stickers vinculados a ele (MakersCount, etc).
4.  Relevância: Para ser listado publicamente, o anime idealmente deve ter pelo menos 1 pacote de stickers associado.
5.  Qualidade de Pacotes: O sistema incentiva pacotes entre 20 a 30 stickers para melhor experiência do usuário final.
6.  Consistência Temática: O anime só deve ter relação com pacotes que contenham stickers do próprio anime.
7.  Destaque: Pode ser marcado como 'IsFeatured' para aparecer em banners e áreas nobres do app.
8.  Ocultação: Pode ser marcado como oculto ('IsVisible=false') sem ser deletado.
9. Moderação: Remoção ou bloqueio ('IsModerated') é ação exclusiva da moderação e exige um ID de registro (IDModeration).
10. Temporada Atual: O sistema deve identificar animes da temporada atual baseando-se na data 'FirstAired' e flag 'IsAired'.
11. Performance: Dados de contagem (Makers, Packs, Stickers) devem ser denormalizados na entidade Anime para leitura rápida.
12. Datas: Obrigatório ter data de lançamento (FirstAired). Encerramento é opcional (para animes em andamento).
13. Score: Deve armazenar a nota de fontes externas (SourceScore) para fins de ordenação por qualidade.
14. Força Maior: Animes moderados por DMCA ou infração grave tornam-se ocultos e bloqueados imediatamente.
15. Hierarquia: Um anime pode ter pacotes, mas um pacote não obrigatoriamente precisa ter um anime (embora recomendado).
16. API Pública: Animes ocultos ou moderados não devem retornar nas listagens padrão da API.
*/