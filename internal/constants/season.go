package constants

import (
	"fmt"
	"strings"
	"time"
)

// Season define o tipo para as estações (Slug no Banco).
type Season string

// Chaves do Banco (Mantenha inglês, simples e sem acentos)
const (
	Spring Season = "spring"
	Summer Season = "summer"
	Autumn Season = "autumn"
	Winter Season = "winter"
)

// Mapa de Traduções (Centralizado)
// TRUQUE: Usamos as constantes de Language como chaves.
// Assim garantimos que "pt_br" aqui é igual ao "pt_br" da language.
var translations = map[Season]map[string]string{
	Spring: {
		string(PT_BR): "Primavera",
		string(EN_US): "Spring",
		string(ES_ES): "Primavera",
	},
	Summer: {
		string(PT_BR): "Verão",
		string(EN_US): "Summer",
		string(ES_ES): "Verano",
	},
	Autumn: {
		string(PT_BR): "Outono",
		string(EN_US): "Autumn",
		string(ES_ES): "Otoño",
	},
	Winter: {
		string(PT_BR): "Inverno",
		string(EN_US): "Winter",
		string(ES_ES): "Invierno",
	},
}

// Cache para normalização rápida (entrada cliente -> slug banco).
// Preenchido uma vez em init-time, O(1) em runtime.
var seasonLowerMap = func() map[string]Season {
	m := make(map[string]Season, len(translations)*4)
	for s, transMap := range translations {
		// Slug exato ("summer")
		m[string(s)] = s
		// Todas as traduções ("verão", "summer", "verano", ...)
		for _, text := range transMap {
			m[strings.ToLower(text)] = s
		}
	}
	return m
}()

// IsValid valida o slug do banco.
// Antes: switch. Agora: consulta em map, mais simples e escalável.
func (s Season) IsValid() bool {
	_, exists := translations[s]
	return exists
}

// GetSeasonByDate: O Sistema define a temporada (Lógica).
func GetSeasonByDate(date time.Time) Season {
	month := date.Month()
	switch month {
	case time.January, time.February, time.March:
		return Winter
	case time.April, time.May, time.June:
		return Spring
	case time.July, time.August, time.September:
		return Summer
	case time.October, time.November, time.December:
		return Autumn
	default:
		return Winter
	}
}

// Translate: SAÍDA (Banco -> Cliente).
// Usa a inteligência de Language.Get para fallback, default, etc.
func (s Season) Translate(lang Language) string {
	if !s.IsValid() {
		return ""
	}

	// Busca o mapa de textos dessa estação.
	t, ok := translations[s]
	if !ok {
		return string(s)
	}

	// DELEGA PARA Language.Get.
	return lang.Get(t)
}

// NormalizeSeason: ENTRADA (Filtro Cliente -> Banco).
// O cliente manda "Verão", a gente converte para "summer" para usar no SQL.
// Agora com cache: de O(n) varrendo o map para O(1).
func NormalizeSeason(input string) (Season, error) {
	input = strings.ToLower(strings.TrimSpace(input))

	if s, ok := seasonLowerMap[input]; ok {
		return s, nil
	}

	return "", fmt.Errorf("Season not found: %s", input)
}

// Helper extra: listar seasons existentes (se quiser usar em SELECTs).
func GetAvailableSeasons() []Season {
	list := make([]Season, 0, len(translations))
	for s := range translations {
		list = append(list, s)
	}
	return list
}
