# ADR-013: Catálogo remoto de instaladores (JSON + cache)

**Data**: 2026-04-04  
**Status**: Aceito  

## Contexto

O catálogo de pacotes instaláveis era mantido em código Go (lista em memória no repositório). Cada novo item ou correção exigia alterar o código e um novo release. O catálogo passou a viver em JSON versionado, embutido no binário com `go:embed` e espelhado por URL raw no GitHub, para que atualizações remotas não dependam de novo build para quem tem rede.

## Decisão

1. **Manifesto** em JSON com `schema_version` e array `packages` alinhado a `entities.Package` (campos em snake_case no JSON). O ficheiro canónico no repo é [`internal/infrastructure/catalog/installer-catalog.json`](../../../internal/infrastructure/catalog/installer-catalog.json); o mesmo conteúdo é **embutido no binário** com `go:embed` no pacote `catalog` (arranque sem rede e baseline alinhado ao build).
2. **URL padrão** (raw GitHub, branch `main`): mesmo caminho no repositório — `internal/infrastructure/catalog/installer-catalog.json` (ver constante `DefaultCatalogURL` em `internal/infrastructure/catalog`). **Override**: variável de ambiente `HOMESTEAD_CATALOG_URL`.
3. **Fetch não bloqueante** ao iniciar o TUI: `Init` dispara comando Bubble Tea que faz HTTP com timeout; a UI abre sem esperar a rede.
4. **Cadeia de dados** (por ordem): ao criar `NewInMemoryPackageRepository`, o repositório é preenchido a partir do JSON **embutido**; em `main`, **merge** do ficheiro de **cache** em disco (se válido e `schema_version` suportado), sobrepondo por `id`; no `Init` do TUI, fetch **remoto** em background — em sucesso, **merge por `id`** e **gravação do cache** com o corpo HTTP recebido.
5. **Caminho do cache**: `filepath.Join(os.UserCacheDir(), "homestead", "installer-catalog.json")` (convenção XDG respeitada via `UserCacheDir` onde aplicável).
6. **Compatibilidade de schema**: cliente suporta `schema_version == 1`. Se o servidor enviar versão maior, **ignorar** o payload remoto (manter embutido + cache anterior) até existir binário compatível.
7. **Categorias no TUI**: grupos do menu de instaladores permanecem definidos no código. Categorias desconhecidas no JSON são mapeadas para `other` e exibidas na secção **Outros**.

## Alternativas consideradas

- **Só catálogo embutido**: simples, mas cada alteração exige release (rejeitado).
- **Fetch bloqueante na abertura**: implementação mais simples, mas UX ruim com rede lenta ou DNS falho (rejeitado).
- **Assinatura criptográfica do manifesto**: adiada para fase futura (fora do âmbito inicial).

## Consequências

**Positivas**

- Novos pacotes e overrides podem ser publicados atualizando apenas o JSON no repositório (ou mirror).
- Utilizadores offline ou com fetch falho continuam com embutido + último cache válido.
- Merge por `id` permite sobrescrever metadados de pacotes embutidos sem duplicar toda a lista no JSON.

**Negativas**

- Dependência de rede e de disponibilidade do host do manifesto.
- Manifesto público expõe URLs e comandos de instalação (aceitável para o caso de uso).
- TUI passa a depender de um módulo de infraestrutura para fetch/parse (trade-off aceite neste CLI).

## Documentação relacionada

- Índice de ADRs: [README.md](README.md)
