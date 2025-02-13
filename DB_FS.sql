-- Tạo cơ sở dữ liệu
CREATE DATABASE fashion_shop;
USE fashion_shop;

-- 1. Bảng Người Dùng (users)
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20) UNIQUE,
    avatar_url VARCHAR(500), -- Ảnh đại diện người dùng
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    role ENUM('admin', 'staff', 'customer') DEFAULT 'customer'
);

-- 2. Bảng Địa Chỉ (shipping_addresses)
CREATE TABLE shipping_addresses (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    full_address TEXT NOT NULL,
    city VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 3. Bảng Loại Sản Phẩm (categories)
CREATE TABLE categories (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL UNIQUE,
    image_url VARCHAR(500) -- Ảnh của loại sản phẩm
);

-- 4. Bảng Sản Phẩm (products)
CREATE TABLE products (
    id INT PRIMARY KEY AUTO_INCREMENT,
    category_id INT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    base_price DECIMAL(10,2) NOT NULL,
    status ENUM('Còn hàng', 'Hết hàng', 'Ngừng sản xuất') DEFAULT 'Còn hàng',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);

-- 5. Bảng Hình Ảnh Sản Phẩm (product_images)
CREATE TABLE product_images (
    id INT PRIMARY KEY AUTO_INCREMENT,
    product_id INT,
    image_url VARCHAR(500) NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- 6. Bảng Biến Thể Sản Phẩm (product_variants)
CREATE TABLE product_variants (
    id INT PRIMARY KEY AUTO_INCREMENT,
    product_id INT,
    size VARCHAR(50), -- Ví dụ: S, M, L, XL
    color VARCHAR(50), -- Ví dụ: Đỏ, Xanh, Đen
    stock INT DEFAULT 0, -- Số lượng tồn kho cho biến thể này
    price DECIMAL(10,2) NOT NULL, -- Giá có thể khác nhau theo biến thể
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- 7. Bảng Giỏ Hàng (cart)
CREATE TABLE cart (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    product_variant_id INT,
    quantity INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (product_variant_id) REFERENCES product_variants(id) ON DELETE CASCADE
);

-- 8. Bảng Đơn Hàng (orders)
CREATE TABLE orders (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    recipient_name VARCHAR(255) NOT NULL, -- Lưu thông tin người nhận để tránh thay đổi
    recipient_phone VARCHAR(20) NOT NULL,
    shipping_address_id INT, -- Tham chiếu địa chỉ giao hàng từ bảng `shipping_addresses`
    total_price DECIMAL(10,2) NOT NULL,
    status ENUM('Chờ xác nhận', 'Chuẩn bị hàng', 'Đang giao', 'Đã giao', 'Hoàn thành', 'Hủy') DEFAULT 'Chờ xác nhận',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (shipping_address_id) REFERENCES shipping_addresses(id) ON DELETE SET NULL
);

-- 9. Bảng Chi Tiết Đơn Hàng (order_details)
CREATE TABLE order_details (
    id INT PRIMARY KEY AUTO_INCREMENT,
    order_id INT,
    product_variant_id INT,
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL, -- Giá tại thời điểm đặt hàng
    total_price DECIMAL(10,2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_variant_id) REFERENCES product_variants(id) ON DELETE CASCADE
);

-- 10. Bảng Thanh Toán (payments)
CREATE TABLE payments (
    id INT PRIMARY KEY AUTO_INCREMENT,
    order_id INT,
    payment_method ENUM('Credit Card', 'Bank Transfer', 'COD', 'E-wallet') NOT NULL,
    payment_status ENUM('Chưa thanh toán', 'Đã thanh toán', 'Hoàn tiền') DEFAULT 'Chưa thanh toán',
    transaction_id VARCHAR(255) UNIQUE, -- Mã giao dịch từ cổng thanh toán
    amount DECIMAL(10,2) NOT NULL,
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);
