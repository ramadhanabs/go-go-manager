CREATE TABLE IF NOT EXISTS employees (
    identity_number VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    gender VARCHAR(10) CHECK (gender IN ('male', 'female')) NOT NULL,
    department_id INTEGER NOT NULL, 
    employee_image_uri TEXT NOT NULL,
    FOREIGN KEY (department_id) REFERENCES department(id) ON DELETE CASCADE
);