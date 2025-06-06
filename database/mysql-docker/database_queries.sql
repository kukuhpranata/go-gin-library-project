-- ---
DROP TABLE IF EXISTS borrowings;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS users;
-- ---
-- Table: books
-- ---
CREATE TABLE books (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(13) UNIQUE NOT NULL,
    publication_year INT,
    quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- ---
-- Table: users
-- ---
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- ---
-- Table: borrowings
-- ---
CREATE TABLE borrowings (
    id INT PRIMARY KEY AUTO_INCREMENT,
    book_id INT NOT NULL,
    user_id INT NOT NULL,
    borrow_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    return_date TIMESTAMP NULL,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE RESTRICT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);
-- ---
-- Dummy Data for books
-- ---
INSERT INTO books (title, author, isbn, publication_year, quantity)
VALUES (
        'The Hitchhiker''s Guide to the Galaxy',
        'Douglas Adams',
        '9780345391803',
        1979,
        5
    ),
    (
        '1984',
        'George Orwell',
        '9780451524935',
        1949,
        7
    ),
    (
        'Pride and Prejudice',
        'Jane Austen',
        '9780141439518',
        1813,
        3
    ),
    (
        'Sapiens: A Brief History of Humankind',
        'Yuval Noah Harari',
        '9780062316097',
        2014,
        4
    ),
    (
        'To Kill a Mockingbird',
        'Harper Lee',
        '9780446310789',
        1960,
        6
    );