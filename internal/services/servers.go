package services

import (
	"fmt"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/charmbracelet/lipgloss"
	"github.com/fredmayer/sentry/internal/models"
	"github.com/fredmayer/sentry/internal/provider"
	"github.com/fredmayer/sentry/internal/styles"
)

type Servers struct {
	config *models.Config
}

func NewServers(config *models.Config) *Servers {
	return &Servers{
		config: config,
	}
}

func (s *Servers) Search(search string) error {
	search = strings.TrimSpace(search)

	for _, server := range s.config.Servers {
		p, err := provider.NewProvider(server)
		if err != nil {
			log.Error(fmt.Errorf("Connection err %s: %e", server.Name, err))
			continue
		}

		hosts, err := p.NginxHosts()
		if hosts != nil && len(hosts) > 0 {
			for _, host := range hosts {
				if search == host.ServerName {
					fmt.Println(styles.StatusBar("\nFinded:", fmt.Sprintf("%s - %s -> %s", server.Name, host.ServerName, host.ProxyPass), ""))
					return nil
				} else {
					// fmt.Printf("|%s != %s|\n", search, host.ServerName)
				}
			}
		}
	}

	fmt.Println(styles.StatusBar("Host not found on any servers", "", ""))
	return nil
}

// checkDNSDelegation - проверяет делегирован ли домен на сервер
func checkDNSDelegation(domain, serverIP string) bool {
	ips, err := net.LookupIP(domain)
	if err != nil {
		log.Printf("Error looking up IP for domain %s: %v", domain, err)
		return false
	}

	for _, ip := range ips {
		if ip.String() == serverIP {
			return true
		}
	}
	return false
}

func (s *Servers) Scan(selected string) error {
	// Рендерим красивый инфрормационный блок
	// TODO вынести в отдельную модель
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).Align(lipgloss.Center)

	style.Width(style.GetMaxWidth())

	r := style.Render(fmt.Sprintf("Load %d server configurations", len(s.config.Servers)))
	fmt.Println(r)

	// servers
	for _, server := range s.config.Servers {
		if len(selected) > 0 && selected != server.Name {
			continue
		}
		styles.H2(server.Name)

		p, err := provider.NewProvider(server)
		if err != nil {
			log.Error(fmt.Errorf("Connection err %s: %e", server.Name, err))
			continue
		}

		// 3. Nginx
		hosts, err := p.NginxHosts()
		if err != nil {
			log.Error(fmt.Errorf("Error running nginx command on %s: %v", server.Name, err))
		}
		if hosts != nil && len(hosts) > 0 {
			for _, host := range hosts {
				isDelegated := checkDNSDelegation(host.ServerName, server.Host)
				cross := lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("#FF6347"))
				ok := lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("#07b804"))

				status := cross.Render("×")
				if isDelegated {
					status = ok.Render("✓")
				}

				// fmt.Println(styles.StatusBar(host.ServerName, host.ProxyPass, status))
				fmt.Printf(" - %s (%s) -> %s  \n", host.ServerName, status, host.ProxyPass)
			}
		}

		// 2. PM2
		res, err := p.Pm2()
		if err != nil {
			log.Error(fmt.Errorf("Error running pm2 command on %s: %v", server.Name, err))
		}
		// TODO parse and vizualizations
		fmt.Println(res)

		// 1. Проверяем docker containers
		err = p.DockerContainers()
		if err != nil {
			log.Error(fmt.Errorf("Error running docker command on %s: %v", server.Name, err))
		}

		defer p.Close()

	}

	return nil
}
