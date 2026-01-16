package models

import (
	"time"
)

// ==========================================================
// 1. ENTIDADE PRINCIPAL (O Ativo/Conteúdo)
// ==========================================================

type Sticker struct {
	ID int64 `json:"id" db:"id" gorm:"primaryKey"`

	// PROPRIEDADE E PROVENIÊNCIA
	// IDMaker: O dono atual do sticker. Pode ser alterado (Transferência).
	IDMaker int64 `json:"id_maker" db:"id_maker"`
	// OriginalMakerID: O criador original. Nunca muda, garante o crédito moral/autoral.
	OriginalMakerID int64 `json:"original_maker_id" db:"original_maker_id"`
	// IDAnime: Opcional. Vincula o sticker a uma obra maior para contexto (ex: "Naruto").
	IDAnime *int64 `json:"id_anime" db:"id_anime"`

	// CONTROLE DE REUSO E VERSÃO
	// IsReusable: Se true, permite que OUTROS makers adicionem este sticker aos packs DELES.
	// O sticker aparece no pack de terceiros, mas a autoria (OriginalMakerID) continua sendo sua.
	IsReusable bool `json:"is_reusable" db:"is_reusable"`
	// ReplacesStickerID: Aponta para um sticker antigo que este substitui.
	// O sistema deve usar isso para redirecionar visualizações antigas para a nova versão remasterizada.
	ReplacesStickerID *int64 `json:"replaces_sticker_id" db:"replaces_sticker_id"`

	// ARQUIVOS E DIMENSÕES (Técnico)
	// ImageURL: URL pública do arquivo original (WebP/PNG/GIF).
	ImageURL string `json:"image_url" db:"image_url"`
	// ImageThumbURL: URL da miniatura otimizada para listagens mobile.
	ImageThumbURL string `json:"image_thumb_url" db:"image_thumb_url"`
	// Width/Height: Dimensões exatas em pixels. Obrigatório para o Frontend calcular
	// o layout Masonry (mosaico) antes da imagem carregar, evitando "pulos" na tela.
	Width  int `json:"width" db:"width"`
	Height int `json:"height" db:"height"`

	// INTELIGÊNCIA DE BUSCA
	// Emojis: Lista de códigos Unicode. Obrigatório para o teclado do WhatsApp sugerir o sticker.
	Emojis []string `json:"emojis" db:"emojis" gorm:"type:text[];serializer:json"`
	// Keywords: Cache de leitura rápida. Contém os slugs normalizados (ex: "sad", "naruto").
	// A inteligência real fica nas tabelas de Keyword, mas o sticker guarda essa cópia para busca veloz.
	Keywords []string `json:"keywords" db:"keywords" gorm:"type:text[];serializer:json"`

	// MÉTRICAS (Contadores Anônimos)
	// Não guardamos QUEM baixou, apenas QUANTAS vezes foi baixado.
	DownloadsCount uint64 `json:"downloads_count" db:"downloads_count"`
	// Total de likes que o sticker recebeu (agregado da tabela do Maker).
	LikesCount uint64 `json:"likes_count" db:"likes_count"`
	// Total de vezes que foi favoritado/salvo por usuários.
	FavoritesCount uint64 `json:"favorites_count" db:"favorites_count"`
	// Viralidade: Em quantos pacotes distintos este sticker está presente.
	PacksCount uint64 `json:"packs_count" db:"packs_count"`

	// VALIDAÇÃO TÉCNICA
	// Tamanho em bytes. Vital para garantir que não ultrapasse os limites rígidos do WhatsApp (ex: 500KB).
	SizeInBytes int64 `json:"size_in_bytes" db:"size_in_bytes"`

	// CONTROLE E MODERAÇÃO
	// IsModerated: True se o sticker foi banido por violar regras.
	IsModerated bool `json:"is_moderated" db:"is_moderated"`
	// IDModeration: Link para o ticket/registro da ação de moderação.
	IDModeration *int64 `json:"id_moderation" db:"id_moderation"`
	// IsVisible: O dono pode ocultar o sticker (Privado) sem deletar.
	IsVisible bool `json:"is_visible" db:"is_visible"`
	// IsDeleted: Soft Delete. O registro fica no banco para histórico/auditoria, mas some do app.
	IsDeleted bool `json:"is_deleted" db:"is_deleted"`
	// DeletedAt: Data da exclusão.
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ==========================================================
// 2. RELACIONAMENTO N:N (Sticker <-> Pack)
// ==========================================================

// PackSticker: A Tabela Pivot.
// É aqui que o Sticker (que é independente) se conecta a um Pacote (que é uma coleção).
type PackSticker struct {
	IDPack    int64 `json:"id_pack" db:"id_pack" gorm:"primaryKey"`
	IDSticker int64 `json:"id_sticker" db:"id_sticker" gorm:"primaryKey"`

	// Position: Define a ordem visual (0, 1, 2...) do sticker dentro DESTE pacote.
	// Permite que o dono do pacote reorganize os stickers como quiser (arrastar e soltar).
	Position int16 `json:"position" db:"position"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// NOTA: As tabelas StickerLike, StickerFavorite e StickerDownload foram removidas daqui.
// Elas pertencem ao contexto do MAKER ("Meus Likes", "Meus Favoritos").
// O Sticker armazena apenas os contadores agregados (LikesCount, etc).

// ==========================================================
// 3. DTOs (Data Transfer Objects)
// ==========================================================

type StickerResponse struct {
	ID             int64    `json:"id"`
	IDAnime        *int64   `json:"id_anime"`
	IDMaker        int64    `json:"id_maker"` // Dono atual
	ImageURL       string   `json:"image_url"`
	ImageThumbURL  string   `json:"image_thumb_url"`
	Width          int      `json:"width"`
	Height         int      `json:"height"`
	Emojis         []string `json:"emojis"`
	Keywords       []string `json:"keywords"`
	DownloadsCount uint64   `json:"downloads_count"`
	LikesCount     uint64   `json:"likes_count"`
	FavoritesCount uint64   `json:"favorites_count"`
	PacksCount     uint64   `json:"packs_count"`
	IsReusable     bool     `json:"is_reusable"`
}

type CreateStickerInput struct {
	// Pode ser nulo se for conteúdo original não atrelado a anime.
	IDAnime *int64 `json:"id_anime" binding:"omitempty,gt=0"`
	// Obrigatório. Quem está subindo o arquivo.
	IDMaker int64 `json:"id_maker" binding:"required,gt=0"`

	// URLs validadas. Geralmente geradas após upload para um bucket S3/GCS.
	ImageURL      string `json:"image_url" binding:"required,url"`
	ImageThumbURL string `json:"image_thumb_url" binding:"required,url"`

	// Dimensões são obrigatórias para performance de renderização no app.
	Width  int `json:"width" binding:"required,gt=0"`
	Height int `json:"height" binding:"required,gt=0"`

	// Regra de Negócio: WhatsApp exige pelo menos 1 emoji associado.
	Emojis []string `json:"emojis" binding:"required,min=1"`

	// Maker envia termos livres ["Naruto", "Chorando"].
	// O Service deve processar isso na tabela de Keywords e salvar os slugs normalizados ["naruto_uzumaki", "crying"].
	Keywords []string `json:"keywords"`

	// Define se nasce público para reuso (Default: true).
	IsReusable *bool `json:"is_reusable"`

	// Obrigatório para validar se cabe nos limites (ex: < 500KB).
	SizeInBytes int64 `json:"size_in_bytes" binding:"required,gt=0"`
}

type UpdateStickerInput struct {
	// Permite transferir a posse do sticker para outro Maker.
	IDMaker *int64 `json:"id_maker" binding:"omitempty,gt=0"`

	// Atualização de metadados de busca.
	Keywords *[]string `json:"keywords"`
	Emojis   *[]string `json:"emojis" binding:"omitempty,min=1"`

	// Controles de visibilidade e reuso.
	IsVisible  *bool `json:"is_visible"`
	IsReusable *bool `json:"is_reusable"`

	// Ações exclusivas de moderação.
	IsModerated  *bool  `json:"is_moderated"`
	IDModeration *int64 `json:"id_moderation" binding:"omitempty,gt=0"`

	// Link para versão corrigida.
	ReplacesStickerID *int64 `json:"replaces_sticker_id" binding:"omitempty,gt=0"`
}

type ReorderStickersInput struct {
	// Lista de IDs na nova ordem desejada.
	StickerIDs []int64 `json:"sticker_ids" binding:"required,min=1"`
}

// Mapper
func (s *Sticker) ToResponse() StickerResponse {
	emojis := s.Emojis
	if emojis == nil {
		emojis = []string{}
	}
	keywords := s.Keywords
	if keywords == nil {
		keywords = []string{}
	}

	return StickerResponse{
		ID:             s.ID,
		IDAnime:        s.IDAnime,
		IDMaker:        s.IDMaker,
		ImageURL:       s.ImageURL,
		ImageThumbURL:  s.ImageThumbURL,
		Width:          s.Width,
		Height:         s.Height,
		Emojis:         emojis,
		Keywords:       keywords,
		DownloadsCount: s.DownloadsCount,
		LikesCount:     s.LikesCount,
		FavoritesCount: s.FavoritesCount,
		PacksCount:     s.PacksCount,
		IsReusable:     s.IsReusable,
	}
}

/*
REGRAS GERAIS DO STICKER:
1.  Independência: O Sticker é uma entidade atômica e independente. Ele não pertence a um pacote, ele é referenciado por pacotes.
2.  Propriedade (Maker): Todo sticker tem um Dono Atual (IDMaker) e um Criador Original (OriginalMakerID).
3.  Transferência: A posse (IDMaker) pode ser transferida para outro usuário, mas o crédito original (OriginalMakerID) é imutável.
4.  Privacidade de Dados: Não armazenamos dados de usuários que baixaram stickers, apenas contadores agregados (DownloadsCount).
5.  Reuso: Se marcado como 'IsReusable', o sticker pode ser incluído em pacotes de outros Makers, aumentando sua viralidade (PacksCount).

REGRAS DE USO E REQUISITOS TÉCNICOS:
6.  Integração WhatsApp: É obrigatório associar pelo menos 1 Emoji (Unicode) para que o sticker apareça nas sugestões do teclado.
7.  Dimensões: O sistema deve exigir e armazenar Width e Height (px) para garantir layouts estáveis no mobile.
8.  Tamanho de Arquivo: O SizeInBytes deve ser validado na entrada para respeitar limites de plataforma (ex: WebP < 500KB).
9.  Busca: As Keywords salvas no sticker devem ser apenas Slugs Normalizados (inglês) gerados pelo sistema de inteligência, garantindo busca global.

REGRAS DE ORGANIZAÇÃO:
10. Ordenação em Packs: A relação Sticker-Pack (PackSticker) possui um campo 'Position' para permitir ordenação manual dentro do pacote.
11. Versionamento: O campo 'ReplacesStickerID' permite lançar correções de imagem sem perder as métricas do sticker original.
*/