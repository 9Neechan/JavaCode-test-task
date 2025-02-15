package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	mockdb "github.com/9Neechan/JavaCode-test-task/db/mock"
	db "github.com/9Neechan/JavaCode-test-task/db/sqlc"
	"github.com/9Neechan/JavaCode-test-task/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// randomWallet создает случайный объект Wallet для тестирования
func randomWallet() db.Wallet {
	return db.Wallet{
		WalletUuid: util.RandomUUID(),
		Balance:    util.RandomMoney(),
	}
}

// requireBodyMatchWallet проверяет, что тело ответа совпадает с ожидаемым объектом Wallet
func requireBodyMatchWallet(t *testing.T, body *bytes.Buffer, wallet db.Wallet) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	balance, err := strconv.ParseInt(string(data), 10, 64)
	require.NoError(t, err)
	require.Equal(t, wallet.Balance, balance)
}

// TestGetWalletAPI тестирует API для получения информации о кошельке
func TestGetWalletAPI(t *testing.T) {
	wallet := randomWallet()

	testCases := []struct {
		name       string
		WalletUuid uuid.UUID
		buildStubs func(store *mockdb.MockStore)
		//buildStubsRabbit func(rabbit *rabbitmq.RabbitMQ)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			WalletUuid: wallet.WalletUuid,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWallet(gomock.Any(), gomock.Eq(wallet.WalletUuid)).
					Times(1).
					Return(wallet, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchWallet(t, recorder.Body, wallet)
			},
		},
		{
			name:       "NotFound",
			WalletUuid: wallet.WalletUuid,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWallet(gomock.Any(), gomock.Eq(wallet.WalletUuid)).
					Times(1).
					Return(db.Wallet{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalError",
			WalletUuid: wallet.WalletUuid,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWallet(gomock.Any(), gomock.Eq(wallet.WalletUuid)).
					Times(1).
					Return(db.Wallet{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "InvalidUUID: UUID is nil",
			WalletUuid: uuid.Nil,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWallet(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			//mockChannel := mockrabbitmq.NewMockAMQPChannel(ctrl)
			//rabbitmq := rabbitmq.NewMockRabbitMQ(ctrl)
			//tc.buildStubsRabbit(rabbitmq)

			//rabbitmq, err := rabbitmq.NewRabbitMQ("amqp://guest:guest@localhos:5672/")
			//require.NoError(t, err)

			server := newTestServer(t, store) // , rabbitmq)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/wallets/%s", tc.WalletUuid)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

// TestGetWalletAPI2 тестирует API для получения информации о кошельке с некорректным UUID
func TestGetWalletAPI2(t *testing.T) {
	testCases := []struct {
		name       string
		WalletUuid string
		buildStubs func(store *mockdb.MockStore)
		//buildStubsRabbit func(rabbit *rabbitmq.RabbitMQ)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "InvalidUUID: wrong format",
			WalletUuid: "hahahah",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWallet(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/wallets/%s", tc.WalletUuid)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

// TestUpdateWalletBalanceAPI тестирует API для обновления баланса кошелька
func TestUpdateWalletBalanceAPI(t *testing.T) {
	wallet := randomWallet()

	testCases := []struct {
		name       string
		body       gin.H
		buildStubs func(store *mockdb.MockStore)
		//buildStubsRabbit func(rabbit *rabbitmq.RabbitMQ)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK (DEPOSIT)",
			body: gin.H{
				"wallet_uuid":    wallet.WalletUuid,
				"amount":         1000,
				"operation_type": "DEPOSIT",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.TransferTxResult{Balance: 2000}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OK (WITHDRAW)",
			body: gin.H{
				"wallet_uuid":    wallet.WalletUuid,
				"amount":         500,
				"operation_type": "WITHDRAW",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.TransferTxResult{Balance: 500}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Invalid UUID",
			body: gin.H{
				"wallet_uuid":    "not-a-valid-uuid",
				"amount":         1000,
				"operation_type": "WITHDRAW",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidUUID: UUID is nil",
			body: gin.H{
				"wallet_uuid":    uuid.Nil,
				"amount":         1000,
				"operation_type": "WITHDRAW",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWallet(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid amount",
			body: gin.H{
				"wallet_uuid":    wallet.WalletUuid,
				"amount":         -1000,
				"operation_type": "WITHDRAW",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid operation type",
			body: gin.H{
				"wallet_uuid":    wallet.WalletUuid,
				"amount":         1000,
				"operation_type": "hah",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Bad Request (Missing field)",
			body: gin.H{
				"wallet_uuid": wallet.WalletUuid,
				"amount":      1000,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Server Error (DB Failure)",
			body: gin.H{
				"wallet_uuid":    wallet.WalletUuid,
				"amount":         1000,
				"operation_type": "DEPOSIT",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.TransferTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(
			tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				store := mockdb.NewMockStore(ctrl)
				tc.buildStubs(store)

				server := newTestServer(t, store)
				recorder := httptest.NewRecorder()

				data, err := json.Marshal(tc.body)
				require.NoError(t, err)

				url := "/api/v1/wallet"
				request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(recorder)
			})
	}
}
