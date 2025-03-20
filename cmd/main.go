package main

import (
	"arctfrex-customers/internal/account"
	"arctfrex-customers/internal/auth"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/conversation"
	"arctfrex-customers/internal/deposit"
	"arctfrex-customers/internal/email"
	"arctfrex-customers/internal/inbox"
	infrastructure "arctfrex-customers/internal/infrastructure"
	market "arctfrex-customers/internal/market"
	"arctfrex-customers/internal/middleware"
	"arctfrex-customers/internal/news"
	"arctfrex-customers/internal/order"
	"arctfrex-customers/internal/otp"
	"arctfrex-customers/internal/report"
	"arctfrex-customers/internal/role"
	"arctfrex-customers/internal/storage"
	user_backoffice "arctfrex-customers/internal/user/backoffice"
	user_mobile "arctfrex-customers/internal/user/mobile"
	"arctfrex-customers/internal/whatsapp"
	"arctfrex-customers/internal/withdrawal"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	//Load environment variable
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	jwtSecretKey := os.Getenv(common.JWT_SECRET_KEY)
	applicationName := os.Getenv(common.APPLICATION_NAME)
	runMarketWorkerPriceUpdates, err := strconv.ParseBool(os.Getenv(common.RUN_MARKET_WORKER_PRICE_UPDATES))
	if err != nil {
		return
	}
	runMarketWorkerLiveMarketUpdates, err := strconv.ParseBool(os.Getenv(common.RUN_MARKET_WORKER_LIVE_MARKET_UPDATES))
	if err != nil {
		return
	}
	runNewsWorkerLatestNewsUpdates, err := strconv.ParseBool(os.Getenv(common.RUN_NEWS_WORKER_LATEST_NEWS_UPDATES))
	if err != nil {
		return
	}
	runNewsWorkerLatestNewsBulletinUpdates, err := strconv.ParseBool(os.Getenv(common.RUN_NEWS_WORKER_LATEST_NEWS_BULLETIN_UPDATES))
	if err != nil {
		return
	}
	runOrderWorkerCloseAllExpiredOrder, err := strconv.ParseBool(os.Getenv(common.RUN_ORDER_WORKER_CLOSE_ALL_EXPIRED_ORDER))
	if err != nil {
		return
	}

	//Database connection
	db := infrastructure.NewDB()

	//Engine
	engine := gin.Default()
	engine.Use(infrastructure.RequestResponseLogger())
	engine.Use(middleware.CORSMiddleware())

	// Use CORS middleware
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://panen.vercel.app", "https://admin.panenkapital.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("ngrok-skip-browser-warning", "true")
		c.Writer.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		c.Next()
	})

	//Auth
	jwtMiddleware := middleware.NewJWTMiddleware(jwtSecretKey)
	jwtService := auth.NewJWTService(jwtSecretKey, applicationName)
	authUsecase := auth.NewAuthUsecase(jwtService)
	auth.NewAuthHandler(engine, authUsecase)

	//Whatsapp
	twilioWhatsappSender := whatsapp.NewTwilioWhatsappSender()

	//User
	userApiClient := user_mobile.NewUserApiclient()
	userRepository := user_mobile.NewUserRepository(db)
	userUsecase := user_mobile.NewUserUseCase(userRepository, jwtService, userApiClient)
	user_mobile.NewUserHandler(engine, jwtMiddleware, userUsecase)

	// Initialize usecase
	marketApiClient := market.NewMarketApiClient()
	marketRepository := market.NewMarketRepository(db)
	marketUsecase := market.NewMarketUsecase(marketRepository, marketApiClient, userRepository)
	// Worker to fetch market price every minute
	marketPriceWorker := market.NewMarketWorker(marketUsecase)
	if runMarketWorkerPriceUpdates {
		marketPriceWorker.PriceUpdates(1 * time.Second)
	}
	if runMarketWorkerLiveMarketUpdates {
		marketPriceWorker.LiveMarketUpdates(3 * time.Second)
	}
	// Initialize handler
	market.NewMarketHandler(engine, jwtMiddleware, marketUsecase)

	newsApiClient := news.NewNewsApiClient()
	newsRepository := news.NewNewsRepository(db)
	newsUsecase := news.NewNewsUsecase(newsRepository, newsApiClient)
	newsWorker := news.NewNewsWorker(newsUsecase)
	if runNewsWorkerLatestNewsUpdates {
		newsWorker.NewsUpdates(10 * time.Minute)
	}
	if runNewsWorkerLatestNewsBulletinUpdates {
		newsWorker.NewsBulletinUpdates(10 * time.Minute)
	}
	news.NewNewsHandler(engine, jwtMiddleware, newsUsecase)

	//Email
	gomailSender := email.NewGomailSender()
	emailUseCase := email.NewEmailUseCase(gomailSender)
	email.NewEmailHandler(engine, emailUseCase)

	//Otp
	otpRepository := otp.NewOtpRepository(db)
	otpUsecase := otp.NewOtpUseCase(
		otpRepository,
		userRepository,
		twilioWhatsappSender,
		gomailSender,
	)
	otp.NewOtpHandler(engine, otpUsecase)

	// Initialize MinIO client
	minioClient, err := storage.NewMinioClient(
		os.Getenv(common.MINIO_ENDPOINT),
		os.Getenv(common.MINIO_ACCESS_KEY),
		os.Getenv(common.MINIO_SECRET_KEY),
		os.Getenv(common.MINIO_BUCKET_NAME),
	)
	if err != nil {
		panic(err)
	}

	//Account
	accountApiclient := account.NewAccountApiclient()
	accountRepository := account.NewAccountRepository(db)
	accountUsecase := account.NewAccountUsecase(accountRepository, accountApiclient)
	account.NewAccountHandler(engine, jwtMiddleware, accountUsecase)

	//Order
	orderRepository := order.NewOrderRepository(db)
	orderUsecase := order.NewOrderUsecase(orderRepository, accountRepository, marketRepository)
	order.NewOrderHandler(engine, jwtMiddleware, orderUsecase)
	orderWorker := order.NewOrderWorker(orderUsecase)
	if runOrderWorkerCloseAllExpiredOrder {
		orderWorker.CloseAllExpiredOrder(5 * time.Minute)
	}

	// Role
	roleRepository := role.NewRoleRepository(db)
	roleUseCase := role.NewRoleUseCase(roleRepository)
	role.NewRoleHandler(engine, jwtMiddleware, roleUseCase)

	//Backoffice user
	backofficeUserRepository := user_backoffice.NewBackofficeUserRepository(db)
	backofficeUserUsecase := user_backoffice.NewBackofficeUsecase(backofficeUserRepository, jwtService)
	user_backoffice.NewBackofficeHandler(engine, jwtMiddleware, backofficeUserUsecase)

	backofficeDepositApiclient := deposit.NewDepositApiclient()
	backofficeDepositRepository := deposit.NewDepositRepository(db)
	backofficeDepositUsecase := deposit.NewDepositUsecase(backofficeDepositRepository, accountRepository, backofficeDepositApiclient, marketRepository)
	deposit.NewDepositHandler(engine, jwtMiddleware, backofficeDepositUsecase)

	backofficeWithdrawalApiclient := withdrawal.NewWithdrawalApiclient()
	backofficeWithdrawalRepository := withdrawal.NewWithdrawalRepository(db)
	backofficeWithdrawalUsecase := withdrawal.NewWithdrawalUsecase(backofficeWithdrawalRepository, accountRepository, backofficeWithdrawalApiclient, marketRepository)
	withdrawal.NewWithdrawalHandler(engine, jwtMiddleware, backofficeWithdrawalUsecase)

	//Conversation Chat
	conversationRepository := conversation.NewConversationRepository(db)
	conversationRepository.SeedConversationSteps()
	conversationUsecase := conversation.NewConversationUsecase(conversationRepository)
	conversation.NewConversationHandler(engine, jwtMiddleware, conversationUsecase, jwtService)

	reportApiClient := report.NewReportApiClient()
	reportRepository := report.NewReportRepository(db)
	reportUsecase := report.NewReportUsecase(reportRepository, accountRepository, backofficeDepositRepository, reportApiClient)
	report.NewReportHandler(engine, jwtMiddleware, reportUsecase)
	reportWorker := report.NewReportWorker(reportUsecase)
	reportWorker.GroupUserLoginsUpdates(5 * time.Second)

	inboxRepository := inbox.NewInboxRepository(db)
	inboxUsecase := inbox.NewInboxUseCase(inboxRepository)
	inbox.NewInboxHandler(engine, jwtMiddleware, inboxUsecase)

	// Create usecase
	storageUsecase := storage.NewStorageUsecase(*minioClient, userRepository, backofficeDepositRepository, accountRepository)
	storage.NewStorageHandler(engine, jwtMiddleware, storageUsecase)

	//gin.SetMode(gin.ReleaseMode)

	//Port
	// log.Fatal(engine.RunTLS(":8543", "certs/cert.pem", "certs/key.pem"))
	log.Fatal(engine.RunTLS(os.Getenv(common.ENGINE_PORT), "certs/cert.pem", "certs/key.pem"))
}
