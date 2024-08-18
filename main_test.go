package main

// import (
// 	"testing"

// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"

// 	"github.com/gin-gonic/gin"

// 	"zt-event-logger/db"
// 	"zt-event-logger/mocks"
// )

// func TestMain(t *testing.T) {
// 	RegisterFailHandler(Fail)
// 	RunSpecs(t, "Main Suite")
// }

// var _ = Describe("Event Receiver", func() {
// 	var (
// 		mockOrgsAPI *mocks.MockOrganizationsAPI
// 		mocksSTSAPI *mocks.MockSTSAPI
// 		ctx         context.Context
// 		controller  *gomock.Controller
// 	)

// 	BeforeEach(func() {
// 		controller = gomock.NewController(GinkgoT())
// 		mockOrgsAPI = mocks.NewMockOrganizationsAPI(controller)
// 		mocksSTSAPI = mocks.NewMockSTSAPI(controller)
// 		ctx = context.Background()
// 	})

// 	AfterEach(func() {
// 		controller.Finish()
// 	})

// 	BeforeEach(func() {
// 		gin.SetMode(gin.TestMode)

// 		mockDBClient = &db.MockDBClient{}
// 		mockProcessor = &events.MockProcessor{}

// 		mockProcessor.ProcessFunc = func(payload []byte, opts ...events.SignatureOpt) (*db.HookBase, error) {
// 			return mockHookBase, nil
// 		}

// 		router = gin.Default()
// 		router.POST("/events/receive", func(c *gin.Context) {
// 			rawPayload, err := io.ReadAll(c.Request.Body)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 				return
// 			}

// 			signature := c.GetHeader("X-ZTC-Signature")
// 			psk := preSharedKey

// 			hookBase, err := mockProcessor.Process(rawPayload, events.WithSignatureInfo(signature, psk))
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 				return
// 			}

// 			c.JSON(http.StatusOK, gin.H{
// 				"message":   "Event received and logged successfully",
// 				"hook_id":   hookBase.HookID,
// 				"org_id":    hookBase.OrgID,
// 				"hook_type": hookBase.HookType,
// 			})
// 		})
// 	})

// 	Context("POST /events/receive", func() {
// 		It("should receive and process the event successfully", func() {
// 			req, _ := http.NewRequest("POST", "/events/receive", bytes.NewBuffer(rawPayload))
// 			req.Header.Set("X-ZTC-Signature", "test-signature")

// 			rr := httptest.NewRecorder()
// 			router.ServeHTTP(rr, req)

// 			Expect(rr.Code).To(Equal(http.StatusOK))
// 			Expect(rr.Body.String()).To(ContainSubstring("Event received and logged successfully"))
// 		})

// 		It("should return an error if the payload is invalid", func() {
// 			req, _ := http.NewRequest("POST", "/events/receive", bytes.NewBuffer([]byte("invalid payload")))
// 			req.Header.Set("X-ZTC-Signature", "test-signature")

// 			rr := httptest.NewRecorder()
// 			router.ServeHTTP(rr, req)

// 			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
// 		})

// 		It("should return an error if the processor fails", func() {
// 			mockProcessor.ProcessFunc = func(payload []byte, opts ...events.SignatureOpt) (*db.HookBase, error) {
// 				return nil, errors.New("processing error")
// 			}

// 			req, _ := http.NewRequest("POST", "/events/receive", bytes.NewBuffer(rawPayload))
// 			req.Header.Set("X-ZTC-Signature", "test-signature")

// 			rr := httptest.NewRecorder()
// 			router.ServeHTTP(rr, req)

// 			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
// 		})
// 	})
// })
