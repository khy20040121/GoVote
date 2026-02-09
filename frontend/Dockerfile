# ================== 1. 构建阶段 ==================
FROM node:18-alpine as builder

WORKDIR /app

# 优先拷贝 package 文件以利用缓存
COPY package.json package-lock.json ./

# 安装依赖 (设置 registry 加速)
RUN npm config set registry https://registry.npmmirror.com/ && npm install

# 拷贝源码并构建
COPY . .
# 这一步会生成 dist 目录
RUN npm run build

# ================== 2. 运行阶段 ==================
FROM nginx:alpine

# 拷贝构建产物到 Nginx 目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 拷贝自定义的 Nginx 配置 (用于转发 API 请求)
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
