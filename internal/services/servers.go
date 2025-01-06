package services

import (
	"fmt"

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
		res, err := p.NginxHosts()
		if err != nil {
			log.Error(fmt.Errorf("Error running nginx command on %s: %v", server.Name, err))
		}
		// TODO parse and vizualizations
		fmt.Println(res)

		// 2. PM2
		res, err = p.Pm2()
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
