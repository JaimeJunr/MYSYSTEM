# Documentação do projeto (índice central)

> **Objetivo**: Ponto de entrada único para navegar a documentação do Homestead. Documentos longos em inglês mantêm o idioma original; este índice e o [contexto de produto](product/CONTEXT.md) estão em português.

## Produto e negócio

- [Produto e contexto](product/CONTEXT.md) — para quem é, problema, valor, limites de escopo

## Índice por tema

### Arquitetura e design

- [Padrões de cima para baixo](architecture/VERTICAL_PATTERNS.md) — mapa vertical (TUI → serviços → domínio → infra)
- [Visão geral da arquitetura](architecture/ARCHITECTURE.md) — princípios, camadas, convenções
- [Layout do TUI](architecture/TUI_LAYOUT.md) — pacotes em `internal/tui`
- [Guia prático de padrões](architecture/PATTERNS_GUIDE.md) — exemplos de código
- [Diagramas](architecture/DIAGRAMS.md) — fluxos e módulos
- [ADRs (decisões)](architecture/adrs/README.md)

### Desenvolvimento

- [Implementação de instaladores](development/INSTALLER_IMPLEMENTATION.md)
- [Histórico de refatoração](development/REFACTORING_SUMMARY.md)
- [AGENTS.md](../AGENTS.md) — convenções para agentes e desenvolvimento no repositório

### Testes

- [Guia de testes](TESTING.md)

### Deploy e infraestrutura

- Release CI: `.github/workflows/release.yml` (artefatos `linux/amd64` e `linux/arm64`).

### APIs

- Não aplicável como API HTTP: CLI/TUI local. Integrações descritas na arquitetura e nos ADRs.

### Ferramentas

- [README do projeto](../README.md) — instalação, uso, estrutura de pastas
- `Makefile` — `build`, `run`, `test`, `test-integration`, `test-coverage`

## Estrutura da pasta `docs/`

```
docs/
├── INDEX.md                    # Este ficheiro
├── TESTING.md
├── product/
│   └── CONTEXT.md              # Produto e contexto de negócio
├── architecture/
│   ├── VERTICAL_PATTERNS.md    # Padrões top-down
│   ├── ARCHITECTURE.md
│   ├── TUI_LAYOUT.md
│   ├── PATTERNS_GUIDE.md
│   ├── DIAGRAMS.md
│   └── adrs/
└── development/
    ├── INSTALLER_IMPLEMENTATION.md
    └── REFACTORING_SUMMARY.md
```

## Por onde começar

| Perfil | Caminho |
| ------ | ------- |
| Novo no projeto | [README.md](../README.md) → [GETTING_STARTED.md](../GETTING_STARTED.md) → [product/CONTEXT.md](product/CONTEXT.md) |
| Arquitetura | [architecture/VERTICAL_PATTERNS.md](architecture/VERTICAL_PATTERNS.md) → [architecture/ARCHITECTURE.md](architecture/ARCHITECTURE.md) |
| Nova funcionalidade | Camadas em [ARCHITECTURE.md](architecture/ARCHITECTURE.md) → padrões em [PATTERNS_GUIDE.md](architecture/PATTERNS_GUIDE.md) → [TESTING.md](TESTING.md) |
| Decisões passadas | [architecture/adrs/README.md](architecture/adrs/README.md) |

## Métricas (referência rápida)

| Item | Ordem de grandeza |
| ---- | ----------------- |
| Pacotes de teste | vários sob `internal/` e raiz |
| Camadas | 4 (TUI, app, domain, infrastructure) |
| ADRs | ver [adrs/README.md](architecture/adrs/README.md) |

*Valores exatos (número de testes, entradas do catálogo) mudam com o tempo — confirme com `go test ./...` e o código-fonte.*

## Contribuindo na documentação

- Atualizar este índice quando criar documentos novos
- Um tema = um sítio canónico; evitar duplicar o mesmo conteúdo
- ADRs novos: [architecture/adrs](architecture/adrs/README.md)

---

**Última atualização:** 2026-04-04
