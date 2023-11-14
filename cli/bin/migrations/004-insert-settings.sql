CREATE TABLE 
    `settings` (
    `id` integer not null primary key autoincrement,
    `created_at` datetime not null default CURRENT_TIMESTAMP,
    `value` varchar(255) null,
    `key` varchar(255) not null
)