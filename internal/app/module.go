package app

import (
	"ditto/internal/ditto"
	"ditto/internal/http/handler"
	"ditto/internal/http/router"
	"ditto/internal/repository"
	"ditto/internal/service"

	"go.uber.org/fx"
)

var Module = fx.Options(
	handler.Module,
	router.Module,
	repository.Module,
	service.Module,
	ditto.Module,
)
