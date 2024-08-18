package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.uber.org/mock/gomock"

	"github.com/gin-gonic/gin"
	"github.com/zerotier/ztchooks"

	"zt-event-logger/mocks"
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

var _ = Describe("Router", func() {

	var (
		controller    *gomock.Controller
		mockDBClient  *mocks.MockDB
		mockProcessor *mocks.MockProcessor
		mockConfig    *config
		router        *gin.Engine
	)

	BeforeEach(func() {
		controller = gomock.NewController(GinkgoT())
		mockDBClient = mocks.NewMockDB(controller)
		mockProcessor = mocks.NewMockProcessor(controller)
		mockConfig = &config{dbFileLocation: "/tmp/zt.db"}
		router = ConfigureRouter(mockConfig, mockDBClient, mockProcessor)
	})

	AfterEach(func() {
		controller.Finish()
	})

	Context("POST /events/receive", func() {
		It("should receive and process the event successfully", func() {
			rawPayload := []byte(`{"hook_id":"abc123","org_id":"org456","hook_type":"NETWORK_JOIN","network_id":"net789","member_id":"mem012"}`)
			mockHookBase := &ztchooks.HookBase{HookID: "abc123", OrgID: "org456", HookType: "NETWORK_JOIN"}

			mockProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Return(mockHookBase, nil)

			req, _ := http.NewRequest("POST", "/events/receive", bytes.NewBuffer(rawPayload))

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Body.String()).To(ContainSubstring("Event received and logged successfully"))
		})
	})
})
