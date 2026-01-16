# Otamaker API (Protótipo)

> [!WARNING]
> **Este projeto é um PROTÓTIPO.**  
> O foco deste código é demonstrar a **modelagem de banco de dados para alta performance** e a estruturação de uma **rede social de nicho (Anime)**. Muitas funcionalidades estão simplificadas ou hardcoded para fins de demonstração da arquitetura de dados.

## 🎯 Objetivo do Projeto
O **Otamaker** é uma API exploratória desenhada para sustentar uma rede social vertical focada em criar, compartilhar e colecionar stickers de anime.

O diferencial técnico deste projeto não é apenas o CRUD básico, mas sim as decisões de **Design de Banco de Dados** voltadas para escalabilidade e engajamento (Gamificação).

## 🏛️ Modelagem para Performance e Escala

A arquitetura do banco de dados (implementada em Go com GORM tags) prioriza a **leitura rápida** e a **segurança dos dados**. As principais estratégias adotadas foram:

### 1. Separação de Identidade (Privado vs Público)
Ao contrário de modelos tradicionais que misturam tudo em uma tabela `User`, separamos estritamente:
- **`Account`**: Contém dados sensíveis (Email, Senha Hash, Tokens). Só é acessível pelo serviço de Autenticação. NUNCA é exposto em endpoints públicos.
- **`Maker`**: É a "máscara social" do usuário. Contém apenas dados públicos (Nickname, Avatar, Bio). É otimizada para ser cacheada e lida milhares de vezes.

### 2. Denormalização Estratégica (Performance de Leitura)
Em redes sociais, *ler* contagens (likes, seguidores) é uma operação infinitamente mais frequente que *escrever*.
Para evitar `COUNT(*)` pesados em tabelas com milhões de linhas a cada request, utilizamos campos denormalizados nas próprias entidades:
- **No Anime**: `MakersCount`, `PacksCount`, `StickersCount`.
- **No Sticker**: `DownloadsCount`, `LikesCount`, `PacksCount`.
- **No Maker**: `FollowersCount`, `PacksCreatedCount`.
*O custo de escrita (atualizar +1) é pago para garantir leitura instantânea.*

### 3. Growth Hacking no Modelo de Dados
Como um protótipo de rede social, o sistema prevê mecanismos para "Impulsionar" perfis artificialmente sem corromper os dados reais (Auditabilidade).
- Os modelos possuem campos `ArtificialXP`, `ArtificialFollowers`, etc.
- A camada de apresentação (Mapper) soma automaticamente `Real + Artificial` antes de entregar o JSON ao frontend.
- Isso permite estratégias de marketing agressivas mantendo a integridade contábil do banco.

### 4. Otimização de Busca (JSONB/Arrays)
Para evitar tabelas de relacionamento complexas (N:N) em buscas simples, utilizamos tipos de dados avançados do PostgreSQL (suportados via Gorm):
- **Keywords e Emojis**: Armazenados como Arrays/JSONB diretamente no registro do Sticker/Anime.
- Isso permite indexação GIN e buscas textuais extremamente velozes sem múltiplos JOINs.

## 🧩 Ecossistema Social

O modelo reflete uma hierarquia clara de conteúdo:
1.  **Anime**: A fonte da verdade (Enciclopédia).
2.  **Pack**: Uma coleção curada de stickers.
3.  **Sticker**: A unidade atômica viral.
    *   *Feature Vital*: Um Sticker é independente. Ele pode ser marcado como `IsReusable` e ser adicionado a pacotes de *outros* Makers, gerando viralidade cruzada enquanto mantém o crédito ao autor original (`OriginalMakerID`).

## 🛠️ Stack Técnica (Contexto)
- **Linguagem**: Go (Golang)
- **ORM**: GORM (com tags customizadas para JSON e Arrays)
- **Banco de Dados Alvo**: PostgreSQL
- **Arquitetura**: Modular Monolith (focado em domínios: Makers, Animes, Stickers)

---
*Este repositório serve como documentação viva de uma arquitetura orientada a dados para redes sociais de alto tráfego.*
