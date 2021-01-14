package nacos

import (
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/datasource/manager"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"net/url"
	"strconv"
)

const DataSourceNacos = "nacos"

func init() {

	manager.Register(DataSourceNacos, func() conf.DataSource {
		var (
			configAddr = flag.String("config")
		)
		if configAddr == "" {
			xlog.Panic("new nacos dataSource, configAddr is empty")
			return nil
		}
		// configAddr is a string in this format:
		// nacos://ip:port?dataId=XXXXX&group=DEFAULT_GROUP&tenant=XXXXX&scheme=http/https&level=debug|info|warn
		urlObj, err := url.Parse(configAddr)
		if err != nil {
			xlog.Panic("parse configAddr error", xlog.FieldErr(err))
			return nil
		}

		schema := urlObj.Query().Get("schema")
		if schema == "" {
			schema = "http"
		}

		var port uint64

		if schema == "http" {
			port = 80
		} else if schema == "https" {
			port = 443
		} else {
			xlog.Panic("unsupported schema: " + schema)
			return nil
		}

		if urlObj.Port() != "" {
			v, err := strconv.Atoi(urlObj.Port())
			if err == nil {
				port = uint64(v)
			}
		}

		sc := []constant.ServerConfig{
			{
				IpAddr: urlObj.Hostname(),
				Port:   port,
			},
		}

		cc := constant.ClientConfig{
			NamespaceId:         urlObj.Query().Get("tenant"), //namespace id
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "/tmp/nacos/log",
			CacheDir:            "/tmp/nacos/cache",
			RotateTime:          "1h",
			MaxAge:              3,
			LogLevel:            "debug",
		}

		level := urlObj.Query().Get("level")
		if level != "" {
			cc.LogLevel = level
		}

		return NewDataSource(sc, cc, urlObj.Query().Get("dataId"), urlObj.Query().Get("group"))
	})

}
