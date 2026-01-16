package constants

import "strings"

// Language é a chave do banco (snake_case).
type Language string

// Constantes Públicas (Exportadas)
const (
	PT_BR Language = "pt_br" // Default
	EN_US Language = "en_us"
	ES_ES Language = "es_es"

	Default = PT_BR
)

// Variáveis Privadas (Uso interno apenas)
var (
	// supported: Validação rápida O(1)
	supported = map[Language]bool{
		PT_BR: true,
		EN_US: true,
		ES_ES: true,
	}

	// parents: Mapeia prefixo (\"en\") para chave completa (\"en_us\")
	parents = map[string]Language{
		"pt": PT_BR,
		"en": EN_US,
		"es": ES_ES,
	}

	// OTIMIZAÇÃO 1: Cache pre-compilado dos ClientFormats
	// Evita converter strings em runtime toda vez que ListSupported() é chamado
	supportedList = func() []string {
		out := make([]string, 0, len(supported))
		for lang := range supported {
			out = append(out, lang.ClientFormat())
		}
		return out
	}()
)

// ==========================================================
// 1. ENTRADA (Normalize)
// ==========================================================

// NormalizeLanguage converte input do cliente (\"pt-BR\", \"en\") para chave do banco (\"pt_br\").
// Usa cache pre-processado, sem alocações desnecessárias.
func NormalizeLanguage(input string) Language {
	if input == "" {
		return Default
	}

	// OTIMIZAÇÃO 2: strings.Builder para manipulações (mais eficiente que ReplaceAll duplo)
	// ReplaceAll cria 2 strings temporárias; Builder cria só 1
	var sb strings.Builder
	sb.Grow(len(input))

	input = strings.TrimSpace(input)
	for _, ch := range strings.ToLower(input) {
		if ch == '-' {
			sb.WriteRune('_')
		} else {
			sb.WriteRune(ch)
		}
	}
	key := sb.String()

	// 2. Match Exato
	if supported[Language(key)] {
		return Language(key)
	}

	// 3. Match por Parente (\"es_mx\" -> \"es\")
	if parts := strings.SplitN(key, "_", 2); len(parts) > 0 {
		if match, ok := parents[parts[0]]; ok {
			return match
		}
	}

	return Default
}

// IsSupported verifica se o idioma é válido sem fazer fallback para o Default.
// Útil para rejeitar idiomas não suportados.
func IsSupported(input string) bool {
	lang := NormalizeLanguage(input)
	// Valida se não retornou Default
	return lang != Default || input == string(Default)
}

// ==========================================================
// 2. SAÍDA (Format)
// ==========================================================

// ClientFormat formata a chave do banco (\"pt_br\") para o padrão Web (\"pt-BR\").
// O(1) otimizado: evita alocações se não houver underscore.
func (l Language) ClientFormat() string {
	s := string(l)
	idx := strings.IndexByte(s, '_')

	// Sem underscore, retorna como está
	if idx == -1 {
		return s
	}

	// Com underscore: aloca apenas quando necessário
	// Exemplo: \"pt_br\" -> \"pt-BR\"
	return s[:idx] + "-" + strings.ToUpper(s[idx+1:])
}

// ListSupported retorna as strings formatadas para o frontend (ex: [\"pt-BR\", ...]).
// OTIMIZAÇÃO 3: Retorna cache pre-computado em init(), não recalcula a cada chamada
func ListSupported() []string {
	return supportedList
}

// ==========================================================
// 3. TEXTO (Get)
// ==========================================================

// (l Language) Get(content map[string]string)
// Fallback inteligente já otimizado.
func (l Language) Get(content map[string]string) string {
	if len(content) == 0 {
		return ""
	}

	key := string(l)

	// 1. Tenta exato (\"pt_br\")
	if val, ok := content[key]; ok && val != "" {
		return val
	}

	// 2. Tenta prefixo (\"pt\")
	if parts := strings.SplitN(key, "_", 2); len(parts) > 0 {
		if val, ok := content[parts[0]]; ok && val != "" {
			return val
		}
	}

	// 3. Tenta Default (\"pt_br\")
	if val, ok := content[string(Default)]; ok && val != "" {
		return val
	}

	// 4. Fallback: Primeiro valor encontrado
	for _, val := range content {
		if val != "" {
			return val
		}
	}

	return ""
}
