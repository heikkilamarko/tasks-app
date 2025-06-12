package shared

import (
	"crypto/tls"
	"log/slog"
	"os"

	"github.com/nats-io/nats.go"
)

func NewNATSConn(config *Config, logger *slog.Logger) (*nats.Conn, error) {
	conn, err := nats.Connect(
		config.Shared.NATSURL,
		nats.Secure(&tls.Config{InsecureSkipVerify: true}), // TODO: Revisit this. Should NOT be used in production.
		nats.UserCredentials(config.Shared.NATSCreds),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(
			func(_ *nats.Conn, err error) {
				logger.Info("nats disconnected", "reason", err)
			}),
		nats.ReconnectHandler(
			func(c *nats.Conn) {
				logger.Info("nats reconnected", "address", c.ConnectedUrl())
			}),
		nats.ErrorHandler(
			func(_ *nats.Conn, _ *nats.Subscription, err error) {
				logger.Error("nats error", "error", err)
				os.Exit(1)
			}),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
