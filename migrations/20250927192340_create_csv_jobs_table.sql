-- +goose Up
-- Create CSV jobs table for tracking CSV processing jobs
CREATE TABLE IF NOT EXISTS csv_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    filename VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    total_rows BIGINT DEFAULT 0,
    processed_rows BIGINT DEFAULT 0,
    failed_rows BIGINT DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_csv_jobs_user_id ON csv_jobs(user_id);
CREATE INDEX IF NOT EXISTS idx_csv_jobs_status ON csv_jobs(status);
CREATE INDEX IF NOT EXISTS idx_csv_jobs_created_at ON csv_jobs(created_at);
CREATE INDEX IF NOT EXISTS idx_csv_jobs_user_status ON csv_jobs(user_id, status);

-- Create a composite index for pagination queries
CREATE INDEX IF NOT EXISTS idx_csv_jobs_user_created_desc ON csv_jobs(user_id, created_at DESC);

-- +goose StatementBegin
-- Add trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_csv_jobs_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_csv_jobs_updated_at
    BEFORE UPDATE ON csv_jobs
    FOR EACH ROW
    EXECUTE FUNCTION update_csv_jobs_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS csv_jobs;
-- +goose StatementEnd