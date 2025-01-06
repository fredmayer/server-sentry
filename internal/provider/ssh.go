package provider

import (
	"fmt"
	"os"
	"strings"

	"github.com/fredmayer/sentry/internal/models"
	"github.com/fredmayer/sentry/internal/styles"
	"golang.org/x/crypto/ssh"
)

type Provider struct {
	client *ssh.Client
	server models.Server
}

func NewProvider(server models.Server) (*Provider, error) {
	var auth []ssh.AuthMethod

	if server.Password != "" {
		auth = append(auth, ssh.Password(server.Password))
	} else if server.Key != "" {
		key, err := os.ReadFile(server.Key)
		if err != nil {
			return nil, fmt.Errorf("unable to read private key: %v", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key: %v", err)
		}

		auth = append(auth, ssh.PublicKeys(signer))
	} else {
		return nil, fmt.Errorf("no authentication method provided for server %s", server.Name)
	}

	config := &ssh.ClientConfig{
		User:            server.User,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	address := fmt.Sprintf("%s:%d", server.Host, server.Port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server %s: %v", server.Name, err)
	}

	return &Provider{
		client: client,
		server: server,
	}, nil
}

func (p *Provider) DockerContainers() error {
	if !p.CommandExists("docker") {
		fmt.Println(styles.ReturnWithX("Docker is not installed"))
		return nil
	}

	session, err := p.client.NewSession()
	if err != nil {
		return fmt.Errorf("Error creating SSH session: %e", err)
	}
	defer session.Close()

	// Run `docker ps` to get container names
	output, err := session.CombinedOutput("docker ps --format \"{{.Names}}\"")
	if err != nil {
		return fmt.Errorf("Error checking Docker containers: %e", err)
	}

	containers := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(containers) == 1 && containers[0] == "" {
		fmt.Println(styles.ReturnWithX("No running Docker containers"))
		return nil
	}

	fmt.Println(styles.ReturnWithOk("Docker containers running:"))
	for _, container := range containers {
		fmt.Printf("- %s\n", container)
	}

	return nil
}

func (p *Provider) Pm2() (string, error) {
	session, err := p.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput("pm2 list")
	if err != nil {
		if p.IsCommandNotFound(err.Error()) {
			return styles.ReturnWithX("PM2 is not installed"), nil
		}
		return "", err
	}

	return fmt.Sprintf("%s\n%s", styles.ReturnWithOk("PM2"), string(output)), nil
}

func (p *Provider) NginxHosts() (string, error) {
	session, err := p.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("Error creating SSH session: %e", err)
	}
	defer session.Close()

	// Check if nginx is installed
	output, err := session.CombinedOutput("command -v nginx")
	if err != nil || strings.TrimSpace(string(output)) == "" {
		return styles.ReturnWithX("Nginx is not installed"), nil
	}

	// Run `nginx -T` to get configuration
	session, err = p.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("Error creating SSH session: %e", err)
	}
	defer session.Close()

	output, err = session.CombinedOutput("nginx -T")
	if err != nil {
		return "", fmt.Errorf("Error running nginx -T on %e", err)
	}

	// Extract server_name values
	config := string(output)
	lines := strings.Split(config, "\n")

	hosts := make([]models.NginxHost, 0)
	var serverName string
	var proxyPass string
	var serverBlock bool

	for _, line := range lines {

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "server_name") {
			// Extract server_name value
			parts := strings.Fields(line)
			if len(parts) > 1 {
				serverName = parts[1]
			}
		}

		if strings.HasPrefix(line, "proxy_pass") {
			// Extract proxy_pass value
			parts := strings.Fields(line)
			if len(parts) > 1 {
				proxyPass = parts[1]
			}
		}

		if strings.HasPrefix(line, "server {") {
			if serverBlock {
				hosts = append(hosts, models.NginxHost{
					ServerName: serverName,
					ProxyPass:  proxyPass,
				})

				serverName = ""
				proxyPass = ""
			}
			serverBlock = true
		}
	}

	for _, host := range hosts {
		fmt.Printf(" - %s -> %s \n", host.ServerName, host.ProxyPass)
	}

	return "", nil
}

func (p *Provider) IsCommandNotFound(response string) bool {
	strings.Contains(response, "Process exited with status 127")
	return true
}

func (p *Provider) CommandExists(command string) bool {
	session, err := p.client.NewSession()
	if err != nil {
		return false
	}
	defer session.Close()

	output, err := session.CombinedOutput(fmt.Sprintf("command -v %s", command))
	if err != nil || strings.TrimSpace(string(output)) == "" {
		return false
	}

	return true
}

func (p *Provider) Close() {
	p.client.Close()
}
