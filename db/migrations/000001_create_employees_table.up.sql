CREATE TABLE employees (
    identity_number VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    gender VARCHAR(10) CHECK (gender IN ('male', 'female')) NOT NULL,
    department_id VARCHAR(50) NOT NULL,
    employee_image_uri TEXT NOT NULL
);