package loadbalancer_test

import (
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/loadbalancer"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestRemoveTargetServer(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := loadbalancer.RemoveTargetCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(&hcloud.Server{ID: 321}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		RemoveServerTarget(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, &hcloud.Server{ID: 321}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"123", "--server", "my-server"})

	expOut := "Target removed from Load Balancer 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestRemoveTargetLabelSelector(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := loadbalancer.RemoveTargetCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		RemoveLabelSelectorTarget(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, "my-label").
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"123", "--label-selector", "my-label"})

	expOut := "Target removed from Load Balancer 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestRemoveTargetIP(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := loadbalancer.RemoveTargetCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		RemoveIPTarget(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, net.ParseIP("192.168.2.1")).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"123", "--ip", "192.168.2.1"})

	expOut := "Target removed from Load Balancer 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
