ALTER TABLE department
ADD COLUMN IF NOT EXISTS userId INT;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE table_name = 'department'
          AND constraint_name = 'fk_user'
    ) THEN
        ALTER TABLE department
        ADD CONSTRAINT fk_user
        FOREIGN KEY (userId)
        REFERENCES users(id)
        ON DELETE SET NULL;
    END IF;
END $$;