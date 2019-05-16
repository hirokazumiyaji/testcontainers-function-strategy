package wait

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"testing"

	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestWaitFunctionStrategy(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "nginx",
		ExposedPorts: []string{"80/tcp"},
		WaitingFor: NewFunctionStrategy(
			func(ctx context.Context, target wait.StrategyTarget) error {
				ip, err := target.Host(ctx)
				if err != nil {
					return err
				}

				port, err := target.MappedPort(ctx, "80")
				if err != nil {
					return err
				}
				address := net.JoinHostPort(ip, strconv.Itoa(port.Int()))

				url := fmt.Sprintf("http://%s", address)
				client := &http.Client{}
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return err
				}

				req = req.WithContext(ctx)
				res, err := client.Do(req)
				if err == nil {
					if http.StatusOK == res.StatusCode {
						return nil
					} else {
						return errors.New("status code is not ok")
					}
				} else {
					return err
				}
			},
			10,
		),
	}
	nginx, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer nginx.Terminate(ctx)
	ip, err := nginx.Host(ctx)
	if err != nil {
		t.Error(err)
	}
	port, err := nginx.MappedPort(ctx, "80")
	if err != nil {
		t.Error(err)
	}
	res, err := http.Get(fmt.Sprintf("http://%s:%s", ip, port.Port()))
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d. Got %d.", http.StatusOK, res.StatusCode)
	}
}
