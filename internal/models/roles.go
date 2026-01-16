package models

import "time"

// ==========================================================
// 1. CARGOS E PERMISSÕES (RBAC)
// ==========================================================

// Role: O nome do cargo (ex: "Moderador", "Suporte", "Admin").
type Role struct {
	ID          int16  `json:"id" db:"id" gorm:"primaryKey"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	
	// Hierarquia (0-100). Um cargo 50 não pode banir um cargo 80.
	Level       int8   `json:"level" db:"level"` 

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Permission: A lista técnica do que pode ser feito.
// Geralmente isso é pré-populado no banco e raramente muda.
type Permission struct {
	ID          int16  `json:"id" db:"id" gorm:"primaryKey"`
	Code        string `json:"code" db:"code"` // Ex: "MOD_BAN_USER", "CONTENT_EDIT"
	Description string `json:"description" db:"description"`
}

// RolePermission: Define o que cada cargo pode fazer.
// Ex: Cargo "Moderador" (1) tem Permissão "Banir" (5).
type RolePermission struct {
	IDRole       int16 `json:"id_role" db:"id_role" gorm:"primaryKey"`
	IDPermission int16 `json:"id_permission" db:"id_permission" gorm:"primaryKey"`
}

// AccountRole: Atribui cargos aos usuários.
// Um usuário pode ter múltiplos cargos (ex: "VIP" e "Moderador").
type AccountRole struct {
	IDAccount int64 `json:"id_account" db:"id_account" gorm:"primaryKey"`
	IDRole    int16 `json:"id_role" db:"id_role" gorm:"primaryKey"`
	
	GrantedBy int64     `json:"granted_by" db:"granted_by"` // Quem deu esse cargo? (Auditoria)
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ==========================================================
// 2. CONSTANTES DE PERMISSÃO (Para uso no código)
// ==========================================================
// Em vez de colunas booleanas (CanBan, CanEdit...), usamos códigos.
// Isso permite adicionar novas permissões sem alterar a estrutura do banco.

const (
	PermCanBanUser          = "USER_BAN"
	PermCanViewPrivate      = "VIEW_PRIVATE"
	PermCanEditContent      = "CONTENT_EDIT"
	PermCanDeleteContent    = "CONTENT_DELETE"
	PermCanManageRoles      = "ROLE_MANAGE"
	PermCanViewReports      = "REPORT_VIEW"
	PermCanResolveReports   = "REPORT_RESOLVE"
	PermCanBoostContent     = "GROWTH_BOOST" // Acesso aos métodos sujos/artificiais
)