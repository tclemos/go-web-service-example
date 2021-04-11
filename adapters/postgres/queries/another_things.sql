-- name: GetAnotherThing :one
SELECT * 
  FROM another_things t
WHERE t.id = $1
LIMIT 1;

-- name: GetAnotherThingByCode :one
SELECT * 
  FROM another_things t
WHERE t.code = $1
LIMIT 1;

-- name: ListAnotherThings :many
SELECT *
  FROM another_things t
 ORDER BY t.name;

-- name: CreateAnotherThing :one
INSERT INTO another_things (
  id, code, name
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateAnotherThing :one
UPDATE another_things
   SET code = $2
     , name = $3
 WHERE id = $1
 RETURNING *;
     
-- name: DeleteAnotherThing :exec
DELETE 
  FROM another_things
 WHERE id = $1;

-- name: DeleteAnotherThingByCode :exec
DELETE 
  FROM another_things
 WHERE code = $1;