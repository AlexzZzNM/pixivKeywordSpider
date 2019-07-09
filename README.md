# p站(pixiv)关键字爬虫

## 项目结构

```bash
|-config/                    #配置文件
|-public/                    #下载图片相关
|-useType/                   #爬取 和 存储所需要的数据类型
|-utils/                     #工具函数
| main.go                    #主函数
```

## 运行项目

```bash
# 克隆项目代码到本地
git clone git@github.com:yangfanjie97/pixivKwywordSpider.git

# 进入项目目录
cd pixivKwywordSpider

#在config目录配置您的账号和密码 以及代理

# 关于版本控制使用的为 go v1.11以上的自带的 mod 具体方法请自行百度

# 运行
go run main.go

备注： 在项目运行后，您可自行决定要搜索的关键字 以及决定是否下载图片
```

```bash
本项目仅用于学术交流！
```
