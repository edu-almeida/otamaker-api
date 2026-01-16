package models

import "time"

// ==========================================================
// 1. ENTIDADE PRINCIPAL (O Pacote)
// ==========================================================

type Pack struct {
	ID int64 `json:"id" db:"id" gorm:"primaryKey"`

	// CONTEXTO E PROPRIEDADE
	// IDAnime: Vincula a um Anime (Contexto). Pode ser 0 se for Original.
	IDAnime int64 `json:"id_anime" db:"id_anime"`
	// IDMaker: O dono do pacote. Pode ser alterado (Transferência).
	IDMaker int64 `json:"id_maker" db:"id_maker"`
	// IDModerationBanned: Se preenchido, o pacote está banido globalmente.
	IDModerationBanned *int64 `json:"id_moderation_banned" db:"id_moderation_banned"`

	// DESCRIÇÃO E VISUAL
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	// TrayImageURL: Ícone que aparece na lista de packs do WhatsApp (96x96px).
	TrayImageURL string `json:"tray_image_url" db:"tray_image_url"`

	// INTELIGÊNCIA DE BUSCA
	// Keywords: Cache de slugs para busca rápida de temas de pacotes (ex: "natal", "memes").
	Keywords []string `json:"keywords" db:"keywords" gorm:"type:text[];serializer:json"`

	// CONFIGURAÇÕES
	IsDeleted  bool `json:"is_deleted" db:"is_deleted"`   // Soft Delete
	IsAnimated bool `json:"is_animated" db:"is_animated"` // Define se o pacote contém animações.
	IsFeatured bool `json:"is_featured" db:"is_featured"` // Destaque editorial.
	IsVisible  bool `json:"is_visible" db:"is_visible"`   // Publicado ou Rascunho.

	// MONETIZAÇÃO E SCORE
	Triage *float32 `json:"triage" db:"triage"` // Nota interna de qualidade.
	Score  *float32 `json:"score" db:"score"`   // Média de avaliação dos usuários.
	Price  *float64 `json:"price" db:"price"`   // Valor de venda (Null = Grátis).

	// MÉTRICAS E DADOS TÉCNICOS
	StickersCount  uint64  `json:"total_stickers" db:"total_stickers"`
	StickersSize   float64 `json:"stickers_size" db:"stickers_size"` // Validação de limite total do pacote.
	DataVersion    string  `json:"data_version" db:"data_version"`   // Hash para notificar atualização no WA.
	AvoidCache     bool    `json:"avoid_cache" db:"avoid_cache"`
	LikesCount     uint64  `json:"likes_count" db:"likes_count"`
	DownloadsCount uint64  `json:"downloads_count" db:"downloads_count"`
	FavoritesCount uint64  `json:"favorites_count" db:"favorites_count"`

	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at" db:"deleted_at"`
	LastUpdateContext string     `json:"last_update_context" db:"last_update_context"`
}

// NOTA: PackLike, PackFavorite e PackDownload foram removidos.
// Pertencem ao contexto do MAKER.

// ==========================================================
// 2. INPUTS (Criação e Atualização)
// ==========================================================

type PackCreate struct {
	IDAnime      int64    `json:"id_anime" binding:"required,gt=0"`
	IsAnimated   bool     `json:"is_animated"`
	IsVisible    bool     `json:"is_visible"`
	Price        float64  `json:"price" binding:"omitempty,gte=0"`
	TrayImageURL string   `json:"tray_image_url" binding:"required,url"`
	Name         string   `json:"name" binding:"required,min=3,max=64"`
	Description  string   `json:"description" binding:"omitempty,max=256"`
	Keywords     []string `json:"keywords" binding:"omitempty,max=10,dive,max=32"`

	// Lista inicial de Stickers. Mínimo 3 exigido pelo WA.
	Stickers []int64 `json:"stickers" binding:"required,min=3,max=30,dive,gt=0"`
}

type PackUpdate struct {
	// Permite transferir o pacote para outro Maker.
	IDMaker *int64 `json:"id_maker" binding:"omitempty,gt=0"`

	IDAnime      *int64    `json:"id_anime" binding:"omitempty,gt=0"`
	TrayImageURL *string   `json:"tray_image_url" binding:"omitempty,url"`
	IsAnimated   *bool     `json:"is_animated" binding:"omitempty"`
	Name         *string   `json:"name" binding:"omitempty,min=3,max=64"`
	Description  *string   `json:"description" binding:"omitempty,max=256"`
	Keywords     *[]string `json:"keywords" binding:"omitempty,max=10,dive,max=32"`

	// Atualiza a lista completa de stickers do pacote.
	Stickers *[]int64 `json:"updated_stickers" binding:"omitempty,min=3,max=30,dive,gt=0"`

	IsVisible *bool    `json:"is_visible" binding:"omitempty"`
	Price     *float64 `json:"price" binding:"omitempty,gte=0"`
}

/*
REGRAS DO PACOTE (PACK):
1.  Composição: Um pacote é um container lógico para 3 a 30 stickers. Ele não "contém" o arquivo do sticker, apenas a referência (ID).
2.  Transferência: Assim como stickers, pacotes podem ser transferidos entre Makers (Update IDMaker).
3.  Homogeneidade: A flag 'IsAnimated' define o comportamento do pacote. Misturar estáticos e animados pode causar rejeição em algumas plataformas.
4.  Metadados Técnicos: 'TrayImageURL' (96x96px) e 'DataVersion' são requisitos estritos para integração com APIs de mensageria (WhatsApp).
5.  Interações: Likes e Favoritos são ações do Maker armazenadas nas tabelas dele. O Pack guarda apenas os totais.
6.  Contexto: Um pacote deve tentar se vincular a um Anime (IDAnime), mas aceita vínculo genérico (ID=0) para conteúdos originais.
*/