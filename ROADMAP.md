# 🧭 Ghast Framework — Core Roadmap

> 📦 Estrutura modular e roadmap de desenvolvimento do núcleo do framework Ghast.

**Legenda:** 

| Status       | Ícone |
| :----------- | :---: |
| Concluído    |   ✅   |
| Em andamento |   🧩   |
| Planejado    |   🚧   |

## Nível 1 — Fundamentos de Injeção e Providers

  ✅ Base do container e providers implementada.

  - [x] **core/provider**
    - [x] Estrutura `Provider` com suporte a:
      - [x] `UseClass`
      - [x] `UseFactory`
      - [x] `UseValue`
    - [x] Função `NameFromType()` para identificação única
    - [x] Helpers (`NewClass`, `NewFactory`, `NewValue`)
  - [x] **core/container**
    - [x] Estrutura `Container` com:
      - [x] Registro de providers (`Register`)
      - [x] Resolução com cache (`Resolve`)
      - [x] Injeção via tag `inject:""`
      - [x] Detecção de dependências circulares
      - [x] Suporte a factories e valores fixos
    - [x] Helper de nível de pacote (`MustResolve`)

## Nível 2 — Modularização

  ✅ Sistema de módulos funcional com builder fluente.

  - [x] **core/module**
    - [x] Estrutura `Module` com:
      - [x] Providers, Controllers, Imports
      - [x] Registro recursivo no container (`Register`)
    - [x] **ModuleBuilder**
      - [x] `AddProviders()`
      - [x] `AddControllers()`
      - [x] `Import()`
      - [x] `Build()`
    - [x] Estrutura reorganizada (`builder.go`, `module.go`)

## Nível 3 — Middleware e Contexto HTTP

  ✅ Pipeline de middlewares genérico e agnóstico de framework.

  - [x] **core/middleware**
    - [x] Estrutura `Middleware` (`type.go`)
    - [x] Interface `Context` (`context.go`)
    - [x] Tipo `Handler` e construtor (`handler.go`)
    - [x] Função global `Chain()` (`package.go`)
    - [x] Modelo agnóstico de framework (Fiber, Gin, etc.)
    - [x] Execução ordenada e segura de middlewares

## Nível 4 — Execução de Requisições

  🧩 Em andamento — próxima etapa da implementação.

  - [x] **core/guard**
    - [x] Interface `Guard` (`CanActivate(ctx Context) bool`)
    - [x] Suporte a múltiplos guards por rota/controller
    - [x] Integração com o pipeline do controller
  - [x] **core/interceptor**
    - [x] Interface `Interceptor` (`Before/After`)
    - [x] Capacidade de transformar requisições e respostas
    - [x] Suporte a encadeamento e async (com goroutines)
  - [ ] **core/controller**
    - [ ] Estrutura `Controller`
    - [ ] **ControllerBuilder** com métodos declarativos:
      - [ ] `UseBasePath()`
      - [ ] `Get()`, `Post()`, `Put()`, `Patch()`, `Delete()`
    - [ ] Suporte a middlewares, guards e interceptors
    - [ ] Sistema de metadata para plugins (ex: Swagger)

## Nível 5 — Bootstrap da Aplicação

  🚧 Planejado — inicialização e ciclo de vida da aplicação.

  - [x] **core/application**
    - [x] Estrutura `Application`
    - [x] Registro e inicialização de módulos (`AppModule`)
    - [x] Integração com servidor HTTP (adapter de contexto)
    - [x] Hooks de ciclo de vida (`OnInit`, `OnDestroy`)
    - [x] Injeção de middlewares globais

## Nível 6 — Extensões e Plugins

  🚧 Planejado — sistema de extensões oficiais e CLI.

  - [ ] **Plugin System**
    - [ ] Interface para registro dinâmico de módulos/plugins
    - [ ] `github.com/leandroluk/ghast/swagger` (geração automática de docs)
    - [ ] `github.com/leandroluk/ghast/graphql` (integração declarativa com schema)
    - [ ] `github.com/leandroluk/ghast/logger` (logging configurável e avançado)
    - [ ] `github.com/leandroluk/ghast/metrics` (coleta de métricas em runtime)
  - [ ] **CLI**
    - [ ] Criação automática de módulos, controllers e providers
    - [ ] Comando `ghast new module <Name>`

## Nível 7 — Qualidade e Performance

  🚧 Planejado — otimização e testes avançados.

  - [ ] Benchmarks (`go test -bench`)
  - [ ] Análise de race conditions (`-race`)
  - [ ] Caching inteligente de dependências
  - [ ] Documentação gerada a partir dos metadados
  - [ ] Testes de integração e snapshots (com mocks HTTP)
