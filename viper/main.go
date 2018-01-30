package main

import (
	"log"

	"time"

	"sync"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

// viper是Golang App的配置解决方案，可以读取多种来源的配置内容
// 使你不用担心配置来源于格式
// 支持配置方式如下:
// 	1. 默认设置配置
//	2. 从JSON，TOML，YAML，HCL 或者 Java Properties配置文件读取配置
//  3. 实时监控与重载配置更新
//  4. 从环境变量读取配置
//  5. 从etcd或consul读取配置
//  6. 从命令行读取配置
//  7. 从缓存读取配置
//
// 参数配置优先级是 直接设置 > flag > env > config > kv > 默认
// 对etcd consul的支持还未成熟
func main() {
	var appName = "myTest"

	// set default
	viper.Set("AppName", appName)

	// read from file
	viper.SetConfigName("app")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/" + appName)
	viper.AddConfigPath("/etc/" + appName)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	log.Println(viper.Get("AppName"))
	log.Println(viper.Get("Author"))
	log.Println(viper.Get("Version"))

	// read from remote maybe overwritten
	var v = viper.New()
	v.AddRemoteProvider("etcd", "http://localhost:2379", "/config/app.yaml")
	v.SetConfigType("yaml")
	err := v.ReadRemoteConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(v.Get("AppName"))
	log.Println(v.Get("Author"))
	log.Println(v.Get("Version"))

	go func() {
		for {
			time.Sleep(5 * time.Second)
			err := v.WatchRemoteConfig()
			v.WatchRemoteConfigOnChannel()
			if err != nil {
				log.Fatal(err)
			}
			log.Println(v.Get("AppName"))
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
