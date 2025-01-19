package infrastructure

import (
	"arctfrex-customers/internal/account"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/conversation"
	"arctfrex-customers/internal/deposit"
	"arctfrex-customers/internal/inbox"
	"arctfrex-customers/internal/market"
	"arctfrex-customers/internal/news"
	"arctfrex-customers/internal/order"
	"arctfrex-customers/internal/otp"
	"arctfrex-customers/internal/report"
	backoffice "arctfrex-customers/internal/user/backoffice"
	mobile "arctfrex-customers/internal/user/mobile"
	"arctfrex-customers/internal/withdrawal"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB initializes a new database connection
func NewDB() *gorm.DB {
	// Open a connection to database
	// connStr := "host=172.17.0.4 user=postgres password=mysecret dbname=arctfrex port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	// connStr := "host=localhost user=postgres password=mysecret dbname=arctfrex port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	// Read database configuration from environment variables
	dbHost := os.Getenv(common.DATABASE_SERVER)
	dbPort := os.Getenv(common.DATABASE_PORT)
	dbUser := os.Getenv(common.DATABASE_USERNAME)
	dbPassword := os.Getenv(common.DATABASE_PASSWORD)
	dbName := os.Getenv(common.DATABASE_NAME)

	// Create the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Connected to database")
	databaseProfiling, _ := strconv.ParseBool(os.Getenv(common.DATABASE_PROFILING))

	if databaseProfiling {
		db = db.Debug()
	}

	// Automatically migrate the schema of the User struct
	db.AutoMigrate(
		&mobile.Users{},
		&mobile.UserProfile{},
		&mobile.UserAddress{},
		&mobile.UserEmployment{},
		&mobile.UserFinance{},
		&mobile.UserEmergencyContact{},
		&otp.Otp{},
		&backoffice.BackofficeUsers{},
		&market.Market{},
		&market.MarketCountry{},
		&market.MarketCurrencyRate{},
		&account.Account{},
		&deposit.Deposit{},
		&withdrawal.Withdrawal{},
		&news.News{},
		&news.NewsBulletin{},
		&order.Order{},
		&conversation.ConversationSession{},
		&conversation.ConversationMessage{},
		&conversation.ConversationStep{},
		&conversation.ConversationOption{},
		&report.Report{},
		&report.ReportOrders{},
		&report.ReportHistoryOrders{},
		&report.ReportDealData{},
		&inbox.Inbox{},
	)

	return db
}
