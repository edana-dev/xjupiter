package nacos

import (
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type nacosDataSource struct {
	client  config_client.IConfigClient
	dataId  string
	group   string
	changed chan struct{}
}

func NewDataSource(sc []constant.ServerConfig, cc constant.ClientConfig, dataId string, group string) conf.DataSource {

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})

	if err != nil {
		panic(err)
	}

	ds := &nacosDataSource{
		client:  client,
		dataId:  dataId,
		group:   group,
		changed: make(chan struct{}, 1),
	}

	err = ds.start()

	if err != nil {
		panic(err)
	}

	return ds
}

func (n *nacosDataSource) start() error {
	return n.client.ListenConfig(vo.ConfigParam{
		DataId: n.dataId,
		Group:  n.group,
		OnChange: func(namespace, group, dataId, data string) {
			xlog.Info("nacos config changed")
			xlog.Debugf("config changed group: %s, dataId: %s, content: \n%s", group, dataId, data)
			// notify config arch that config change, and will get config again?
			n.changed <- struct{}{}
		},
	})
}

func (n *nacosDataSource) stop() error {
	return n.client.CancelListenConfig(vo.ConfigParam{
		DataId: n.dataId,
		Group:  n.group,
	})
}

func (n *nacosDataSource) ReadConfig() ([]byte, error) {
	content, err := n.client.GetConfig(vo.ConfigParam{
		DataId: n.dataId,
		Group:  n.group,
	})
	xlog.Debug("GetConfig,config :" + content)
	if content == "" {
		xlog.Error("Failed to get config, please check.")
	}

	return []byte(content), err
}

func (n *nacosDataSource) IsConfigChanged() <-chan struct{} {
	return n.changed
}

func (n *nacosDataSource) Close() error {
	return n.stop()
}
