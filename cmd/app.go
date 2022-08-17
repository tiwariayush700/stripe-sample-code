package main

import (
	`context`
	`log`

	`github.com/gin-gonic/gin`

	`stripe.com/docs/payments/core/config`
	`stripe.com/docs/payments/handler`
	repositoryimpl `stripe.com/docs/payments/repository/impl`
	serviceimpl `stripe.com/docs/payments/service/impl`
)

type app struct {
	Configuration *config.Config
	Router        *gin.Engine
}

func (a *app) Start(ctx context.Context) {

	//NewPaymentRepositoryImpl
	paymentRepository, err := repositoryimpl.NewPaymentRepositoryImpl(ctx, mongoUri(a.Configuration), "cb-payments")
	if err != nil {
		log.Fatalf("NewPaymentRepositoryImpl err %v", err)
	}

	//NewPaymentServiceImpl
	paymentService := serviceimpl.NewPaymentServiceImpl(a.Configuration.PaymentConfig, paymentRepository)

	//NewPaymentHandler
	paymentHandler := handler.NewPayment(paymentService)

	paymentHandler.RegisterRoutes(ctx, a.Router)

	log.Printf("Application loaded successfully on port : %s", a.Configuration.Port)
	log.Fatal(a.Router.Run(":" + a.Configuration.Port))
}

func NewApp(configuration *config.Config, router *gin.Engine) *app {
	return &app{Configuration: configuration, Router: router}
}

func mongoUri(appConfig *config.Config) string {

	auth := ""
	m := "mongodb://" + auth + appConfig.MongoServer

	return m
}
