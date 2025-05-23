# Anki Flashcard Extractor API (Go + MVC)

Este projeto implementa uma API RESTful em Go para extrair flashcards de arquivos de coleção Anki (`.anki21`, `.anki2` ou `.apkg`). A API segue a arquitetura MVC (Model-View-Controller) para uma organização robusta e de fácil manutenção, utilizando o Gorilla Mux para roteamento HTTP.

## 🚀 Funcionalidades

* **Upload de Coleções Anki:** Permite o upload de arquivos `.anki21`, `.anki2` ou `.apkg` (pacotes Anki).
* **Extração Inteligente:**
    * Se for um `.apkg`, descompacta automaticamente para encontrar o arquivo de banco de dados (`.anki21` ou `.anki2`) interno.
    * **Prioriza sempre o formato `.anki21`** para garantir compatibilidade com as versões mais recentes do Anki e evitar problemas com notas de migração.
* **Extração de Flashcards:** Extrai notas (Notes), cartões (Cards) e informações de decks (Deck Names) do banco de dados do Anki.
* **Retorno JSON:** Retorna os flashcards extraídos em formato JSON, incluindo detalhes da nota, do cartão e o nome do deck associado.
* **Estrutura MVC:** Organização clara e modular do código, facilitando o desenvolvimento, teste e manutenção.

## 🏗️ Arquitetura

O projeto adota uma estrutura de pastas e organização de código baseada nos princípios da arquitetura MVC, com camadas adicionais para maior robustez:

```
my-anki-api/
├── cmd/                          # Ponto de entrada da aplicação
│   └── api/
│       └── main.go
│
├── internal/                     # Código privado da aplicação (lógica de negócio e camadas)
│   ├── config/                   # Configurações da aplicação
│   ├── models/                   # Definições das structs (Flashcard, Note, Card, Deck)
│   ├── repositories/             # Camada de abstração de dados (interage com o DB SQLite)
│   ├── services/                 # Lógica de negócio (extração, combinação de dados)
│   ├── handlers/                 # Controladores HTTP (lida com requisições, chama serviços)
│   ├── routes/                   # Definição das rotas da API
│   └── utils/                    # Funções utilitárias (descompactação, manipulação de strings)
│
├── tmp/                          # Diretório para uploads e descompactação temporária
├── go.mod                        # Dependências Go
├── go.sum
└── README.md
```

## ⚙️ Como Rodar

### Pré-requisitos

* **Go:** Versão 1.18 ou superior.
* **Coleção Anki:** Tenha um arquivo de coleção Anki (`.anki21`, `.anki2` ou `.apkg`) para testar. Recomenda-se exportar do **Anki Desktop 2.1.x ou superior** para garantir o formato `.anki21` dentro do `.apkg`.

### Passos

1.  **Clone o Repositório:**

    ```bash
    git clone https://github.com/seu-usuario/my-anki-api.git # Substitua pelo seu repositório
    cd my-anki-api
    ```

2.  **Baixe as Dependências:**

    ```bash
    go mod tidy
    ```

3.  **Crie o Diretório Temporário:**
    A aplicação tentará criar este diretório, mas é bom garantir.

    ```bash
    mkdir -p tmp
    ```

4.  **Execute a API:**

    ```bash
    go run ./cmd/api/main.go
    ```

    A API estará rodando em `http://localhost:8080`.

## 💻 Endpoints da API

### `POST /upload-anki`

Este endpoint permite que você faça o upload de um arquivo de coleção Anki e receba os flashcards extraídos em JSON.

* **Método:** `POST`
* **URL:** `http://localhost:8080/upload-anki`
* **Content-Type:** `multipart/form-data`
* **Campo do Formulário:**
    * `anki_file` (tipo `file`): O arquivo `.anki21`, `.anki2` ou `.apkg` da sua coleção Anki.

#### Exemplo de Requisição (usando `curl`):

```bash
curl -X POST \
  http://localhost:8080/upload-anki \
  -H 'Content-Type: multipart/form-data' \
  -F 'anki_file=@/caminho/para/seu/arquivo/minha_colecao.apkg'
```

**Substitua `/caminho/para/seu/arquivo/minha_colecao.apkg` pelo caminho absoluto ou relativo do seu arquivo Anki.**

#### Exemplo de Resposta (JSON):

```json
[
  {
    "note": {
      "id": 1678901234567,
      "guid": "abcdefGHIjK",
      "model": 1678901234000,
      "fields": [
        "Pergunta do Flashcard",
        "Resposta do Flashcard"
      ],
      "tags": [
        "Go",
        "Anki"
      ]
    },
    "card": {
      "id": 1678901234568,
      "note_id": 1678901234567,
      "deck_id": 1,
      "interval": 1,
      "due": 1716447600,
      "template_order": 0
    },
    "deck_name": "Minhas Notas de Estudo"
  },
  {
    "note": {
      "id": 1678901234569,
      "guid": "LmNoPqrStUv",
      "model": 1678901234000,
      "fields": [
        "Outra Pergunta?",
        "Outra Resposta!"
      ],
      "tags": [
        "Programação"
      ]
    },
    "card": {
      "id": 1678901234570,
      "note_id": 1678901234569,
      "deck_id": 123456789,
      "interval": 5,
      "due": 1716879600,
      "template_order": 0
    },
    "deck_name": "Estrutura de Dados"
  }
]
```

---

## 🤝 Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou pull requests.

## 📄 Licença

Este projeto está licenciado sob a licença MIT.
