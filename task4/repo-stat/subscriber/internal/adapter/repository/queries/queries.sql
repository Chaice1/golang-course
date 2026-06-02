-- name: CreateSubscription :exec

INSERT INTO subscriptions (repo_name,owner_name) VALUES($1,$2);

-- name: DeleteSubscription :exec

DELETE FROM subscriptions WHERE repo_name = $1 AND owner_name = $2;

-- name: GetSubscriptions :many

SELECT * FROM subscriptions;

-- name: GetSubscription :one

SELECT * FROM subscriptions WHERE repo_name = $1 AND owner_name = $2;

