package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/gin-gonic/gin"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, Gin Gonic!",
		})
	})

	router.GET("/ip", func(c *gin.Context) {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			fmt.Println("Error getting IP addresses:", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": addrs,
		})
	})

	router.GET("/cpu", func(c *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		c.JSON(http.StatusOK, gin.H{
			"cpu": runtime.NumCPU(),
			"mem": m.Alloc,
			"m":   m,
		})
	})

	router.GET("/env", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"envs": os.Environ(),
		})
	})

	router.GET("/hostname", func(c *gin.Context) {
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Println("Error getting hostname:", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"hostname": hostname,
		})
	})

	router.GET("/aux", func(c *gin.Context) {
		// Run the "ps" command
		cmd := exec.Command("ps", "aux")
		output, err := cmd.Output()

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"hostname": string(output),
		})
	})

	router.GET("/kconf", func(c *gin.Context) {
		config, err := rest.InClusterConfig()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"err": err,
			})

			return
		}

		clienset, err := kubernetes.NewForConfig(config)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"err": err,
			})

			return
		}

		// Get the current namespace from the configuration
		namespace, err := clienset.CoreV1().Namespaces().List(c, v1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"err": err,
			})

			return
		}

		var p []coreV1.Pod

		for _, v := range namespace.Items {
			pods, err := clienset.CoreV1().Pods(v.Name).List(c, v1.ListOptions{})
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"err": err,
				})

				return
			}

			for _, pd := range pods.Items {
				p = append(p, pd)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"namespace": namespace.Items,
			"pods":      p,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}

	fmt.Printf("Server is running on http://localhost:%s\n", port)

	err := router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
