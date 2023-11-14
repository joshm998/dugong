CREATE TABLE
    `fields` (
    `id` integer not null primary key autoincrement,
    `created_at` datetime not null default CURRENT_TIMESTAMP,
    `value` varchar(255) null,
    `name` varchar(255) not null,
    `field_type` varchar(50) not null
)