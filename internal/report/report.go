package report

import "arctfrex-customers/internal/base"

type Report struct {
	Code        string `gorm:"primary_key" json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`

	base.BaseModel
}

type ReportOrders struct {
	OrderID          int64   `json:"Order"`            // Maps to "Order" field
	ExternalID       string  `json:"ExternalID"`       // Maps to "ExternalID" field
	Login            int64   `json:"Login"`            // Maps to "Login" field
	Dealer           int     `json:"Dealer"`           // Maps to "Dealer" field
	Symbol           string  `json:"Symbol"`           // Maps to "Symbol" field
	Digits           int     `json:"Digits"`           // Maps to "Digits" field
	DigitsCurrency   int     `json:"DigitsCurrency"`   // Maps to "DigitsCurrency" field
	ContractSize     float64 `json:"ContractSize"`     // Maps to "ContractSize" field
	State            int     `json:"State"`            // Maps to "State" field
	Reason           int     `json:"Reason"`           // Maps to "Reason" field
	TimeSetup        int64   `json:"TimeSetup"`        // Maps to "TimeSetup" field
	TimeExpiration   int64   `json:"TimeExpiration"`   // Maps to "TimeExpiration" field
	TimeDone         int64   `json:"TimeDone"`         // Maps to "TimeDone" field
	TimeSetupMsc     int64   `json:"TimeSetupMsc"`     // Maps to "TimeSetupMsc" field
	TimeDoneMsc      int64   `json:"TimeDoneMsc"`      // Maps to "TimeDoneMsc" field
	ModifyFlags      int     `json:"ModifyFlags"`      // Maps to "ModifyFlags" field
	Type             int     `json:"Type"`             // Maps to "Type" field
	TypeFill         int     `json:"TypeFill"`         // Maps to "TypeFill" field
	TypeTime         int     `json:"TypeTime"`         // Maps to "TypeTime" field
	PriceOrder       float64 `json:"PriceOrder"`       // Maps to "PriceOrder" field
	PriceTrigger     float64 `json:"PriceTrigger"`     // Maps to "PriceTrigger" field
	PriceCurrent     float64 `json:"PriceCurrent"`     // Maps to "PriceCurrent" field
	PriceSL          float64 `json:"PriceSL"`          // Maps to "PriceSL" field
	PriceTP          float64 `json:"PriceTP"`          // Maps to "PriceTP" field
	VolumeInitial    int64   `json:"VolumeInitial"`    // Maps to "VolumeInitial" field
	VolumeInitialExt int64   `json:"VolumeInitialExt"` // Maps to "VolumeInitialExt" field
	VolumeCurrent    int64   `json:"VolumeCurrent"`    // Maps to "VolumeCurrent" field
	VolumeCurrentExt int64   `json:"VolumeCurrentExt"` // Maps to "VolumeCurrentExt" field
	ExpertID         int64   `json:"ExpertID"`         // Maps to "ExpertID" field
	ExpertPositionID int64   `json:"ExpertPositionID"` // Maps to "ExpertPositionID" field
	PositionByID     int64   `json:"PositionByID"`     // Maps to "PositionByID" field
	Comment          string  `json:"Comment"`          // Maps to "Comment" field
	ActivationMode   int     `json:"ActivationMode"`   // Maps to "ActivationMode" field
	ActivationTime   int64   `json:"ActivationTime"`   // Maps to "ActivationTime" field
	ActivationPrice  float64 `json:"ActivationPrice"`  // Maps to "ActivationPrice" field
	ActivationFlags  int     `json:"ActivationFlags"`  // Maps to "ActivationFlags" field

	base.BaseModel
}

type ReportHistoryOrders struct {
	OrderID          int64   `json:"Order" gorm:"column:order_id;primaryKey"`
	ExternalID       string  `json:"ExternalID" gorm:"column:external_id"`
	Login            int64   `json:"Login" gorm:"column:login"`
	Dealer           int     `json:"Dealer" gorm:"column:dealer"`
	Symbol           string  `json:"Symbol" gorm:"column:symbol"`
	Digits           int     `json:"Digits" gorm:"column:digits"`
	DigitsCurrency   int     `json:"DigitsCurrency" gorm:"column:digits_currency"`
	ContractSize     float64 `json:"ContractSize" gorm:"column:contract_size"`
	State            int     `json:"State" gorm:"column:state"`
	Reason           int     `json:"Reason" gorm:"column:reason"`
	TimeSetup        int64   `json:"TimeSetup" gorm:"column:time_setup"`
	TimeExpiration   int64   `json:"TimeExpiration" gorm:"column:time_expiration"`
	TimeDone         int64   `json:"TimeDone" gorm:"column:time_done"`
	TimeSetupMsc     int64   `json:"TimeSetupMsc" gorm:"column:time_setup_msc"`
	TimeDoneMsc      int64   `json:"TimeDoneMsc" gorm:"column:time_done_msc"`
	ModifyFlags      int     `json:"ModifyFlags" gorm:"column:modify_flags"`
	Type             int     `json:"Type" gorm:"column:type"`
	TypeFill         int     `json:"TypeFill" gorm:"column:type_fill"`
	TypeTime         int     `json:"TypeTime" gorm:"column:type_time"`
	PriceOrder       float64 `json:"PriceOrder" gorm:"column:price_order"`
	PriceTrigger     float64 `json:"PriceTrigger" gorm:"column:price_trigger"`
	PriceCurrent     float64 `json:"PriceCurrent" gorm:"column:price_current"`
	PriceSL          float64 `json:"PriceSL" gorm:"column:price_sl"`
	PriceTP          float64 `json:"PriceTP" gorm:"column:price_tp"`
	VolumeInitial    int64   `json:"VolumeInitial" gorm:"column:volume_initial"`
	VolumeInitialExt int64   `json:"VolumeInitialExt" gorm:"column:volume_initial_ext"`
	VolumeCurrent    int64   `json:"VolumeCurrent" gorm:"column:volume_current"`
	VolumeCurrentExt int64   `json:"VolumeCurrentExt" gorm:"column:volume_current_ext"`
	ExpertID         int64   `json:"ExpertID" gorm:"column:expert_id"`
	ExpertPositionID int64   `json:"ExpertPositionID" gorm:"column:expert_position_id"`
	PositionByID     int64   `json:"PositionByID" gorm:"column:position_by_id"`
	Comment          string  `json:"Comment" gorm:"column:comment"`
	ActivationMode   int     `json:"ActivationMode" gorm:"column:activation_mode"`
	ActivationTime   int64   `json:"ActivationTime" gorm:"column:activation_time"`
	ActivationPrice  float64 `json:"ActivationPrice" gorm:"column:activation_price"`
	ActivationFlags  int     `json:"ActivationFlags" gorm:"column:activation_flags"`
}

type ReportDealData struct {
	Deal            int     `json:"deal"`
	ExternalID      string  `json:"external_id"`
	Login           int     `json:"login"`
	Dealer          int     `json:"dealer"`
	Order           int     `json:"order"`
	Action          int     `json:"action"`
	Entry           int     `json:"entry"`
	Reason          int     `json:"reason"`
	Digits          int     `json:"digits"`
	DigitsCurrency  int     `json:"digits_currency"`
	ContractSize    int     `json:"contract_size"`
	Time            int64   `json:"time"`
	TimeMsc         int64   `json:"time_msc"`
	Symbol          string  `json:"symbol"`
	Price           float64 `json:"price"`
	Volume          int     `json:"volume"`
	VolumeExt       int     `json:"volume_ext"`
	Profit          int     `json:"profit"`
	Storage         int     `json:"storage"`
	Commission      float64 `json:"commission"`
	CommissionAgent int     `json:"commission_agent"`
	RateProfit      float64 `json:"rate_profit"`
	RateMargin      int     `json:"rate_margin"`
	ExpertID        int     `json:"expert_id"`
	PositionID      int     `json:"position_id"`
	Comment         string  `json:"comment"`
	ProfitRaw       int     `json:"profit_raw"`
	PricePosition   int     `json:"price_position"`
	VolumeClosed    int     `json:"volume_closed"`
	VolumeClosedExt int     `json:"volume_closed_ext"`
	TickValue       int     `json:"tick_value"`
	TickSize        int     `json:"tick_size"`
	Flags           int     `json:"flags"`
	Gateway         string  `json:"gateway"`
	PriceGateway    float64 `json:"price_gateway"`
	ModifyFlags     int     `json:"modify_flags"`
	PriceSL         float64 `json:"price_sl"`
	PriceTP         float64 `json:"price_tp"`
}

type ReportData struct {
	Code   string   `json:"code"`
	Column []string `json:"column"`
	Data   any      `json:"data"`
}

type Presentasi struct {
	Value float64 `json:"value"`
	Color string  `json:"color"`
}

type AccountGetManifestResponse struct {
	Login          int     `json:"Login"`
	Balance        float64 `json:"Balance"`
	Equity         float64 `json:"Equity"`
	Profit         float64 `json:"Profit"`
	Floating       float64 `json:"Floating"`
	MarginInitial  float64 `json:"MarginInitial"`
	MarginLeverage int     `json:"MarginLeverage"`
	Margin         float64 `json:"Margin"`
	Commission     float64 `json:"Commission"`
	Swap           float64 `json:"Swap"`
	MarginLevel    float64 `json:"MarginLevel"`
	SOEquity       float64 `json:"SOEquity"`
	SOMargin       float64 `json:"SOMargin"`
	SOLevel        float64 `json:"SOLevel"`
	// PresentasiRunningProfit float64 `json:"presentasi_running_profit"`
	// PresentasiAllProfit     float64 `json:"presentasi_all_profit"`
	PresentasiRunningProfit Presentasi `json:"PresentasiRunningProfit"`
	PresentasiAllProfit     Presentasi `json:"PresentasiAllProfit"`
}

type ReportProfitLossData struct {
	MetaLoginID           int64   `json:"MetaLoginID"`
	Name                  string  `json:"Name"`
	DomCity               string  `json:"DomCity"`
	Currency              string  `json:"Currency"`
	CurrencyRate          float64 `json:"CurrencyRate"`
	TotalDepositAmount    float64 `json:"TotalDepositAmount"`
	TotalWithdrawalAmount float64 `json:"TotalWithdrawalAmount"`
	PrevEquity            float64 `json:"PrevEquity"`
	Nmii                  float64 `json:"Nmii"`
	LastEquity            float64 `json:"LastEquity"`
	GrossProfit           float64 `json:"GrossProfit"`
	GrossProfitUSD        float64 `json:"GrossProfitUSD"`
	SingleSideLot         float64 `json:"SingleSideLot"`
	Commission            float64 `json:"Commission"`
	Rebate                float64 `json:"Rebate"`
	PrevBadDebt           float64 `json:"PrevBadDebt"`
	LastBadDebt           float64 `json:"LastBadDebt"`
	NetProfit             float64 `json:"NetProfit"`
	NetProfitUSD          float64 `json:"NetProfitUSD"`
	AccountID             int64   `json:"AccountID"`
	UserID                int64   `json:"UserID"`
}

type ReportApiResponse struct {
	base.ApiResponse
}

type ReportRepository interface {
	GetActiveReports() (*[]Report, error)
}
