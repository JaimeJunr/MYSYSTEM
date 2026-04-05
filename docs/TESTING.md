# 🧪 Testing Guide - Homestead

Documentação completa sobre os testes do projeto Homestead.

## 📋 Índice

- [Visão Geral](#visão-geral)
- [Estrutura dos Testes](#estrutura-dos-testes)
- [Executando Testes](#executando-testes)
- [Tipos de Testes](#tipos-de-testes)
- [Cobertura de Código](#cobertura-de-código)
- [Escrevendo Novos Testes](#escrevendo-novos-testes)
- [Boas Práticas](#boas-práticas)

## 📊 Visão Geral

O projeto Homestead possui uma suite completa de testes incluindo:
- ✅ Testes unitários para packages individuais
- ✅ Testes de integração entre componentes
- ✅ Benchmarks de performance
- ✅ Utilitários de teste reutilizáveis

### Estatísticas Atuais

- **Packages testados**: inclui `scripts`, `tui` (raiz + subpacotes sem `_test`), `testutil`, `services`, `entities`, etc. (`go test ./...`)
- **Testes unitários**: ~30 testes
- **Testes de integração**: 3 testes
- **Benchmarks**: 5 benchmarks
- **Cobertura**: Execute `make test-coverage` para ver

## 🗂️ Estrutura dos Testes

```
Homestead/
├── integration_test.go              # Testes de integração do sistema
├── internal/
│   ├── scripts/
│   │   ├── script.go
│   │   └── script_test.go          # Testes do gerenciador de scripts
│   ├── tui/
│   │   ├── model.go
│   │   ├── model_test.go           # Testes do pacote raiz tui
│   │   ├── cmds/                   # (sem testes dedicados)
│   │   ├── items/
│   │   ├── msg/
│   │   ├── theme/
│   │   └── sysurl/
│   └── testutil/
│       └── testutil.go              # Helpers e utilitários de teste
└── coverage.out                     # Relatório de cobertura (gerado)
```

## 🚀 Executando Testes

### Comandos Rápidos

```bash
# Executar todos os testes
make test

# Testes com output detalhado
make test-verbose

# Testes com cobertura
make test-coverage

# Apenas testes unitários (rápidos)
make test-short

# Apenas testes de integração
make test-integration

# Benchmarks de performance
make benchmark

# Cobertura em HTML (abre no navegador)
make test-coverage-html
```

### Comandos Go Diretos

```bash
# Todos os testes
go test ./...

# Com verbose
go test -v ./...

# Package específico
go test ./internal/scripts
go test ./internal/tui/...

# Com cobertura
go test -cover ./...

# Modo short (pula testes longos)
go test -short ./...

# Executar teste específico
go test -v -run TestGetAllScripts ./internal/scripts

# Benchmarks
go test -bench=. ./...
go test -bench=BenchmarkGetAllScripts ./internal/scripts
```

## 🔍 Tipos de Testes

### 1. Testes Unitários

Testam funções e métodos isoladamente.

**Exemplo - internal/scripts/script_test.go:**
```go
func TestGetAllScripts(t *testing.T) {
    scripts := GetAllScripts()

    if len(scripts) == 0 {
        t.Error("Expected at least one script")
    }
}
```

**Localização:**
- `internal/scripts/script_test.go` - Testes do gerenciador de scripts
- `internal/tui/model_test.go` - Testes da interface TUI

### 2. Testes de Integração

Testam interação entre componentes.

**Exemplo - integration_test.go:**
```go
func TestIntegration_FullWiring(t *testing.T) {
    // Ver integration_test.go — espelha cmd/homestead/main.go e usa tui.NewModel(...)
}
```

**Executar:** `make test-integration`

### 3. Benchmarks

Medem performance de operações críticas.

**Exemplo:**
```go
func BenchmarkGetAllScripts(b *testing.B) {
    for i := 0; i < b.N; i++ {
        GetAllScripts()
    }
}
```

**Executar:** `make benchmark`

**Output esperado:**
```
BenchmarkGetAllScripts-8        1000000    1234 ns/op    512 B/op    10 allocs/op
```

### 4. Testes com Table-Driven Pattern

Testes parametrizados para múltiplos casos.

**Exemplo:**
```go
func TestGetScriptsByCategory(t *testing.T) {
    tests := []struct {
        category    ScriptCategory
        expectedMin int
    }{
        {CategoryCleanup, 3},
        {CategoryMonitoring, 2},
    }

    for _, tt := range tests {
        t.Run(string(tt.category), func(t *testing.T) {
            scripts := GetScriptsByCategory(tt.category)
            // assertions...
        })
    }
}
```

## 📈 Cobertura de Código

### Ver Cobertura

```bash
# Gerar e ver resumo
make test-coverage

# Ver em HTML (recomendado)
make test-coverage-html
```

### Output Esperado

```
github.com/JaimeJunr/Homestead/internal/scripts    coverage: 85.2% of statements
github.com/JaimeJunr/Homestead/internal/tui        coverage: 72.1% of statements
```

### Meta de Cobertura

- **Mínimo**: 70% de cobertura
- **Ideal**: 80%+ de cobertura
- **Crítico**: 90%+ para código de segurança/execução

## ✍️ Escrevendo Novos Testes

### Template Básico

```go
package mypackage

import "testing"

func TestMyFunction(t *testing.T) {
    // Arrange - preparar dados
    input := "test"

    // Act - executar função
    result := MyFunction(input)

    // Assert - verificar resultado
    if result != "expected" {
        t.Errorf("Expected 'expected', got '%s'", result)
    }
}
```

### Usando Test Helpers

```go
import "github.com/JaimeJunr/Homestead/internal/testutil"

func TestWithHelpers(t *testing.T) {
    // Usar helpers do testutil
    testutil.AssertEqual(t, "expected", actual)
    testutil.AssertTrue(t, condition, "should be true")

    // Criar script temporário
    scriptPath := testutil.CreateMockScript(t)
    // scriptPath é automaticamente limpo após o teste
}
```

### Subtests

```go
func TestMyFeature(t *testing.T) {
    t.Run("success case", func(t *testing.T) {
        // teste de sucesso
    })

    t.Run("error case", func(t *testing.T) {
        // teste de erro
    })
}
```

### Pular Testes Condicionalmente

```go
func TestExpensiveOperation(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping expensive test in short mode")
    }
    // teste demorado
}
```

## 🎯 Boas Práticas

### ✅ Faça

1. **Use t.Helper()** em funções auxiliares
   ```go
   func assertValid(t *testing.T, value string) {
       t.Helper()  // Marca como helper
       if value == "" {
           t.Error("value is empty")
       }
   }
   ```

2. **Nomes descritivos** para testes
   ```go
   func TestGetScriptsByCategory_ReturnsOnlyMatchingScripts(t *testing.T)
   ```

3. **Table-driven tests** para múltiplos casos

4. **Cleanup automático**
   ```go
   tmpFile := createTempFile()
   t.Cleanup(func() { os.Remove(tmpFile) })
   ```

5. **Mensagens de erro claras**
   ```go
   t.Errorf("Expected %d scripts, got %d", expected, len(actual))
   ```

### ❌ Evite

1. **Testes dependentes** (um teste depende de outro)
2. **Estado compartilhado** entre testes
3. **Hardcoded paths** (use t.TempDir())
4. **Testes sem assertions** (sempre verifique algo)
5. **Testes muito longos** (divida em subtests)

## 🔧 Configuração de CI/CD

### GitHub Actions (exemplo)

```yaml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run tests
        run: make test-coverage

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

## 🐛 Debugging Testes

### Ver output detalhado

```bash
go test -v ./internal/scripts
```

### Executar apenas um teste

```bash
go test -v -run TestGetAllScripts ./internal/scripts
```

### Com debug logs

```bash
go test -v ./... 2>&1 | grep "your debug message"
```

### Usar delve (debugger)

```bash
dlv test ./internal/scripts -- -test.run TestGetAllScripts
```

## 📚 Recursos Adicionais

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [Go Test Comments](https://go.dev/blog/examples)
- [Testify Library](https://github.com/stretchr/testify) (opcional)

## 🎓 Exemplos Práticos

### Testar Erro

```go
func TestScriptExecute_FileNotFound(t *testing.T) {
    script := Script{Path: "/nonexistent/script.sh"}

    err := script.Execute()

    if err == nil {
        t.Error("Expected error for nonexistent script")
    }
}
```

### Testar com Mock

```go
func TestWithMockScript(t *testing.T) {
    scriptPath := testutil.CreateMockScript(t)

    script := Script{Path: scriptPath}
    err := script.Execute()

    testutil.AssertNil(t, err)
}
```

### Benchmark com Setup

```go
func BenchmarkComplexOperation(b *testing.B) {
    // Setup (não é medido)
    data := prepareTestData()

    // Reset timer após setup
    b.ResetTimer()

    // Operação medida
    for i := 0; i < b.N; i++ {
        ComplexOperation(data)
    }
}
```

---

**Última atualização**: 2026-03-14

Para mais informações, consulte:
- **[docs/INDEX.md](INDEX.md)** — índice da documentação
- **[README.md](../README.md)** — visão geral do projeto
- **[GETTING_STARTED.md](../GETTING_STARTED.md)** — como começar
- **[architecture/ARCHITECTURE.md](architecture/ARCHITECTURE.md)** — arquitetura
