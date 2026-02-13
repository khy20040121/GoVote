
基于 Go 和 React 的社区投票项目。

## 技术栈

- **后端**: Go, Gin, MySQL, Redis
- **前端**: React, TypeScript, Vite, Tailwind CSS
- **部署**: Docker, Docker Compose, Nginx

## 功能特性

- 用户注册与登录 (JWT 认证)
- 帖子发布、查看详情
- 帖子列表 (支持按时间或热度排序)
- 帖子投票 (使用 Redis ZSet 实现排行榜)

## 快速开始 

使用 Docker Compose 可以一键启动完整环境（前端、后端、数据库）。

1. 确保已安装 Docker 和 Docker Compose。
2. 在项目根目录下运行：

   ```bash
   docker-compose up -d --build
   ```

3. 访问服务：
   - http://47.111.18.217


## 目录结构

- `controller/`: 处理路由请求
- `logic/`: 业务逻辑层
- `dao/`: 数据访问层 (MySQL/Redis)
- `models/`: 数据模型定义
- `frontend/`: 前端 React 项目
- `config/`: 配置文件
- `compose.yaml`: Docker 编排文件

## 文档

- 接口文档: 请参考 `API.md`
