INSERT INTO forums (name, description)
VALUES (
    'Parent',
    'This is a parent forum. It will have children forums.'
  );
INSERT INTO forums (name, description, parent_id)
VALUES (
    'Child',
    'This is a child forum. It has a parent.',
    last_insert_rowid()
  );
INSERT INTO forums (name, description)
VALUES (
    'Another Parent',
    'This is a another parent forum. It will have children forums.'
  );
INSERT INTO forums (name, description, parent_id)
VALUES (
    'Another Child',
    'This is a another child forum. It has a parent.',
    last_insert_rowid()
  );