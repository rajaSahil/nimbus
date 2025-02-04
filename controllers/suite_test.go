// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of Nimbus

package controllers

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	intentv1 "github.com/5GSEC/nimbus/api/v1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config             // cfg will hold the Kubernetes rest configuration.
var k8sClient client.Client      // k8sClient is the Kubernetes client for the test environment.
var testEnv *envtest.Environment // testEnv is the environment for running the tests.

// TestControllers is the entry point for testing the controllers package.
func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail) // Register a Ginkgo fail handler.

	RunSpecs(t, "Controller Suite") // Run the Ginkgo specs for the 'Controller Suite'.
}

var _ = BeforeSuite(func() {
	// Setup for the test suite, executed before any specs are run.
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true))) // Set up the logger.

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		// CRDDirectoryPaths specifies the paths to the CRD manifests.
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,

		// The BinaryAssetsDirectory is only required if you want to run the tests directly
		// without call the makefile target test. If not informed it will look for the
		// default path defined in controller-runtime which is /usr/local/kubebuilder/.
		// Note that you must have the required binaries setup under the bin directory to perform
		// the tests directly. When we run make test it will be setup and used automatically.
		BinaryAssetsDirectory: filepath.Join("..", "..", "bin", "k8s",
			fmt.Sprintf("1.28.3-%s-%s", runtime.GOOS, runtime.GOARCH)),
	}

	var err error
	cfg, err = testEnv.Start()        // Start the test environment.
	Expect(err).NotTo(HaveOccurred()) // Assert that starting the environment does not produce an error.
	Expect(cfg).NotTo(BeNil())        // Assert that the config is not nil.

	// Add the intentv1 API scheme to the runtime scheme.
	err = intentv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred()) // Assert that adding the scheme does not produce an error.

	// Scaffold additional schemes here if needed.

	// Initialize the Kubernetes client for the test environment.
	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred()) // Assert that creating the client does not produce an error.
	Expect(k8sClient).NotTo(BeNil())  // Assert that the client is not nil.

})

var _ = AfterSuite(func() {
	// Teardown for the test suite, executed after all specs have run.
	By("tearing down the test environment")
	err := testEnv.Stop()             // Stop the test environment.
	Expect(err).NotTo(HaveOccurred()) // Assert that stopping the environment does not produce an error.
})
