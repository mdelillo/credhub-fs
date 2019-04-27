package server_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/mdelillo/credhub-fs/test/helpers"
	"github.com/mdelillo/credhub-fs/test/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	var (
		listenAddr string
		certPath   string
		keyPath    string
	)

	BeforeEach(func() {
		listenAddr = helpers.GetFreeAddr()
		certPath, keyPath = helpers.GenerateSelfSignedCert("127.0.0.1")
	})

	AfterEach(func() {
		os.Remove(certPath)
		os.Remove(keyPath)
	})

	It("starts and shuts down a server with the given handlers", func() {
		handlerResponse := "response from test handler"
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, handlerResponse)
		})

		s := server.NewServer(listenAddr, certPath, keyPath, handler)

		serverDone := make(chan interface{})
		go func() {
			Expect(s.Start()).To(Succeed())
			close(serverDone)
		}()

		Expect(helpers.WaitForServerToBeAvailable(listenAddr, 5*time.Second)).To(Succeed())

		body := get("https://" + listenAddr)
		Expect(body).To(Equal(handlerResponse))

		Expect(s.Shutdown()).To(Succeed())
		Eventually(serverDone).Should(BeClosed(), "s.Start() did not return")

		Expect(helpers.ServerIsAvailable(listenAddr)).To(BeFalse())
	})
})

func get(url string) string {
	resp, err := helpers.HTTPClient.Get(url)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return string(body)
}
