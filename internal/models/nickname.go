package models

import (
	"time"
)

// ==========================================================
// 1. TERMOS RESERVADOS (Blocklist)
// ==========================================================

type ReasonReserved int8

const (
	ReasonReservedSystem    ReasonReserved = 1
	ReasonReservedBrand     ReasonReserved = 2
	ReasonReservedOffensive ReasonReserved = 3
	ReasonReservedAnime     ReasonReserved = 4
	ReasonReservedOther     ReasonReserved = 5
)

type ReservedTerm struct {
	ID int64 `json:"id" db:"id" gorm:"primaryKey"`

	// Termo sempre salvo em lowercase.
	Term string `json:"term" db:"term" gorm:"uniqueIndex"`

	IsExactMatch bool           `json:"is_exact_match" db:"is_exact_match"`
	Reason       ReasonReserved `json:"reason" db:"reason"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ==========================================================
// 2. VERIFICAÇÃO DE DISPONIBILIDADE (DTOs)
// ==========================================================

// NicknameCheckInput: O que o Front envia enquanto o usuário digita.
type NicknameCheckInput struct {
	// REGRA ATUALIZADA: 'lowercase'
	// O nick DEVE ser a-z, 0-9. Sem letras maiúsculas.
	Nickname string `json:"nickname" binding:"required,min=3,max=30,alphanum,lowercase"`
}

// NicknameCheckResponse: O relatório completo sobre aquele nick.
type NicknameCheckResponse struct {
	Nickname string `json:"nickname"`

	// Resumo
	IsAvailable bool `json:"is_available"` // True = Pode usar!
	IsValid     bool `json:"is_valid"`     // True = Formato correto (a-z0-9)

	// Detalhes do Erro (se houver)
	Reason     string `json:"reason,omitempty"`     // "taken", "reserved", "invalid_format"
	Suggestion string `json:"suggestion,omitempty"` // Ex: "naruto_br", "naruto1"
}

// ==========================================================
// 3. REIVINDICAÇÃO (Legado/Suporte)
// ==========================================================
// Mantemos apenas DTOs de admin, já que o fluxo é via suporte/email.

type AdminForceNicknameInput struct {
	IDMaker        int64  `json:"id_maker" binding:"required"`
	NewNickname    string `json:"new_nickname" binding:"required,alphanum,lowercase"`
	IgnoreReserved bool   `json:"ignore_reserved"` // God Mode: Ignora a blocklist
}

/*
REGRAS DE NICKNAME:
1.  Validação Alfanumérica: Continua valendo. Apenas a-z, 0-9.
2.  Verificação de Reserva:
    - Antes de salvar um novo Maker, o sistema consulta 'ReservedTerm'.
    - Se encontrar match, bloqueia.
    - Exceção: Se quem está fazendo a alteração for um ADMIN, o sistema ignora essa tabela.
3.  Processo de Reivindicação:
    - Não existe tabela no banco. É feito via suporte (email).
    - O Admin usa o endpoint 'AdminUpdateMaker' para setar o nick manualmente, ignorando a trava.
*/
