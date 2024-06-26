CREATE TABLE IF NOT EXISTS items (
  id VARCHAR(255) PRIMARY KEY,
  access_token VARCHAR(255),

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS accounts (
  id VARCHAR(255) PRIMARY KEY,
  item_id VARCHAR(255),
  account_name VARCHAR(255) NOT NULL,
  account_type VARCHAR(255) NOT NULL,
  balance FLOAT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
