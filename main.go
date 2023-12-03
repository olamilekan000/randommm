package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	router.GET("get-roles-for-sa", func(c *gin.Context) {
		// Getname(c.Request.URL.String())

		config, err := rest.InClusterConfig()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"err": err,
			})

			return
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"err": err,
			})

			return
		}

		// serviceAccountName := "your-service-account-name"
		namespace := "default"

		roles, err := clientset.RbacV1().Roles(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		roleBindings, err := clientset.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"roles":       roles,
			"roleBinding": roleBindings,
			// "roleBinding": roleBinding,
		})

	})

	router.GET("/create-sa", func(c *gin.Context) {
		config, err := rest.InClusterConfig()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"err": err,
			})

			return
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"err": err,
			})

			return
		}

		// Define ServiceAccount
		serviceAccount := &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pipeops-63484a87",
				Namespace: "default",
			},
		}

		// Create ServiceAccount
		createdServiceAccount, err := clientset.CoreV1().ServiceAccounts("default").Create(context.TODO(), serviceAccount, metav1.CreateOptions{})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"err": err,
			})

			return
		}
		fmt.Printf("ServiceAccount created: %s\n", createdServiceAccount.Name)

		// Define Role
		role := &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "example-role",
				Namespace: "default",
			},
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{""},
					Resources: []string{"*"},
					Verbs:     []string{"*"},
				},
			},
		}

		// Create Role
		createdRole, err := clientset.RbacV1().Roles("default").Create(context.TODO(), role, metav1.CreateOptions{})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"err": err,
			})

			return
		}
		fmt.Printf("Role created: %s\n", createdRole.Name)

		// Define RoleBinding
		roleBinding := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "example-rolebinding",
				Namespace: "default",
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      "pipeops-63484a87",
					Namespace: "default",
				},
			},
			RoleRef: rbacv1.RoleRef{
				Kind:     "Role",
				Name:     "example-role",
				APIGroup: "rbac.authorization.k8s.io",
			},
		}

		// Create RoleBinding
		createdRoleBinding, err := clientset.RbacV1().RoleBindings("default").Create(context.TODO(), roleBinding, metav1.CreateOptions{})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"err": err,
			})

			return
		}

		fmt.Printf("RoleBinding created: %s\n", createdRoleBinding.Name)

		c.JSON(http.StatusOK, gin.H{
			"createdRoleBinding": createdRoleBinding,
			"roleBinding":        roleBinding,
			// "roleBinding": roleBinding,
		})
	})

	router.GET("/ns", func(c *gin.Context) {
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
		namespace, err := clienset.CoreV1().Namespaces().List(c, metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"err": err,
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"namespace": namespace.Items,
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

		// // Get the current namespace from the configuration
		// namespace, err := clienset.CoreV1().Namespaces().List(c, metav1.ListOptions{})
		// if err != nil {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"err": err,
		// 	})

		// 	return
		// }

		// var p []coreV1.Pod

		// for _, v := range namespace.Items {
		// 	pods, err := clienset.CoreV1().Pods(v.Name).List(c, metav1.ListOptions{})
		// 	if err != nil {
		// 		c.JSON(http.StatusOK, gin.H{
		// 			"err": err,
		// 		})

		// 		return
		// 	}

		// 	for _, pd := range pods.Items {
		// 		p = append(p, pd)
		// 	}
		// }

		pods, err := clienset.CoreV1().Pods("default").List(c, metav1.ListOptions{})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"err": err,
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"pods": pods,
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

func Getname(urlStr string) {
	subdomain, err := extractSubdomain(urlStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Subdomain:", subdomain)
}

func extractSubdomain(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	hostParts := strings.Split(parsedURL.Hostname(), ".")
	if len(hostParts) < 2 {
		return "", fmt.Errorf("invalid URL format")
	}

	return hostParts[0], nil
}
