CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    github_url VARCHAR(255),
    live_url VARCHAR(255),
    technologies JSONB DEFAULT '[]'::jsonb,
    featured BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_projects_featured ON projects(featured);
CREATE INDEX idx_projects_created_at ON projects(created_at);

