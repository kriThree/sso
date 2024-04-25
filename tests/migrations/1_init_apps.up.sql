INSERT INTO apps (id, name, secret) 
VALUES (1,'test','secret1')
ON CONFLICT DO NOTHING; 