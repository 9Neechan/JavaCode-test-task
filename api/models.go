package api

// gin не может корректно распарсить uuid из uri, пожтому используем string
type getWalletRequest struct {
	WalletUuid string `uri:"id" binding:"required,uuid"`
}

type UpdateWalletBalanceRequest struct {
	WalletUuid    string `json:"wallet_uuid" binding:"required,uuid"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	OperationType string `json:"operation_type" binding:"required,oneof=DEPOSIT WITHDRAW"`
}

