package pgdb

// TODO: возможно сделать, чтобы выражаения подготавливались перед выполнением запроса

const insertPost = `INSERT INTO posts (text, author_id, created_at) VALUES ($1, $2, current_timestamp) RETURNING id`

const updatePostById = `UPDATE posts SET text = $2 WHERE id = $1`

const selectAllPosts = `
SELECT p.id, p.text, p.created_at, p.updated_at, u.id as author_id, u.username as author_name, coalesce(u.avatar, '')
FROM posts p
INNER JOIN users u ON p.author_id = u.id
`

const selectPostById = `
SELECT p.id, p.text, p.created_at, p.updated_at, u.id as author_id, u.username as author_name, coalesce(u.avatar, '')
FROM posts p
INNER JOIN users u ON p.author_Id = u.id
`

const selectPostByAuthor = `
SELECT p.id, p.text, p.created_at, p.updated_at, u.id as author_id, u.username as author_name, coalesce(u.avatar, '')
FROM posts p
INNER JOIN users u ON p.author_Id = u.id
WHERE u.username = $1
ORDER BY p.created_at
`

const deletePostById = `DELETE FROM posts WHERE id = $1`

const selectCommentsByPost = `
SELECT c.id, c.text, c.created_at, c.updated_at, u.id as author_id, u.username as author_name, coalesce(u.avatar, '')
FROM comments c
INNER JOIN users u ON c.author_id = u.id
INNER JOIN posts p ON c.post_id = p.id
WHERE c.post_id = $1
`
