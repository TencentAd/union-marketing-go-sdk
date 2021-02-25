# SDK 例子使用

#### 运行

```bash
go run sdk_example.go -config_path=config.json
```


#### 配置修改

以下配置需要自定义修改

| 配置项                 | 说明                                                    |
| ---------------------- | ------------------------------------------------------- |
| ams.auth.client_id     | 第三方应用ID， https://developers.e.qq.com/app 可以查看 |
| ams.auth.client_secret | 第三方应用密码                                          |
| ams.auth.redirect_uri  | 处理授权回调的地址                                      |



#### redirect_uri

sdk_example.go 可以处理授权回调，通过修改本地host，使回调请求发到sdk_example，就能测试授权回调处理功能。

