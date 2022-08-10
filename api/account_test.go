package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/ozan1338/db/mock"
	db "github.com/ozan1338/db/sqlc"
	"github.com/ozan1338/util"
	"github.com/stretchr/testify/require"
)

// type getAccountRequest struct {
// 	ID int64 `uri:"id" binding:"required,min=1"`
// }

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCase := []struct{
		name string
		accountID int64
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			accountID: account.ID,
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
			name: "NotFound",
			accountID: account.ID,
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
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d",tc.accountID)
			request,err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}

	
	
}

func randomAccount() db.Account {
	return db.Account{
		ID: util.RandomInt(1, 10000),
		Owner: util.RandomOwner(),
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

	account := randomAccount()
	// accountTestArg := accountParams{
	// 	Owner: account.Owner,
	// 	Currency: account.Currency,
	// 	Balance: 0,
	// }

	testCase := []struct{
		name string
		account accountParams
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "BadRequest",
			account: accountParams{
				Owner: account.Owner,
				Currency: "IDR",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				CreateAccount(gomock.Any(), gomock.Any()).
				Times(0)
				store.EXPECT().
				GetLastInsertId(gomock.Any()).
				Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// requireBodyMatchAccount(t,recorder.Body,account)
			},
		},
	}

	for i := range testCase{
		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			fmt.Println(tc)

			//setData
			data := url.Values{}
			data.Set("owner", tc.account.Owner)
			data.Set("currency", tc.account.Currency)

			store := mockdb.NewMockStore(ctrl)
			//build stubs
			tc.buildStubs(store)
			
			//start test server
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts")
			request,err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

type testlistAccountRequest struct {
	page_size int32 
	page_id int32 
}

type tesListParamsss struct {
	Limit int32
	Offset int32
}

func TestListAccount(t *testing.T) {
	account := make([]db.Account,5)
	for i:=0 ; i < 5 ; i++ {
		account[i] = randomAccount()
	}

	accounts := []db.Account{}

	testCase := []struct{
		name string
		page tesListParamsss
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "BadRequest",
			page: tesListParamsss{
				Offset: 0,
				Limit: 1,
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
				Limit: 5,
				Offset: 1,
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
				Limit: 5,
				Offset: 1,
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
				Limit: 5,
				Offset: 100,
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
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts?page_id=%d&page_size=%d",tc.page.Offset,tc.page.Limit)
			request,err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}