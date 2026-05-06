# Fish Game Admin Dashboard & Gameplay Hub

## Giới thiệu
Đây là dự án Frontend dành cho Hệ thống Quản trị (CMS) và Sảnh Game (Gameplay Hub) của tựa game bắn cá 2D. 
Hệ thống cho phép:
- Quản lý thông tin người chơi, tài khoản, số vàng, theo dõi tỷ lệ thắng/thua.
- Cấu hình chỉ số cá 2D (dựa trên hệ số nhân Multiplier và tỷ lệ nổ cơ bản thay vì cơ chế HP truyền thống).
- Quản lý các phòng chơi, kiểm soát tỷ lệ trả thưởng (RTP) cho từng phòng.
- Quan sát và tham gia trực tiếp vào màn hình chơi game ngay trên cùng một hệ sinh thái.

*Lưu ý: Đây là dự án quản lý giao diện tĩnh (Frontend). Logic xử lý (Backend / WebSocket Server) hoạt động trên một hệ thống server riêng.*

## Yêu cầu hệ thống
- Node.js (phiên bản 18+ khuyên dùng, nếu chạy code trực tiếp)
- Docker (nếu chạy ứng dụng trên container)

---

## 🚀 Hướng dẫn chạy dự án

### Cách 1: Chạy trực tiếp trên máy (Môi trường Development)

1. Cài đặt các thư viện phụ thuộc:
   ```bash
   npm install
   ```

2. Khởi động server phát triển:
   ```bash
   npm run dev
   ```

3. Mở trình duyệt và truy cập vào đường dẫn hiển thị trên terminal (thường là `http://localhost:5173` hoặc `http://localhost:3000`).

---

### Cách 2: Chạy bằng công nghệ Docker (Mọi môi trường / Production)

Dự án đã được cấu hình sẵn `Dockerfile` để đóng gói thành một image chuẩn bị cho việc triển khai (deploy).

1. **Build Docker Image**
   Chạy lệnh sau tại thư mục gốc của dự án (nơi chứa file `Dockerfile`):
   ```bash
   docker build -t fish-game-admin:latest .
   ```

2. **Khởi chạy Docker Container**
   Sau khi build xong, khởi chạy container và ánh xạ với port `3000`:
   ```bash
   docker run -d -p 3000:3000 --name fish-admin-container fish-game-admin:latest
   ```

3. Mở trình duyệt và truy cập:
   👉 **http://localhost:3000**

*(Để dừng server Docker, chạy lệnh `docker stop fish-admin-container`. Để xóa, chạy lệnh `docker rm fish-admin-container`)*

## 🛠 Công nghệ sử dụng
- **React.js 19** & **Vite** (Build tool)
- **Tailwind v4** (Utility-First CSS)
- **Recharts** (Vẽ biểu đồ thống kê)
- **Lucide-react** (Hệ thống icon)
- **Docker** (Đóng gói và triển khai phần mềm)
