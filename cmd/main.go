package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"arctfrex-customers/internal/api"
	"arctfrex-customers/internal/auth"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/conversation"
	"arctfrex-customers/internal/email"
	"arctfrex-customers/internal/handler"
	"arctfrex-customers/internal/inbox"
	infrastructure "arctfrex-customers/internal/infrastructure"
	market "arctfrex-customers/internal/market"
	"arctfrex-customers/internal/middleware"
	"arctfrex-customers/internal/news"
	"arctfrex-customers/internal/order"
	"arctfrex-customers/internal/otp"
	"arctfrex-customers/internal/report"
	"arctfrex-customers/internal/repository"
	"arctfrex-customers/internal/role"
	"arctfrex-customers/internal/storage"
	"arctfrex-customers/internal/usecase"
	"arctfrex-customers/internal/whatsapp"
)

func main() {

	//Load environment variable
	// If not found, try loading from the parent directory (e.g., for 'cmd' folder)
	// Get the parent directory
	// Adjust paths for terminal usage
	// If not in 'cmd', we are probably in a debugger or a different environment
	// Resolve the paths dynamically based on the current working directory
	certPath, keyPath, err := initEnv()
	if err != nil {
		return
	}
	common.InitSonyflake()

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
		AllowOrigins:     []string{"http://localhost:3000", "https://panen.vercel.app", "https://admin.panenkapital.com", "https://dev-admin.panenkapital.com"},
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
	userApiClient := api.NewUserApiclient()
	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUseCase(userRepository, jwtService, userApiClient)
	handler.NewUserHandler(engine, jwtMiddleware, userUsecase)

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

	//api
	marketApiClient := api.NewMarketApiClient()
	accountApiclient := api.NewAccountApiclient()
	backofficeDepositApiclient := api.NewDepositApiclient()
	backofficeWithdrawalApiclient := api.NewWithdrawalApiclient()

	//repository
	marketRepository := repository.NewMarketRepository(db)
	accountRepository := repository.NewAccountRepository(db)
	workflowSettingRepository := repository.NewWorkflowSettingRepository(db)
	workflowApproverRepository := repository.NewWorkflowApproverRepository(db)
	backofficeDepositRepository := repository.NewDepositRepository(db)
	backofficeWithdrawalRepository := repository.NewWithdrawalRepository(db)
	backofficeUserRepository := repository.NewBackofficeUserRepository(db)

	//usecase
	marketUsecase := usecase.NewMarketUsecase(marketRepository, marketApiClient, userRepository)
	accountUsecase := usecase.NewAccountUsecase(accountRepository, accountApiclient)
	backofficeUserUsecase := usecase.NewBackofficeUsecase(backofficeUserRepository, jwtService)
	backofficeWithdrawalUsecase := usecase.NewWithdrawalUsecase(backofficeWithdrawalRepository, accountRepository, backofficeWithdrawalApiclient, marketRepository, workflowSettingRepository, workflowApproverRepository)
	backofficeDepositUsecase := usecase.NewDepositUsecase(
		backofficeDepositRepository,
		accountRepository,
		backofficeDepositApiclient,
		marketRepository,
		workflowSettingRepository,
		workflowApproverRepository,
	)
	workflowApproverUsecase := usecase.NewWorkflowApproverUsecase(
		workflowApproverRepository,
		backofficeDepositRepository,
		backofficeWithdrawalRepository,
		accountRepository,
		marketRepository,
		backofficeDepositApiclient,
		backofficeWithdrawalApiclient,
		db,
	)

	//handler
	handler.NewMarketHandler(engine, jwtMiddleware, marketUsecase)
	handler.NewAccountHandler(engine, jwtMiddleware, accountUsecase)
	handler.NewDepositHandler(engine, jwtMiddleware, backofficeDepositUsecase)
	handler.NewWithdrawalHandler(engine, jwtMiddleware, backofficeWithdrawalUsecase)
	handler.NewWorkflowApproverHandler(engine, jwtMiddleware, workflowApproverUsecase)
	handler.NewBackofficeHandler(engine, jwtMiddleware, backofficeUserUsecase)

	// Worker to fetch market price every minute
	marketPriceWorker := market.NewMarketWorker(marketUsecase)
	if runMarketWorkerPriceUpdates {
		marketPriceWorker.PriceUpdates(1 * time.Second)
	}
	if runMarketWorkerLiveMarketUpdates {
		marketPriceWorker.LiveMarketUpdates(5 * time.Second)
	}

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
	log.Fatal(engine.RunTLS(os.Getenv(common.ENGINE_PORT), certPath, keyPath))
}

func initEnv() (string, string, error) {
	certPath := "certs/cert.pem"
	keyPath := "certs/key.pem"
	err := godotenv.Load(".env")
	if err != nil {

		cwd, _ := os.Getwd()
		parentDir := filepath.Dir(cwd)
		err = godotenv.Load(filepath.Join(parentDir, ".env"))

		if filepath.Base(cwd) == "cmd" {

			certPath = filepath.Join(cwd, "../certs/cert.pem")
			keyPath = filepath.Join(cwd, "../certs/key.pem")
		} else {

			certPath, err = filepath.Abs(certPath)
			if err != nil {
				log.Fatal("Error getting absolute cert path:", err)
			}
			keyPath, err = filepath.Abs(keyPath)
			if err != nil {
				log.Fatal("Error getting absolute key path:", err)
			}
		}
	}

	if err != nil {
		log.Fatal("Error loading .env file")
		return "", "", err
	}

	return certPath, keyPath, err
}
