package models

import (
	"time"
)

// ==========================================================
// 1. CONTA E SEGURANÇA (Privado)
// ==========================================================

// Account: Guarda as credenciais e o acesso global.
// USO: Apenas o serviço de Auth deve ler/escrever aqui. Nunca retorne isso num JSON público.
type Account struct {
	ID                 int64      `json:"-" db:"id" gorm:"primaryKey"`
	Email              string     `json:"email" db:"email" gorm:"uniqueIndex"`
	HashedPassword     string     `json:"-" db:"hashed_password"`
	
	// Tokens para manter a sessão ou recuperar senha
	Token              *string    `json:"-" db:"token"`
	TokenExpiresAt     *time.Time `json:"-" db:"token_expires_at"`
	
	// Nível de acesso ao painel administrativo (0=User, 9=Admin)
	AccessLevel        int8       `json:"access_level" db:"access_level"`
	
	// Se preenchido, o usuário está banido e não consegue logar.
	IDModerationBanned *int64     `json:"id_moderation_banned" db:"id_moderation_banned"`
	
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	UpdateNote         *string    `json:"-" db:"updated_note"`
}

// ==========================================================
// 2. PERFIL DO MAKER (Público + Analytics)
// ==========================================================

// Maker: É a "máscara" social do usuário.
type Maker struct {
	// PK compartilhada com Account (Relação 1:1)
	IDAccount int64   `json:"id" db:"id_account" gorm:"primaryKey"`

	// --- IDENTIDADE ---
	Nickname  string  `json:"nickname" db:"nickname" gorm:"uniqueIndex"` // @usuario
	Name      string  `json:"name" db:"name"`                            // Nome de exibição
	LastName  string  `json:"last_name" db:"last_name"`
	Bio       *string `json:"bio" db:"bio"`
	Website   *string `json:"website" db:"website"`

	// --- VISUAL ---
	AvatarURL        *string `json:"avatar_url" db:"avatar_url"`
	AvatarPreviewURL *string `json:"avatar_preview_url" db:"avatar_preview_url"`
	BannerURL        *string `json:"banner_url" db:"banner_url"`

	// --- GAMIFICAÇÃO (Progresso) ---
	// XP: Acumulado com ações (Upload, Like recebido). Define o Rank.
	XP     uint64 `json:"xp" db:"xp"`
	// Coins: Moeda de troca ganha em missões. Usada na Loja.
	Coins  uint64 `json:"coins" db:"coins"`
	// IDRank: Cache do nível atual (ex: "Veterano"). O Service atualiza isso quando o XP muda.
	IDRank int16  `json:"id_rank" db:"id_rank"`

	// --- GROWTH HACKING (Campos Artificiais) ---
	// USO: O Admin pode injetar números aqui para inflar a percepção de popularidade.
	// O Mapper 'ToPublicResponse' DEVE somar (Real + Artificial) antes de entregar pro Front.
	ArtificialXP        uint64 `json:"artificial_xp" db:"artificial_xp"`
	ArtificialFollowers uint64 `json:"artificial_followers" db:"artificial_followers"`
	ArtificialCreated   uint64 `json:"artificial_created" db:"artificial_created"` // Finge que criou mais itens

	// --- STATUS E PRIVILÉGIOS ---
	IsSuspended   int8 `json:"is_suspended" db:"is_suspended"` // 1=Shadowban (invisível), 2=Suspenso
	Verified      int8 `json:"verified" db:"verified"`         // 1=Verificado (Selo Azul)
	Vip           int8 `json:"vip" db:"vip"`                   // 1=Pro, 2=Legend (Acelerador de ganhos)
	
	// STATUS ESPECIAIS (O "Reconhecimento" que você pediu)
	// USO: Flags manuais para destacar usuários sem depender de algoritmo.
	IsPartner     bool `json:"is_partner" db:"is_partner"`         // Ex: "Parceiro Comercial"
	IsContributor bool `json:"is_contributor" db:"is_contributor"` // Ex: "Top Criador da Comunidade"

	// --- CUSTOMIZAÇÃO ATIVA (Skins) ---
	IDAvatarStyle  *int16 `json:"id_avatar_style" db:"id_avatar_style"`
	IDPackStyle    *int16 `json:"id_pack_style" db:"id_pack_style"`
	IDProfileStyle *int16 `json:"id_profile_style" db:"id_profile_style"`

	// Configurações Gerais
	IDSettings int64 `json:"id_settings" db:"id_settings"`

	// --- ANALYTICS DE CRIAÇÃO (O que ele FEZ) ---
	// USO: Incrementados quando ele cria/publica algo.
	// Serve para badges do tipo "Criador Prolífico" (Crie 100 packs).
	PacksCreatedCount    uint64 `json:"packs_created_count" db:"packs_created_count"`
	StickersCreatedCount uint64 `json:"stickers_created_count" db:"stickers_created_count"`
	
	// --- ANALYTICS DE MÉRITO (O que ele RECEBEU) ---
	// USO: Incrementados quando OUTROS interagem com o conteúdo dele.
	// Serve para ranking de popularidade.
	FollowersCount     uint64 `json:"followers_count" db:"followers_count"`
	LikesReceivedCount uint64 `json:"likes_received_count" db:"likes_received_count"`
	
	// --- ANALYTICS DE CONSUMO (O que ele USOU) ---
	// USO: Incrementados quando ele baixa/usa algo.
	// Serve para badges do tipo "Colecionador" ou "Fã".
	// NOTA: Respeita a privacidade pois é apenas um número, não diz QUAIS packs baixou.
	DownloadsPerformedCount uint64 `json:"downloads_performed_count" db:"downloads_performed_count"`
	FollowingsCount         uint64 `json:"followings_count" db:"followings_count"`
	LoginsStreak            int    `json:"logins_streak" db:"logins_streak"` // Dias seguidos

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ==========================================================
// 3. SETTINGS (Configurações)
// ==========================================================

type MakerSettings struct {
	ID        int64 `json:"id" db:"id" gorm:"primaryKey"`
	IDMaker   int64 `json:"id_maker" db:"id_maker"`

	// Privacidade
	IsProfilePrivate   bool `json:"is_profile_private" db:"is_profile_private"`
	// Backdoor: Permite ocultar a aba de "Meus Favoritos" para outros usuários.
	AreFavoritesPublic bool `json:"are_favorites_public" db:"are_favorites_public"`

	// Preferências
	ShowAdultContent bool `json:"show_adult_content" db:"show_adult_content"`
	AllowDirectMsg   bool `json:"allow_direct_msg" db:"allow_direct_msg"`

	// Notificações
	NotifyOnLike    bool `json:"notify_on_like" db:"notify_on_like"`
	NotifyOnFollow  bool `json:"notify_on_follow" db:"notify_on_follow"`
	NotifyOnNewPack bool `json:"notify_on_new_pack" db:"notify_on_new_pack"`

	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ==========================================================
// 4. BIBLIOTECA PESSOAL (Interações)
// ==========================================================

// MakerStickerFavorite: Stickers que o usuário salvou para usar.
type MakerStickerFavorite struct {
	IDMaker   int64 `json:"id_maker" db:"id_maker" gorm:"primaryKey"`
	IDSticker int64 `json:"id_sticker" db:"id_sticker" gorm:"primaryKey"`
	
	// Permite ordenar a coleção pessoal.
	Position  int16 `json:"position" db:"position"`
	
	// Backdoor Granular: "Quero que minha lista seja pública, MENOS esse sticker aqui".
	IsPublic  bool  `json:"is_public" db:"is_public"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// MakerPackFavorite: Packs que o usuário salvou.
type MakerPackFavorite struct {
	IDMaker   int64 `json:"id_maker" db:"id_maker" gorm:"primaryKey"`
	IDPack    int64 `json:"id_pack" db:"id_pack" gorm:"primaryKey"`
	Position  int16 `json:"position" db:"position"`
	IsPublic  bool  `json:"is_public" db:"is_public"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Likes (Sempre Públicos)
type MakerStickerLike struct {
	IDMaker   int64     `json:"id_maker" db:"id_maker" gorm:"primaryKey"`
	IDSticker int64     `json:"id_sticker" db:"id_sticker" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type MakerPackLike struct {
	IDMaker   int64     `json:"id_maker" db:"id_maker" gorm:"primaryKey"`
	IDPack    int64     `json:"id_pack" db:"id_pack" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type MakerFollow struct {
	IDMaker       int64     `json:"id_maker" db:"id_maker" gorm:"primaryKey"`
	IDMakerFollow int64     `json:"id_maker_follow" db:"id_maker_follow" gorm:"primaryKey"`
	Since         time.Time `json:"since" db:"since"`
}

// ==========================================================
// 5. DTOs & INPUTS (Incluindo Admin)
// ==========================================================

// MakerPublicResponse: A resposta pública (com os dados "maquiados").
type MakerPublicResponse struct {
	ID               int64     `json:"id"`
	Nickname         string    `json:"nickname"`
	Name             string    `json:"name"`
	Bio              *string   `json:"bio"`
	AvatarURL        *string   `json:"avatar_url"`
	AvatarPreviewURL *string   `json:"avatar_preview_url"`
	BannerURL        *string   `json:"banner_url"`
	Website          *string   `json:"website"`
	
	// Gamificação (Real + Artificial)
	XP               uint64    `json:"xp"`
	Rank             RankResponse `json:"rank"`
	
	// Status
	Verified         int8      `json:"verified"`
	Vip              int8      `json:"vip"`
	IsPartner        bool      `json:"is_partner"`
	IsContributor    bool      `json:"is_contributor"`
	
	// Métricas (Soma Real + Artificial)
	FollowersCount       uint64    `json:"followers_count"`
	FollowingsCount      uint64    `json:"followings_count"`
	PacksCreatedCount    uint64    `json:"packs_created_count"`
	StickersCreatedCount uint64    `json:"stickers_created_count"`
	
	// Visual
	AvatarStyle      *StyleResponse `json:"avatar_style"` 
	PackStyle        *StyleResponse `json:"pack_style"`
}

// AdminUpdateMakerInput: Ferramenta de "Deus" para manipular o perfil.
type AdminUpdateMakerInput struct {
	IDMaker int64 `json:"id_maker" binding:"required"`
	
	// Status Manual
	SetVip           *int8 `json:"set_vip"`
	SetVerified      *int8 `json:"set_verified"`
	SetPartner       *bool `json:"set_partner"`
	SetContributor   *bool `json:"set_contributor"`
	
	// Injeção de Números (Growth Hacking)
	// Valores aqui serão salvos em 'ArtificialXP', 'ArtificialFollowers', etc.
	SetArtificialXP        *uint64 `json:"set_artificial_xp"`
	SetArtificialFollowers *uint64 `json:"set_artificial_followers"`
	SetArtificialCreated   *uint64 `json:"set_artificial_created"`
	
	// Ações Punitivas
	SetSuspended     *int8 `json:"set_suspended"`
}

// Mapper Inteligente (Soma Real + Artificial)
func (m *Maker) ToPublicResponse(rankData RankResponse, avatarStyle, packStyle *StyleResponse) MakerPublicResponse {
	return MakerPublicResponse{
		ID:               m.IDAccount,
		Nickname:         m.Nickname,
		Name:             m.Name,
		Bio:              m.Bio,
		AvatarURL:        m.AvatarURL,
		AvatarPreviewURL: m.AvatarPreviewURL,
		BannerURL:        m.BannerURL,
		Website:          m.Website,
		
		// Soma dos valores reais com os artificiais
		XP:                   m.XP + m.ArtificialXP,
		FollowersCount:       m.FollowersCount + m.ArtificialFollowers,
		PacksCreatedCount:    m.PacksCreatedCount + m.ArtificialCreated,
		StickersCreatedCount: m.StickersCreatedCount + m.ArtificialCreated, // Simplificação: usa o mesmo boost
		
		FollowingsCount: m.FollowingsCount,
		
		Rank:          rankData,
		Verified:      m.Verified,
		Vip:           m.Vip,
		IsPartner:     m.IsPartner,
		IsContributor: m.IsContributor,
		
		AvatarStyle:   avatarStyle,
		PackStyle:     packStyle,
	}
}

/*
REGRAS DE NEGÓCIO - MAKER:

1.  Identidade Dual:
    - 'Account' é estritamente privada (Auth/Segurança).
    - 'Maker' é estritamente pública (Social/Gamificação).

2.  Analytics & Privacidade:
    - O sistema rastreia 'DownloadsPerformedCount' (quantos baixou) para fins de missões, mas NÃO rastreia QUAIS pacotes foram baixados. Isso protege o histórico do usuário.
    - 'LikesReceivedCount' aumenta quando alguém curte um pack/sticker deste Maker.
    - 'PacksCreatedCount' aumenta quando o Maker publica um novo pacote.

3.  Growth Hacking (Métodos Sujos):
    - O sistema permite injeção de dados via 'Artificial...' fields.
    - O Mapper 'ToPublicResponse' é o ÚNICO ponto de saída de dados para o front, e ele DEVE sempre somar (Real + Artificial).
    - Isso garante que o banco mantenha a integridade dos dados reais enquanto o marketing pode exibir números inflados.

4.  Reconhecimento (Fanático/Especialista):
    - Para identificar se um maker é "Fã de Naruto", o sistema NÃO usa colunas aqui.
    - O sistema usa a tabela 'MakerKeyword' (no arquivo keyword.go) para associar o Maker à tag "Naruto" com um peso alto.
    - Badges específicas ("Criador de Anime") são concedidas baseadas nessas keywords.

5.  Controle de Acesso:
    - Se 'IDModerationBanned' em Account estiver preenchido, o login é rejeitado no middleware.
    - Se 'IsSuspended' em Maker for true, o perfil fica visível mas com limitações (ex: não pode postar).
*/