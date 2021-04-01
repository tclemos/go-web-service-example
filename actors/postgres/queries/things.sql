-- name: GetThing :one
SELECT * 
  FROM things t
WHERE t.id = $1
LIMIT 1;

-- name: GetThingByCode :one
SELECT * 
  FROM things t
WHERE t.code = $1
LIMIT 1;

-- name: ListThings :many
SELECT *
  FROM things t
 ORDER BY t.name;

-- name: CreateThing :one
INSERT INTO things (
  id, code, name
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateThing :one
UPDATE things
   SET code = $2
     , name = $3
 WHERE id = $1
 RETURNING *;
     
-- name: DeleteThing :exec
DELETE 
  FROM things
 WHERE id = $1;

-- name: DeleteThingByCode :exec
DELETE 
  FROM things
 WHERE code = $1;