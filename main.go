package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/gin-gonic/gin"
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
			"mem":  m.Alloc,
			"m":  m,
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


	port := 8080

	fmt.Printf("Server is running on http://localhost:%d\n", port)

	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
