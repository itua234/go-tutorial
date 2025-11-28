package middlewares

import (
	models "confam-api/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "postgres",
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	return gormDB, mock, err
}

func TestAuthAppBySecretKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("missing API key", func(t *testing.T) {
		gormDb, _, _ := setupMockDB()
		w := httptest.NewRecorder() // Records the response
		c, _ := gin.CreateTestContext(w)

		// 2. Create Request without Header
		req, _ := http.NewRequest("GET", "/test", nil)
		c.Request = req

		// 3. Execute Middleware
		middleware := AuthAppBySecretKey(gormDb)
		middleware(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		//assert.Contains(t, w.Body.String(), "Missing API key")
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, true, response["error"])
		assert.Contains(t, response["message"], "Missing API key")
	})

	t.Run("Invalid API Key (Not found in DB)", func(t *testing.T) {
		gormDb, mock, _ := setupMockDB()
		// Expect a query to find the key, but return no rows (error)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "api_keys"`)).
			WithArgs("invalid-token").
			WillReturnError(gorm.ErrRecordNotFound)

		w := httptest.NewRecorder() // Records the response
		c, _ := gin.CreateTestContext(w)

		// 2. Create Request without Header
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("x-api-key", "invalid-token")
		c.Request = req

		// 3. Execute Middleware
		middleware := AuthAppBySecretKey(gormDb)
		middleware(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid API key")
	})

	t.Run("Valid API Key", func(t *testing.T) {
		gormDb, mock, _ := setupMockDB()
		mock.ExpectQuery("SELECT .* FROM `api_keys`").
			WithArgs("valid-token", 1).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "key", "app_id"}).
					AddRow(1, "valid-token", "app-uuid-123"),
			)

		mock.ExpectQuery("SELECT .* FROM `apps`").
			WithArgs("app-uuid-123").
			WillReturnRows(
				sqlmock.NewRows([]string{"ID", "name", "company_id"}).
					AddRow("app-uuid-123", "My App", "company-uuid-456"),
			)

		mock.ExpectQuery("SELECT .* FROM `companies`").
			WithArgs("company-uuid-456").
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name"}).
					AddRow("company-uuid-456", "My Company"),
			)

		w := httptest.NewRecorder() // Records the response
		c, _ := gin.CreateTestContext(w)

		// 2. Create Request without Header
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("x-api-key", "valid-token")
		c.Request = req

		// 3. Execute Middleware
		middleware := AuthAppBySecretKey(gormDb)
		middleware(c)

		assert.Equal(t, http.StatusOK, w.Code)

		app, exists := c.Get("app")
		assert.True(t, exists, "Context should contain 'app' key")

		appModel, ok := app.(*models.App)
		assert.True(t, ok, "Should be *models.App type")
		assert.NotNil(t, appModel)

		assert.Equal(t, "app-uuid-123", appModel.ID)
		assert.Equal(t, "My App", appModel.Name)
	})
}
