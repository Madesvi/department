package handlers_test

import (
	"department/internal/api/handlers"
	"department/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"department/internal/api/handlers/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetDepartmentHandler(t *testing.T) {
	mockDept := models.Department{
		Employee: []models.Employee{{FullName: "Alex Jones"}},
		Children: []models.Department{{}},
	}

	tests := []struct {
		name           string
		urlPathID      string
		queryParams    string
		mockSetup      func(m *mocks.MockDepartmentGetter)
		expectedStatus int
		expectedBody   bool
	}{
		{
			name:      "Success case with standard defaults",
			urlPathID: "10",
			mockSetup: func(m *mocks.MockDepartmentGetter) {
				m.EXPECT().
					GetDepartment(mock.Anything, 10, 1, true).
					Return(mockDept, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   true,
		},
		{
			name:        "Success case with custom query configuration",
			urlPathID:   "20",
			queryParams: "?depth=2&include_employees=false",
			mockSetup: func(m *mocks.MockDepartmentGetter) {
				m.EXPECT().
					GetDepartment(mock.Anything, 20, 2, false).
					Return(mockDept, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   true,
		},
		{
			name:        "Fallback guard when depth limits are exceeded",
			urlPathID:   "30",
			queryParams: "?depth=99", // Handler max constraint caps this at 5
			mockSetup: func(m *mocks.MockDepartmentGetter) {
				m.EXPECT().
					GetDepartment(mock.Anything, 30, 5, true).
					Return(mockDept, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   true,
		},
		{
			name:           "Bad request when passing alphabetical ID format",
			urlPathID:      "xyz",
			mockSetup:      func(m *mocks.MockDepartmentGetter) {}, // Assert no core service call
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "Internal server error gracefully bubbles up",
			urlPathID: "50",
			mockSetup: func(m *mocks.MockDepartmentGetter) {
				m.EXPECT().
					GetDepartment(mock.Anything, 50, 1, true).
					Return(models.Department{}, errors.New("unexpected database connectivity drop")).
					Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGetter := mocks.NewMockDepartmentGetter(t)
			tt.mockSetup(mockGetter)

			handler := handlers.GetDepartmentHandler(mockGetter)

			reqUrl := fmt.Sprintf("/departments/%s%s", tt.urlPathID, tt.queryParams)
			req := httptest.NewRequest(http.MethodGet, reqUrl, nil)
			req.SetPathValue("id", tt.urlPathID) // Compatible with native Go 1.22+ routing mechanics

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody {
				assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

				var actualResponse struct {
					Department models.Department   `json:"department"`
					Employee   []models.Employee   `json:"employees"`
					Children   []models.Department `json:"children"`
				}

				err := json.NewDecoder(rr.Body).Decode(&actualResponse)
				assert.NoError(t, err)
				assert.Equal(t, mockDept.Employee, actualResponse.Employee)
				assert.Equal(t, mockDept.Children, actualResponse.Children)
			}
		})
	}
}
