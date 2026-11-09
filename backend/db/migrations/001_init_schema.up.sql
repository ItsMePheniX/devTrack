CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 1. users table
CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255)UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'dev',
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 2.project table
CREATE TABLE IF NOT EXISTS projects(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100)NOT NULL,
    key VARCHAR(10)UNIQUE NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

--3. Project members table
CREATE TABLE IF NOT EXISTS project_members(
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'dev',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY(project_id, user_id)
);

--4. Tasks table
CREATE TABLE IF NOT EXISTS tasks(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id)ON DELETE CASCADE,
    task_number INT NOT NULL,
    title VARCHAR(255)NOT NULL,
    description TEXT,
    status VARCHAR(20)NOT NULL DEFAULT 'TODO',
    priority VARCHAR(20)NOT NULL DEFAULT 'MEDIUM',
    position DOUBLE PRECISION NOT NULL DEFAULT 65536.0,
    assignee_id UUID REFERENCES users(id)ON DELETE SET NULL,
    reporter_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    estimated_hours NUMERIC(6,2)DEFAULT 0.0,
    actual_hours NUMERIC(6,2)DEFAULT 0.0,
    due_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (project_id,task_number)
);

--5. Task comments
CREATE TABLE IF NOT EXISTS task_comments(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks(id)ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id)ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

--6. Task Dependencies Table(DAG's)
CREATE TABLE IF NOT EXISTS task_dependencies(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id)ON DELETE CASCADE,
    predecessor_id UUID NOT NULL REFERENCES tasks(id)ON DELETE CASCADE,
    successor_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    dependency_type VARCHAR(20)NOT NULL DEFAULT 'BLOCKS',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (predecessor_id, successor_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_projects_owner ON projects(owner_id);
CREATE INDEX IF NOT EXISTS idx_tasks_project_status ON tasks(project_id,status);
CREATE INDEX IF NOT EXISTS idx_tasks_assignee ON tasks(assignee_id);
CREATE INDEX IF NOT EXISTS idx_task_comments_task ON task_comments(task_id);
CREATE INDEX IF NOT EXISTS idx_task_dependencies_project ON task_dependencies(project_id);
