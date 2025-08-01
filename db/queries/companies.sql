-- name: CreateCompany :one
INSERT INTO companies (name, subscription_plan)
VALUES ($1, $2)
RETURNING id, name, subscription_plan, created_at, updated_at;

-- name: GetCompany :one
SELECT id, name, subscription_plan, created_at, updated_at
FROM companies
WHERE id = $1;

-- name: UpdateCompany :one
UPDATE companies
SET
    name = COALESCE($2, name),
    subscription_plan = COALESCE($3, subscription_plan),
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, subscription_plan, created_at, updated_at;

-- name: DeleteCompany :exec
DELETE FROM companies
WHERE id = $1;

-- name: ListCompanies :many
SELECT id, name, subscription_plan, created_at, updated_at
FROM companies
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountCompanies :one
SELECT COUNT(*) FROM companies;