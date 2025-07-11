package common

const (
	ENGINE_PORT                                  = "ENGINE_PORT"
	DATABASE_PROFILING                           = "DATABASE_PROFILING"
	DATABASE_SERVER                              = "DATABASE_SERVER"
	DATABASE_USERNAME                            = "DATABASE_USERNAME"
	DATABASE_PASSWORD                            = "DATABASE_PASSWORD"
	DATABASE_NAME                                = "DATABASE_NAME"
	DATABASE_PORT                                = "DATABASE_PORT"
	STRING_EMPTY                                 = ""
	STRING_SPACE                                 = " "
	STRING_DOUBLE_SPACE                          = "  "
	DASH                                         = "-"
	JWT_SECRET_KEY                               = "JWT_SECRET_KEY"
	APPLICATION_NAME                             = "APPLICATION_NAME"
	HTTP_METHOD_GET                              = "GET"
	HTTP_METHOD_POST                             = "POST"
	OTP_GENERATOR_SECRET                         = "OTP_GENERATOR_SECRET"
	EMAIL_FROM                                   = "EMAIL_FROM"
	OTP_EMAIL_SUBJECT                            = "OTP_EMAIL_SUBJECT"
	OTP_SEND_WITH_EMAIL                          = "OTP_SEND_WITH_EMAIL"
	RUN_MARKET_WORKER_PRICE_UPDATES              = "RUN_MARKET_WORKER_PRICE_UPDATES"
	RUN_MARKET_WORKER_LIVE_MARKET_UPDATES        = "RUN_MARKET_WORKER_LIVE_MARKET_UPDATES"
	RUN_NEWS_WORKER_LATEST_NEWS_UPDATES          = "RUN_NEWS_WORKER_LATEST_NEWS_UPDATES"
	RUN_NEWS_WORKER_LATEST_NEWS_BULLETIN_UPDATES = "RUN_NEWS_WORKER_LATEST_NEWS_BULLETIN_UPDATES"
	RUN_ORDER_WORKER_CLOSE_ALL_EXPIRED_ORDER     = "RUN_ORDER_WORKER_CLOSE_ALL_EXPIRED_ORDER"
	LOG_CONSOLE_RESET                            = "\033[0m"
	LOG_CONSOLE_RED                              = "\033[31m"
	LOG_CONSOLE_GREEN                            = "\033[32m"
	LOG_CONSOLE_YELLOW                           = "\033[33m"
	LOG_CONSOLE_BLUE                             = "\033[34m"
	LOG_CONSOLE_CYAN                             = "\033[36m"
	MINIO_ENDPOINT                               = "MINIO_ENDPOINT"
	MINIO_BASEURL                                = "MINIO_BASEURL"
	MINIO_ENDPOINT_SECURED                       = "MINIO_ENDPOINT_SECURED"
	MINIO_ACCESS_KEY                             = "MINIO_ACCESS_KEY"
	MINIO_SECRET_KEY                             = "MINIO_SECRET_KEY"
	MINIO_BUCKET_NAME                            = "BUCKET_NAME"
	lowerCharset                                 = "abcdefghijklmnopqrstuvwxyz"
	upperCharset                                 = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberCharset                                = "0123456789"
	specialCharset                               = "!@#$%^&*()-_+=<>?~"
	allCharset                                   = lowerCharset + upperCharset + numberCharset + specialCharset
	passwordLength                               = 12 // Minimum recommended length
	ARC_META_INTEGRATOR_BASEURL                  = "ARC_META_INTEGRATOR_BASEURL"
)

const (
	WorkflowDepositApprover    = "deposit-approver"
	WorkflowWithdrawalApprover = "withdrawal-approver"
)

const (
	Finance        = "finance"
	Settlement     = "settlement"
	CreditIn       = "credit-in"
	CreditOut      = "credit-out"
	DepositType    = "deposit"
	WithdrawalType = "withdrawal"
	SPA            = "spa"
	Multi          = "multi"
)
