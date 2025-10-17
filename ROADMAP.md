# P0 – núcleo que falta

  - Request scope no DI — ✅ feito (v1)
    - Container.BeginRequest/EndRequest, cache por reqID, ResolveWithReq, MustResolveRequest[T](ctn, ctx).
    - Wrapper wrapWithRequestScope abre/fecha o escopo por requisição.
  - Ciclo de vida completo — 🟨 parcial
    - ✅ OnModuleInit (chamado ao registrar singletons).
    - ✅ OnApplicationShutdown(ctx) e ✅ OnModuleDestroy() (chamados no shutdown).
    - ✅ Graceful shutdown (adapters em goroutine + Shutdown(ctx) no Fiber/Gin).
    - ⛔ Falta: OnApplicationBootstrap (hook pós-boot).
  - Pipeline real — ✅ feito
    - Ordem: Interceptors → Guards → Middlewares → Handler → Exception Filters → (wrapper externo) Request Scope.
    - Filtro default global configurável (DefaultFilter / DisableDefaultFilter).
  - Fronteiras de módulo & exports — ❌ faltando
    - Ainda é container global; não há exports/resolução por fronteira.
  - Tokens & variantes de provider — ❌ faltando
    - Sem useExisting, sem tokens arbitrários, sem @Optional, sem forwardRef.
  - Pipes (transform/validation) — ❌ faltando
    - Sem Pipe/ValidationPipe/transforms antes do handler.
  - Param binding “decorator-like” — ❌ faltando
    - Sem @Param/@Query/@Body/@Headers equivalentes (helpers) e pipes por parâmetro.

# P1 – DX e roteamento

  - Escopos intermediários 🟡
    Controller-scoped e Module-scoped (guards/interceptors/filters/middlewares).
    Hoje: global e por rota; falta nível controller/módulo.

  - Prefixo global & versionamento ❌
    setGlobalPrefix('/v1'), version por URI/header.

  - Exception types prontos 🟡
    Mapas prontos tipo BadRequest(), NotFound(), etc.
    Hoje: HttpError genérico + DefaultFilter.

  - Config de adapter via app ❌
    CORS, body-limit, timeouts etc expostos pelo Builder, não só no Fiber/Gin.

  - Graceful shutdown real ❌
    Capturar sinais, server.Shutdown(ctx), chamar hooks de módulos/providers.

# P2 – ecossistema

  - Microservices (TCP, NATS, RMQ, Kafka, gRPC) ❌

  - WebSockets (Gateways) ❌

  - GraphQL ❌

  - Swagger/OpenAPI ❌

  - Cache/Config/Schedule/CQRS/EventEmitter/i18n ❌

  - CLI & Testing utilities ❌ (geradores, TestingModule, overrides, e2e harness)