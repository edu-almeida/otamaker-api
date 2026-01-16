package models

import "time"

// ==========================================================
// 1. RANKING (Nível do Usuário)
// ==========================================================

// Rank: A escada de evolução.
// Lógica: Service monitora 'Maker.XP'. Se XP > Rank.MinXP, atualiza 'Maker.IDRank'.
type Rank struct {
	ID       int16  `json:"id" db:"id" gorm:"primaryKey"`
	Name     string `json:"name" db:"name"`     // Ex: "Iniciante", "Lenda"
	MinXP    uint64 `json:"min_xp" db:"min_xp"` // Gatilho de subida
	IconURL  string `json:"icon_url" db:"icon_url"`
	ColorHex string `json:"color_hex" db:"color_hex"`
}

// ==========================================================
// 2. MISSÕES (Tarefas Ativas)
// ==========================================================

type MissionType string

const (
	MissionTypeLogin    MissionType = "login"    // Gatilho: Auth
	MissionTypeCreate   MissionType = "create"   // Gatilho: Upload
	MissionTypeShare    MissionType = "share"    // Gatilho: Share Button
	MissionTypeDownload MissionType = "download" // Gatilho: Download Action
	MissionTypeLike     MissionType = "like"     // Gatilho: Like Action
)

// Mission: Configuração de uma tarefa.
type Mission struct {
	ID          int64       `json:"id" db:"id" gorm:"primaryKey"`
	Title       string      `json:"title" db:"title"`
	Description string      `json:"description" db:"description"`
	Type        MissionType `json:"type" db:"type"` // Define qual evento do sistema dispara a checagem

	// Meta (ex: Criar 5 packs)
	TargetCount int `json:"target_count" db:"target_count"`

	// Prêmios
	RewardXP    int `json:"reward_xp" db:"reward_xp"`
	RewardCoins int `json:"reward_coins" db:"reward_coins"`

	// Regras
	IsDaily   bool   `json:"is_daily" db:"is_daily"`       // Se true, 'MakerMission' reseta em 24h
	IsOneTime bool   `json:"is_one_time" db:"is_one_time"` // Se true, só pode fazer uma vez na vida
	MinRankID *int16 `json:"min_rank_id" db:"min_rank_id"` // Filtro: Só para ranks altos
	IsVipOnly bool   `json:"is_vip_only" db:"is_vip_only"` // Filtro: Só para VIPs
}

// MakerMission: Estado da missão para um usuário específico.
type MakerMission struct {
	IDMaker   int64 `json:"id_maker" db:"id_maker" gorm:"primaryKey"`
	IDMission int64 `json:"id_mission" db:"id_mission" gorm:"primaryKey"`

	CurrentCount int        `json:"current_count" db:"current_count"` // Ex: Fez 2 de 5
	IsCompleted  bool       `json:"is_completed" db:"is_completed"`   // Se true, já pegou o prêmio
	CompletedAt  *time.Time `json:"completed_at" db:"completed_at"`
}

// ==========================================================
// 3. CONQUISTAS (Badges e Insígnias)
// ==========================================================

// Badge: Conquista automática baseada em estatísticas acumuladas.
// USO: O sistema compara 'Maker.PacksCreatedCount' >= 'Badge.RequirementValue'.
type Badge struct {
	ID          int64  `json:"id" db:"id" gorm:"primaryKey"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	IconURL     string `json:"icon_url" db:"icon_url"`

	// Regra Automática
	RequirementType  string `json:"req_type" db:"req_type"`   // Ex: "packs_created_count"
	RequirementValue int    `json:"req_value" db:"req_value"` // Ex: 100
}

// Insignia: Medalha especial concedida MANUALMENTE ou por EVENTOS.
// USO: Serve para "Creator Reconhecido", "Staff", "Vencedor Evento X".
// Não depende de contadores, depende de ação administrativa (Backdoor de Mérito).
type Insignia struct {
	ID          int64  `json:"id" db:"id" gorm:"primaryKey"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	IconURL     string `json:"icon_url" db:"icon_url"`
	Rarity      string `json:"rarity" db:"rarity"` // Common, Epic, Legendary
}

// Tabelas de Ligação (Inventário de Conquistas)
type MakerBadge struct {
	IDMaker  int64     `json:"id_maker" db:"id_maker" gorm:"primaryKey"`
	IDBadge  int64     `json:"id_badge" db:"id_badge" gorm:"primaryKey"`
	EarnedAt time.Time `json:"earned_at" db:"earned_at"`
}

type MakerInsignia struct {
	IDMaker    int64     `json:"id_maker" db:"id_maker" gorm:"primaryKey"`
	IDInsignia int64     `json:"id_insignia" db:"id_insignia" gorm:"primaryKey"`
	EarnedAt   time.Time `json:"earned_at" db:"earned_at"`
}

// ==========================================================
// 4. LOJA E ESTILO (Cosméticos)
// ==========================================================

type StyleType string

const (
	StyleAvatar  StyleType = "avatar_frame"
	StylePack    StyleType = "pack_card"
	StyleProfile StyleType = "profile_bg"
)

// MakerStyle: Item cosmético.
type MakerStyle struct {
	ID       int16     `json:"id" db:"id" gorm:"primaryKey"`
	Name     string    `json:"name" db:"name"`
	Type     StyleType `json:"type" db:"type"`
	AssetURL string    `json:"asset_url" db:"asset_url"`

	// Regras de Compra/Uso
	PriceCoins *int   `json:"price_coins" db:"price_coins"` // Custo em Coins
	MinRankID  *int16 `json:"min_rank_id" db:"min_rank_id"` // Precisa ser nível X
	IsVipOnly  bool   `json:"is_vip_only" db:"is_vip_only"` // Exclusivo VIP
}

type MakerUnlockedStyle struct {
	IDMaker    int64     `json:"id_maker" db:"id_maker" gorm:"primaryKey"`
	IDStyle    int16     `json:"id_style" db:"id_style" gorm:"primaryKey"`
	UnlockedAt time.Time `json:"unlocked_at" db:"unlocked_at"`
}

// DTOs
type RankResponse struct {
	Name     string `json:"name"`
	IconURL  string `json:"icon_url"`
	ColorHex string `json:"color_hex"`
}

type StyleResponse struct {
	Name     string `json:"name"`
	AssetURL string `json:"asset_url"`
}

/*
REGRAS DE GAMIFICAÇÃO:

1.  Reconhecimento Automático vs Manual:
    - 'Badges' são para esforço quantitativo (ex: Baixou 1000 vezes). O sistema dá sozinho observando os contadores do Maker.
    - 'Insignias' são para status qualitativo (ex: Creator Reconhecido, Parceiro). O Admin dá manualmente via 'Backdoor'.

2.  Economia (XP vs Coins):
    - XP é eterno e define o Rank.
    - Coins são consumíveis e compram Styles.
    - Ambos são gerados completando Missões.

3.  Missões:
    - São o motor de engajamento diário.
    - O backend deve ter listeners (ex: OnDownload, OnLike) que incrementam 'MakerMission.CurrentCount'.

4.  O Fator "Fanático":
    - Para reconhecer um usuário como "Fã de One Piece", o sistema monitora a tabela 'MakerKeyword' (definida em keyword.go).
    - Se o usuário cria/baixa muito conteúdo com a tag "one_piece", o peso dele nessa keyword sobe.
    - Isso pode desbloquear Badges específicas ou Insígnias customizadas no futuro.
*/
