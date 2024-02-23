-- 删除用户名重复的用户，只保留一条
DELETE FROM chatgpt_users
WHERE username IN (
    SELECT username
    FROM (
             SELECT username
             FROM chatgpt_users
             GROUP BY username
             HAVING COUNT(*) > 1
         ) AS temp
) AND id NOT IN (
    SELECT MIN(id)
    FROM (
             SELECT id, username
             FROM chatgpt_users
             GROUP BY id, username
             HAVING COUNT(*) > 1
         ) AS temp
    GROUP BY username
);

-- 给 username 字段建立唯一索引
ALTER TABLE `chatgpt_users` ADD UNIQUE(`username`)