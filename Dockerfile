# Lớp Build (Build Stage)
FROM node:22-alpine AS build

WORKDIR /app

# Sao chép package.json và package-lock.json (nếu có)
COPY package*.json ./

# Cài đặt các thư viện (dependencies)
RUN npm install

# Sao chép toàn bộ mã nguồn vào container
COPY . .

# Chạy lệnh build ứng dụng React/Vite
RUN npm run build

# Lớp Production (Production Stage)
FROM node:22-alpine

WORKDIR /app

# Cài đặt server tĩnh để phục vụ file (serve)
RUN npm install -g serve

# Sao chép các file đã được build từ lớp build sang lớp production
COPY --from=build /app/dist ./dist

# Khai báo port 3000 (port mặc định của dự án)
EXPOSE 3000

# Chạy server phục vụ thư mục dist trên port 3000
CMD ["serve", "-s", "dist", "-l", "3000"]
