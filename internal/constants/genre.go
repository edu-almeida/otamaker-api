package constants

import (
	"fmt"
	"strings"
)

// Genre define o tipo para os gêneros.
type Genre string

// --- CHAVES DO BANCO (SLUGS EM INGLÊS) ---
// Valores em Inglês para padronização internacional do banco.
const (
	Acao          Genre = "action"
	Aventura      Genre = "adventure"
	Comedia       Genre = "comedy"
	Drama         Genre = "drama"
	Fantasia      Genre = "fantasy"
	Terror        Genre = "horror"
	Romance       Genre = "romance"
	FiccaoCien    Genre = "sci-fi"
	VidaCotidiana Genre = "slice_of_life"
	Esportes      Genre = "sports"
	Misterio      Genre = "mystery"
	Sobrenatural  Genre = "supernatural"
	Suspense      Genre = "thriller"
	Isekai        Genre = "isekai"
	Mecha         Genre = "mecha"
	Harem         Genre = "harem"
	Ecchi         Genre = "ecchi"
	Psicologico   Genre = "psychological"
	Escolar       Genre = "school"
	Militar       Genre = "military"
	Musical       Genre = "music"
	Historico     Genre = "historical"
	Shonen        Genre = "shonen"
	Shojo         Genre = "shojo"
	Seinen        Genre = "seinen"
	Josei         Genre = "josei"
	Kodomo        Genre = "kids"
)

// --- MAPA DE TRADUÇÕES ---
// Chave do mapa interno é sempre o Language (string(PT_BR), etc.).
var genreTranslations = map[Genre]map[string]string{
	Acao:          {string(PT_BR): "Ação", string(EN_US): "Action", string(ES_ES): "Acción"},
	Aventura:      {string(PT_BR): "Aventura", string(EN_US): "Adventure", string(ES_ES): "Aventura"},
	Comedia:       {string(PT_BR): "Comédia", string(EN_US): "Comedy", string(ES_ES): "Comedia"},
	Drama:         {string(PT_BR): "Drama", string(EN_US): "Drama", string(ES_ES): "Drama"},
	Fantasia:      {string(PT_BR): "Fantasia", string(EN_US): "Fantasy", string(ES_ES): "Fantasía"},
	Terror:        {string(PT_BR): "Terror", string(EN_US): "Horror", string(ES_ES): "Terror"},
	Romance:       {string(PT_BR): "Romance", string(EN_US): "Romance", string(ES_ES): "Romance"},
	FiccaoCien:    {string(PT_BR): "Ficção Científica", string(EN_US): "Sci-Fi", string(ES_ES): "Ciencia Ficción"},
	VidaCotidiana: {string(PT_BR): "Vida Cotidiana", string(EN_US): "Slice of Life", string(ES_ES): "Recuentos de la vida"},
	Esportes:      {string(PT_BR): "Esportes", string(EN_US): "Sports", string(ES_ES): "Deportes"},
	Misterio:      {string(PT_BR): "Mistério", string(EN_US): "Mystery", string(ES_ES): "Misterio"},
	Sobrenatural:  {string(PT_BR): "Sobrenatural", string(EN_US): "Supernatural", string(ES_ES): "Sobrenatural"},
	Suspense:      {string(PT_BR): "Suspense", string(EN_US): "Thriller", string(ES_ES): "Suspenso"},
	Isekai:        {string(PT_BR): "Isekai", string(EN_US): "Isekai", string(ES_ES): "Isekai"},
	Mecha:         {string(PT_BR): "Mecha", string(EN_US): "Mecha", string(ES_ES): "Mecha"},
	Harem:         {string(PT_BR): "Harem", string(EN_US): "Harem", string(ES_ES): "Harem"},
	Ecchi:         {string(PT_BR): "Ecchi", string(EN_US): "Ecchi", string(ES_ES): "Ecchi"},
	Psicologico:   {string(PT_BR): "Psicológico", string(EN_US): "Psychological", string(ES_ES): "Psicológico"},
	Escolar:       {string(PT_BR): "Escolar", string(EN_US): "School", string(ES_ES): "Escolar"},
	Militar:       {string(PT_BR): "Militar", string(EN_US): "Military", string(ES_ES): "Militar"},
	Musical:       {string(PT_BR): "Musical", string(EN_US): "Music", string(ES_ES): "Musical"},
	Historico:     {string(PT_BR): "Histórico", string(EN_US): "Historical", string(ES_ES): "Histórico"},
	Shonen:        {string(PT_BR): "Shounen", string(EN_US): "Shonen", string(ES_ES): "Shonen"},
	Shojo:         {string(PT_BR): "Shoujo", string(EN_US): "Shojo", string(ES_ES): "Shojo"},
	Seinen:        {string(PT_BR): "Seinen", string(EN_US): "Seinen", string(ES_ES): "Seinen"},
	Josei:         {string(PT_BR): "Josei", string(EN_US): "Josei", string(ES_ES): "Josei"},
	Kodomo:        {string(PT_BR): "Infantil", string(EN_US): "Kids", string(ES_ES): "Infantil"},
}

// Cache para normalização rápida (entrada cliente -> slug banco).
// Evita varrer todo o map genreTranslations a cada chamada.
var (
	// slug exato: "sci-fi" -> FiccaoCien
	genreSlugMap = func() map[string]Genre {
		m := make(map[string]Genre, len(genreTranslations))
		for g := range genreTranslations {
			m[string(g)] = g
		}
		return m
	}()

	// qualquer tradução em qualquer idioma: "ficção científica" -> FiccaoCien
	genreLowerMap = func() map[string]Genre {
		m := make(map[string]Genre, len(genreTranslations)*3)
		for g, translations := range genreTranslations {
			for _, text := range translations {
				m[strings.ToLower(text)] = g
			}
		}
		return m
	}()
)

// IsValid verifica se o gênero existe (segurança).
//
//	O(1) consultando o map central.
func (g Genre) IsValid() bool {
	_, exists := genreTranslations[g]
	return exists
}

// Translate traduz para o idioma do cliente.
// Usa Language.Get, que cuida de fallback e default.
func (g Genre) Translate(lang Language) string {
	if !g.IsValid() {
		return ""
	}

	trans, ok := genreTranslations[g]
	if !ok {
		// Fallback: "slice_of_life" -> "Slice Of Life"
		return strings.Title(strings.ReplaceAll(string(g), "_", " "))
	}

	return lang.Get(trans)
}

// NormalizeGenre converte Input -> Slug Inglês (Banco).
// Ex.: "Ação" / "acción" / "action" -> "action".
//
// Varria o map inteiro em dois for.
// Consulta em cache O(1) na maioria absoluta dos casos.
func NormalizeGenre(input string) (Genre, error) {
	cleanInput := strings.TrimSpace(strings.ToLower(input))

	// 1. Match slug exato ("sci-fi").
	if g, ok := genreSlugMap[cleanInput]; ok {
		return g, nil
	}

	// 2. Match qualquer tradução ("ficção científica", "ciencia ficción", "sci-fi").
	if g, ok := genreLowerMap[cleanInput]; ok {
		return g, nil
	}

	return "", fmt.Errorf("gênero não encontrado: %s", input)
}

// GetAvailableGenres retorna todos os gêneros conhecidos.
// Útil para popular SELECTs ou validações server-side.
func GetAvailableGenres() []Genre {
	list := make([]Genre, 0, len(genreTranslations))
	for g := range genreTranslations {
		list = append(list, g)
	}
	return list
}
