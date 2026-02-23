package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/floodguard/internal/config"
	"github.com/yourusername/floodguard/internal/monitor"
	"go.uber.org/zap"
)

// version is set at build time via -ldflags "-X main.version=x.y.z".
var version = "v1.0.3"

var (
	cfgFile string
	logger  *zap.Logger
)

func main() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	if err := rootCmd.Execute(); err != nil {
		logger.Error("Command execution failed", zap.Error(err))
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:     "floodguard",
	Version: version,
	Short:   "A lightweight firewall tool to protect against CC and DDoS attacks",
	Long:    `FloodGuard monitors network connections and automatically blocks malicious IPs.`,
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the firewall protection",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		mon := monitor.New(cfg, logger)
		return mon.Start()
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current protection status",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("FloodGuard Status: Running")
		fmt.Println("Monitoring connections...")
		return nil
	},
}

var (
	publicIP string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Auto-detect public IP if not provided
		if publicIP == "" {
			publicIP = autoDetectIPs()
		}
		return config.InitDefault(publicIP)
	},
}

func autoDetectIPs() string {
	var ips []string
	
	// Try to get public IP from common services
	publicIP := getPublicIP()
	if publicIP != "" {
		ips = append(ips, publicIP)
	}
	
	// Get local network IPs
	localIPs := getLocalIPs()
	ips = append(ips, localIPs...)
	
	return strings.Join(ips, " ")
}

func getPublicIP() string {
	// Try multiple services
	services := []string{
		"https://api.ipify.org",
		"https://ifconfig.me",
		"https://icanhazip.com",
	}
	
	client := &http.Client{Timeout: 5 * time.Second}
	
	for _, service := range services {
		resp, err := client.Get(service)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		
		ip := strings.TrimSpace(string(body))
		if ip != "" {
			return ip
		}
	}
	
	return ""
}

func getLocalIPs() []string {
	var ips []string
	
	ifaces, err := net.Interfaces()
	if err != nil {
		return ips
	}
	
	for _, iface := range ifaces {
		// Skip loopback and down interfaces
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}
		
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			
			if ip == nil || ip.IsLoopback() {
				continue
			}
			
			// Add both IPv4 and IPv6
			ips = append(ips, ip.String())
		}
	}
	
	return ips
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "/etc/floodguard/config.yaml", "config file path")
	
	initCmd.Flags().StringVar(&publicIP, "public-ip", "", "public IP address to add to whitelist")
	
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(initCmd)
}
