package app

import (
	"github.com/p1xray/pxr-sso/internal/controller/grpc"
	"github.com/p1xray/pxr-sso/internal/oidc/application/builder"
	"github.com/p1xray/pxr-sso/internal/oidc/application/cookie"
	"github.com/p1xray/pxr-sso/internal/oidc/application/generator"
	"github.com/p1xray/pxr-sso/internal/oidc/application/usecase/authorize"
	"github.com/p1xray/pxr-sso/internal/oidc/application/usecase/consent"
	"github.com/p1xray/pxr-sso/internal/oidc/application/usecase/grant"
	"github.com/p1xray/pxr-sso/internal/oidc/application/usecase/login"
	"github.com/p1xray/pxr-sso/internal/oidc/application/usecase/register"
	"github.com/p1xray/pxr-sso/internal/oidc/application/usecase/token"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/cache"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/cache/redis"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/repository"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage"
	"github.com/p1xray/pxr-sso/internal/oidc/infrastructure/storage/postgresql"
	"github.com/p1xray/pxr-sso/pkg/grpcserver"
	"github.com/p1xray/pxr-sso/pkg/logger"
	"log/slog"
)

// diContainer is a lazy-initialized dependency container.
type diContainer struct {
	// Infrastructure:
	// 	config, logger, metrics, etc
	cfg *Config
	log *slog.Logger

	// 	storages
	db    storage.Storage
	cache cache.Cache

	// 	repositories
	repo repository.Repository

	// Application:
	// 	builders
	uriBuilder builder.URIBuilder

	//	cookie
	sessionCookieEncoding cookie.SessionEncoding

	//	generators
	tokensGenerator generator.TokensGenerator

	// 	use cases:
	//		grant
	authorizedGrantUseCase grant.AuthorizationGrantFlow

	//		authorize
	authorizeLoginFlowProcessor         authorize.LoginFlowProcessor
	authorizeNoneFlowProcessor          authorize.NoneFlowProcessor
	authorizeConsentFlowProcessor       authorize.ConsentFlowProcessor
	authorizeSelectAccountFlowProcessor authorize.SelectAccountFlowProcessor
	authorizeFlowSwitcher               authorize.FlowSwitcher
	authorizeUseCase                    authorize.Authorize

	//		login
	loginUseCase login.Login

	//		register
	registerUseCase register.Register

	//		consent
	consentUseCase consent.Consent

	//		consent card
	consentCardReaderUseCase consent.CardReader

	//		token
	tokenUseCase token.RequestProcessor

	// API:
	// 	gRPC
	grpcServer grpcserver.Server
}

// newDIContainer creates a new empty DI container.
// All fields are nil, dependencies will be created lazily on the first access.
func newDIContainer() *diContainer {
	return &diContainer{}
}

// Config returns the application configuration.
func (d *diContainer) Config() *Config {
	if d.cfg == nil {
		cfgLoader := newConfigLoader()
		cfg := cfgLoader.MustLoad()

		d.cfg = cfg
	}

	return d.cfg
}

// Logger returns the logger configured depending on the environment specified in the config.
func (d *diContainer) Logger() *slog.Logger {
	if d.log == nil {
		d.log = logger.SetupLogger(d.Config().Env)
	}

	return d.log
}

// DB returns the connection to the database.
func (d *diContainer) DB() storage.Storage {
	if d.db == nil {
		postgresDB, err := postgresql.New(d.Config().Postgres)
		if err != nil {
			panic(err)
		}

		d.db = postgresDB
	}

	return d.db
}

// Cache returns the connection to the cache.
func (d *diContainer) Cache() cache.Cache {
	if d.cache == nil {
		redisCache, err := redis.New(d.Config().Redis)
		if err != nil {
			panic(err)
		}

		d.cache = redisCache
	}

	return d.cache
}

// Repository returns the repository.
func (d *diContainer) Repository() repository.Repository {
	if d.repo == nil {
		d.repo = repository.New(d.DB())
	}

	return d.repo
}

// URIBuilder returns the URI builder to redirect the user agent to a specified URI.
func (d *diContainer) URIBuilder() builder.URIBuilder {
	if d.uriBuilder == nil {
		d.uriBuilder = builder.NewURIBuilder(d.Config().URIBuilder)
	}

	return d.uriBuilder
}

// SessionCookieEncoding returns the encoding for session cookie.
func (d *diContainer) SessionCookieEncoding() cookie.SessionEncoding {
	if d.sessionCookieEncoding == nil {
		d.sessionCookieEncoding = cookie.NewSessionEncoding(d.Config().SessionCookie)
	}

	return d.sessionCookieEncoding
}

// TokensGenerator returns the JWT generator for OIDC.
func (d *diContainer) TokensGenerator() generator.TokensGenerator {
	if d.tokensGenerator == nil {
		d.tokensGenerator = generator.NewTokensGenerator(d.Config().Token)
	}

	return d.tokensGenerator
}

// AuthorizedGrantUseCase returns the use case for processing authorized grant.
func (d *diContainer) AuthorizedGrantUseCase() grant.AuthorizationGrantFlow {
	if d.authorizedGrantUseCase == nil {
		d.authorizedGrantUseCase = grant.NewUseCase(d.URIBuilder(), d.Repository(), d.Repository(), d.Cache())
	}

	return d.authorizedGrantUseCase
}

// AuthorizeLoginFlowProcessor returns the processor for authorize flow with logging in interaction.
func (d *diContainer) AuthorizeLoginFlowProcessor() authorize.LoginFlowProcessor {
	if d.authorizeLoginFlowProcessor == nil {
		d.authorizeLoginFlowProcessor = authorize.NewLoginFlowProcessor(d.URIBuilder(), d.Cache())
	}

	return d.authorizeLoginFlowProcessor
}

// AuthorizeNoneFlowProcessor returns the processor for authorize without interaction flow.
func (d *diContainer) AuthorizeNoneFlowProcessor() authorize.NoneFlowProcessor {
	if d.authorizeNoneFlowProcessor == nil {
		d.authorizeNoneFlowProcessor = authorize.NewNoneFlowProcessor(d.AuthorizedGrantUseCase())
	}

	return d.authorizeNoneFlowProcessor
}

// AuthorizeConsentFlowProcessor returns the processor for authorize flow with confirming consent interaction.
func (d *diContainer) AuthorizeConsentFlowProcessor() authorize.ConsentFlowProcessor {
	if d.authorizeConsentFlowProcessor == nil {
		d.authorizeConsentFlowProcessor = authorize.NewConsentFlowProcessor(d.URIBuilder(), d.Cache())
	}

	return d.authorizeConsentFlowProcessor
}

// AuthorizeSelectAccountFlowProcessor returns the processor for authorize flow with selecting account interaction.
func (d *diContainer) AuthorizeSelectAccountFlowProcessor() authorize.SelectAccountFlowProcessor {
	if d.authorizeSelectAccountFlowProcessor == nil {
		d.authorizeSelectAccountFlowProcessor = authorize.NewSelectAccountFlowProcessor(d.URIBuilder(), d.Cache())
	}

	return d.authorizeSelectAccountFlowProcessor
}

// AuthorizeFlowSwitcher returns the processor for defining the interaction flow.
func (d *diContainer) AuthorizeFlowSwitcher() authorize.FlowSwitcher {
	if d.authorizeFlowSwitcher == nil {
		d.authorizeFlowSwitcher = authorize.NewFlowSwitcher(
			d.AuthorizeLoginFlowProcessor(),
			d.AuthorizeNoneFlowProcessor(),
			d.AuthorizeConsentFlowProcessor(),
			d.AuthorizeSelectAccountFlowProcessor(),
		)
	}

	return d.authorizeFlowSwitcher
}

// AuthorizeUseCase returns the use case for processing authorize request.
func (d *diContainer) AuthorizeUseCase() authorize.Authorize {
	if d.authorizeUseCase == nil {
		d.authorizeUseCase = authorize.NewUseCase(
			d.Logger(),
			d.URIBuilder(),
			d.Repository(),
			d.Repository(),
			d.SessionCookieEncoding(),
			d.AuthorizeFlowSwitcher(),
		)
	}

	return d.authorizeUseCase
}

// LoginUseCase returns the use case for processing logging in user.
func (d *diContainer) LoginUseCase() login.Login {
	if d.loginUseCase == nil {
		d.loginUseCase = login.New(
			d.Logger(),
			d.URIBuilder(),
			d.Repository(),
			d.Repository(),
			d.Repository(),
			d.Cache(),
			d.SessionCookieEncoding(),
		)
	}

	return d.loginUseCase
}

// RegisterUseCase returns the use case for processing registering user.
func (d *diContainer) RegisterUseCase() register.Register {
	if d.registerUseCase == nil {
		d.registerUseCase = register.New(
			d.Logger(),
			d.URIBuilder(),
			d.Repository(),
			d.Repository(),
			d.Repository(),
			d.Repository(),
			d.Cache(),
			d.SessionCookieEncoding(),
		)
	}

	return d.registerUseCase
}

// ConsentUseCase returns the use case for processing confirming consent.
func (d *diContainer) ConsentUseCase() consent.Consent {
	if d.consentUseCase == nil {
		d.consentUseCase = consent.New(
			d.Logger(),
			d.Repository(),
			d.Repository(),
			d.Cache(),
			d.SessionCookieEncoding(),
			d.AuthorizedGrantUseCase(),
		)
	}

	return d.consentUseCase
}

func (d *diContainer) ConsentCardReaderUseCase() consent.CardReader {
	if d.consentCardReaderUseCase == nil {
		d.consentCardReaderUseCase = consent.NewCardReader(
			d.Logger(),
			d.Cache(),
			d.Repository(),
		)
	}

	return d.consentCardReaderUseCase
}

// TokenUseCase returns the use case for processing token request.
func (d *diContainer) TokenUseCase() token.RequestProcessor {
	if d.tokenUseCase == nil {
		d.tokenUseCase = token.New(
			d.Logger(),
			d.Cache(),
			d.Cache(),
			d.Repository(),
			d.Repository(),
			d.TokensGenerator(),
		)
	}

	return d.tokenUseCase
}

// GRPCServer returns the gRPC server with configured handlers.
func (d *diContainer) GRPCServer() grpcserver.Server {
	if d.grpcServer == nil {
		grpcServer := grpcserver.New(grpcserver.WithPort(d.Config().Port))

		grpc.NewRouter(
			grpcServer.Registrar(),
			d.AuthorizeUseCase(),
			d.LoginUseCase(),
			d.RegisterUseCase(),
			d.ConsentUseCase(),
			d.ConsentCardReaderUseCase(),
			d.TokenUseCase(),
		)

		d.grpcServer = grpcServer
	}

	return d.grpcServer
}
