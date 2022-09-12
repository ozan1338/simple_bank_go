package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/ozan1338/db/mock"
	db "github.com/ozan1338/db/sqlc"
	"github.com/ozan1338/token"
	"github.com/ozan1338/util"
	"github.com/stretchr/testify/require"
)

// type getAccountRequest struct {
// 	ID int64 `uri:"id" binding:"required,min=1"`
// }

func TestGetAccountAPI(t *testing.T) {
	user,_ := randomUser(t)
	account := randomAccount(user.Username)

	testCase := []struct{
		name string
		accountID int64
		setupAuth func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationtypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t,recorder.Body,account)
			},
		},
		//TODO: add more case
		{
			name:      "UnauthorizedUser",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationtypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NotFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationtypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				// requireBodyMatchAccount(t,recorder.Body,account)
			},
		},
		{
			name: "InternalError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationtypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				// requireBodyMatchAccount(t,recorder.Body,account)
			},
		},
		{
			name: "BadRequest",
			accountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationtypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Any()).
				Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// requireBodyMatchAccount(t,recorder.Body,account)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			//build stubs
			tc.buildStubs(store)
			
			//start test server
			server := newTestServer(t,store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d",tc.accountID)
			request,err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}

	
	
}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID: util.RandomInt(1, 10000),
		Owner: owner,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount (t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t,err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)

	require.NoError(t,err)
	require.Equal(t, account, gotAccount)
}

type accountParams struct {
	Owner    string `json:"owner"`
	Currency  string  `json:"currency"`
	Balance int64 `json:balance`
}

func TestCreateAccount(t *testing.T) {
	user,_ := randomUser(t)
	account := randomAccount(user.Username)
	// accountTestArg := accountParams{
	// 	Owner: account.Owner,
	// 	Currency: account.Currency,
	// 	Balance: 0,
	// }

	testCase := []struct{
		name string
		body gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Invalid Currency",
			body: gin.H{
				"owner": account.Owner,
				"currency": "invalid",
				"balance": account.Balance,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationtypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				CreateAccount(gomock.Any(), gomock.Any()).
				Times(0)
				store.EXPECT().
				GetLastInsertId(gomock.Any()).
				Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// requireBodyMatchAccount(t,recorder.Body,account)
			},
		},
		{
			name: "User Doesnt Exist",
			body: gin.H{
				"owner": "Error",
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationtypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    "Error",
					Currency: account.Currency,
					Balance:  0,
				}
				store.EXPECT().UserExist(gomock.Any(), gomock.Eq(user.Username)).Times(1)
				store.EXPECT().
				CreateAccount(gomock.Any(), gomock.Eq(arg)).
				Times(0)
				// store.EXPECT().
				// GetLastInsertId(gomock.Any()).
				// Times(1).Return(account.ID, nil)
			},
			// checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			// 	require.Equal(t, http.StatusBadRequest, recorder.Code)
			// 	// requireBodyMatchAccount(t,recorder.Body,account)
			// },
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
	}

	for i := range testCase{
		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			// fmt.Println(tc)

			store := mockdb.NewMockStore(ctrl)
			//build stubs
			tc.buildStubs(store)
			
			//start test server
			server := newTestServer(t,store)
			recorder := httptest.NewRecorder()

			//Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/accounts")
			request,err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(recorder)
		})
	}
}

type testlistAccountRequest struct {
	page_size int32 
	page_id int32 
}

type tesListParamsss struct {
	Owner string
	Limit int32
	Offset int32
}

func TestListAccount(t *testing.T) {
	user,_ := randomUser(t)
	account := make([]db.Account,5)

	for i:=0 ; i < 5 ; i++ {
		account[i] = randomAccount(user.Username)
	}

	accounts := []db.Account{}

	testCase := []struct{
		name string
		page tesListParamsss
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "BadRequest",
			page: tesListParamsss{
				Owner: user.Username,
				Offset: 0,
				Limit: 1,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationtypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				
				store.EXPECT().ListAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "OK",
			page: tesListParamsss{
				Owner: user.Username,
				Limit: 5,
				Offset: 1,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationtypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// pageList := tesListParamsss{
				// 	Limit: 5,
				// 	Offset: 0,
				// }
				// fmt.Println("NEXT>>>",account)
				store.EXPECT().ListAccount(gomock.Any(), gomock.Any()).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalError",
			page: tesListParamsss{
				Owner: user.Username,
				Limit: 5,
				Offset: 1,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationtypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// pageList := tesListParamsss{
				// 	Limit: 5,
				// 	Offset: 0,
				// }
				// fmt.Println("NEXT>>>",account)
				store.EXPECT().ListAccount(gomock.Any(), gomock.Any()).Times(1).Return(accounts, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			page: tesListParamsss{
				Owner: user.Username,
				Limit: 5,
				Offset: 100,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationtypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// pageList := tesListParamsss{
				// 	Limit: 5,
				// 	Offset: 0,
				// }
				// fmt.Println("NEXT>>>",account)
				store.EXPECT().ListAccount(gomock.Any(), gomock.Any()).Times(1).Return(accounts, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			//build stubs
			tc.buildStubs(store)
			
			//start test server
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts?page_id=%d&page_size=%d",tc.page.Offset,tc.page.Limit)
			request,err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}