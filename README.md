# backend-lib

用于存放公用代码，如lib、helper等


## 如何引用gitlab私有库

### 1. 获取 gitlab 的AccessToken(已配置gitlab的ssh可以直接跳到步骤3)
进入 gitlab 个人 Settings —> Access Tokens，然后创建一个personal access token，这里权限最好选择只读 (read_repository)。

### 2. 配置本地Git
将下面命令中的 <YOUR_PRIVATE_TOKEN> 替换为上一步创建的 token.
```
git config --global http.extraheader "PRIVATE-TOKEN: <YOUR_PRIVATE_TOKEN>"
```

### 3. 将 git 请求从 https 转为 ssh(未在gitlab配置ssh的话可以跳过该步骤)
```
git config --global url."git@git-plat.tecorigin.net:".insteadof "http://git-plat.tecorigin.net/"
```
### 4. 修改go env
```
go env -w GOPRIVATE=git-plat.tecorigin.net
go env -w GONOSUMDB=git-plat.tecorigin.net
go env -w GONOPROXY=git-plat.tecorigin.net
go env -w GOINSECURE=git-plat.tecorigin.net
```

