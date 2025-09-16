CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  first_name TEXT,
  last_name TEXT,
  email TEXT UNIQUE,
  password TEXT,
  usm_pesos INTEGER DEFAULT 0
);

CREATE TABLE books (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  book_name TEXT,
  book_category TEXT,
  transaction_type TEXT, -- 'Venta' o 'Arriendo'
  price INTEGER,
  status TEXT,           -- 'Disponible' / 'Agotado'
  popularity_score INTEGER DEFAULT 0
);

CREATE TABLE inventory (
  book_id INTEGER PRIMARY KEY,
  available_quantity INTEGER,
  FOREIGN KEY(book_id) REFERENCES books(id)
);

CREATE TABLE loans (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER,
  book_id INTEGER,
  start_date TEXT,   -- DD/MM/YYYY
  return_date TEXT,  -- fecha devuelto o NULL
  due_date TEXT,     -- fecha esperada de devolución
  status TEXT,       -- 'pendiente' / 'finalizado'
  FOREIGN KEY(user_id) REFERENCES users(id),
  FOREIGN KEY(book_id) REFERENCES books(id)
);

CREATE TABLE sales (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER,
  book_id INTEGER,
  sale_date TEXT,
  FOREIGN KEY(user_id) REFERENCES users(id),
  FOREIGN KEY(book_id) REFERENCES books(id)
);

CREATE TABLE transactions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER,
  book_id INTEGER,
  type TEXT, -- 'Compra' | 'Arriendo' | 'Abono' | 'Multa'
  date TEXT,
  amount INTEGER
);
