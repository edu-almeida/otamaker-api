package models

import "time"

// ==========================================================
// 1. ENTIDADE DE MODERAÇÃO (O Ticket)
// ==========================================================

type Moderation struct {
	ID int64 `json:"id" db:"id" gorm:"primaryKey"`

	// O QUE foi denunciado?
	IDTarget   int64      `json:"id_target" db:"id_target"`
	TargetType TargetType `json:"target_type" db:"target_type"`

	// QUEM denunciou?
	IDMakerReporter *int64 `json:"id_maker_reporter" db:"id_maker_reporter"` // Null = Sistema ou Anônimo

	// POR QUE?
	Reason      ReasonType `json:"reason" db:"reason"`
	Description string     `json:"description" db:"description"`   // Texto livre do usuário
	SnapshotURL *string    `json:"snapshot_url" db:"snapshot_url"` // Print ou estado do objeto na hora

	// RESOLUÇÃO
	IDMakerResolver *int64     `json:"id_maker_resolver" db:"id_maker_resolver"`
	Status          ModStatus  `json:"status" db:"status"`
	ActionTaken     ModAction  `json:"action_taken" db:"action_taken"`
	ModNote         *string    `json:"mod_note" db:"mod_note"` // Justificativa interna
	ResolvedAt      *time.Time `json:"resolved_at" db:"resolved_at"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ==========================================================
// 2. ENUMS E TIPOS AUXILIARES
// ==========================================================

// TargetType: O que está sendo denunciado?
type TargetType int8

const (
	TargetMaker   TargetType = 1 // Perfil (Avatar, Bio, Nickname)
	TargetPack    TargetType = 2 // Pacote inteiro (Capa, Nome)
	TargetSticker TargetType = 3 // Imagem específica do Sticker
	TargetAnime   TargetType = 4 // Dados do Anime (Sinopse, Capa)
	TargetReview  TargetType = 5 // Comentários ou Avaliações
)

// ReasonType: Motivo Detalhado da Denúncia
type ReasonType int8

const (
	// --- COMPORTAMENTO E SEGURANÇA (Grave) ---
	ReasonSpam           ReasonType = 1 // Spam, Bots, Links maliciosos
	ReasonNudity         ReasonType = 2 // Pornografia, Hentai explícito (se proibido)
	ReasonViolence       ReasonType = 3 // Gore, Automutilação, Violência real
	ReasonHateSpeech     ReasonType = 4 // Racismo, Homofobia, Xenofobia
	ReasonHarassment     ReasonType = 5 // Bullying, Assédio, Stalking
	ReasonIllegalContent ReasonType = 6 // Drogas, Terrorismo, Crimes

	// --- INTEGRIDADE E AUTORIA (Maker/Pack) ---
	ReasonCopyright      ReasonType = 10 // Roubo de arte, DMCA
	ReasonImpersonation  ReasonType = 11 // Fingir ser outro Maker ou Staff (Fake)
	ReasonMisinformation ReasonType = 12 // Fake News ou dados falsos sobre animes

	// --- QUALIDADE E METADADOS (Sticker/Pack/Info) ---
	ReasonLowQuality     ReasonType = 20 // Imagem quebrada, resolução ruim, recorte mal feito
	ReasonWrongContext   ReasonType = 21 // Sticker no pacote errado / Anime errado
	ReasonMisleadingTags ReasonType = 22 // "Clickbait": Tag "Naruto" num pack de "One Piece"
	ReasonSpoiler        ReasonType = 23 // Spoiler de Anime sem aviso (Crítico para a comunidade!)
	ReasonOffTopic       ReasonType = 24 // Conteúdo que não é Anime/Geek (se for regra do app)

	ReasonOther ReasonType = 99
)

// ModStatus: Estado do fluxo
type ModStatus int8

const (
	StatusPending       ModStatus = 0
	StatusInvestigating ModStatus = 1 // Em análise (pode demorar)
	StatusResolved      ModStatus = 2 // Ação tomada (Ban/Hide)
	StatusRejected      ModStatus = 3 // Denúncia improcedente
	StatusIgnored       ModStatus = 4 // Spam de denúncia
)

// ModAction: A punição aplicada
type ModAction int8

const (
	ActionNone        ModAction = 0 // Nenhuma ação (Inocente)
	ActionWarning     ModAction = 1 // Envia alerta/aviso ao Maker
	ActionEditForce   ModAction = 2 // Admin editou/removeu apenas o conteúdo ofensivo (ex: removeu tags)
	ActionHideContent ModAction = 3 // Ocultou o item (Soft Ban)
	ActionBanContent  ModAction = 4 // Removeu o item e marcou infração (Hard Ban)
	ActionSuspendUser ModAction = 5 // Usuário não pode postar temporariamente
	ActionBanAccount  ModAction = 6 // Banimento total da conta
)

// ==========================================================
// 3. DTOs
// ==========================================================

type ModerationResponse struct {
	ID          int64  `json:"id"`
	TargetType  string `json:"target_type"` // Ex: "Sticker"
	TargetID    int64  `json:"target_id"`
	Reason      string `json:"reason"`      // Ex: "Spoiler sem aviso"
	Description string `json:"description"` // Texto do delator
	Status      string `json:"status"`

	CreatedAt  time.Time  `json:"created_at"`
	ResolvedAt *time.Time `json:"resolved_at"`

	// Preview do item denunciado para facilitar a vida do Admin no Front
	TargetPreview interface{} `json:"target_preview"`
}
