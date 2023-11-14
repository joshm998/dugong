#! /usr/bin/env node 
const path = require('path');
// const dugong = require.resolve('dugong-cms').replace("dist/index.js", "")
const fs = require('fs')

const db = require('better-sqlite3')(path.join(process.cwd(), 'sqlite.db'));

let migrations = [];
console.log("Generating Table Schema...")

try {
    migrations = db.prepare(`select * from migrations`).all()
    
}
catch {
    console.log("Creating Migrations Table")
    db.exec(`create table
    'migrations' (
      'id' integer not null primary key autoincrement,
      'migrated_at' datetime not null default CURRENT_TIMESTAMP,
      'name' varchar(255) not null
    )`)
}

const dir = fs.opendirSync(path.join(__dirname, '/migrations'))
let dirent
while ((dirent = dir.readSync()) !== null) {
    if (dirent.isFile()) {
        const name = dirent.name.replace(".sql", "");
        if (migrations.filter(e => e.name == name).length == 0) {
            console.log(`Starting Migration: ${name}`)
            var file = fs.readFileSync(path.join(dirent.path, dirent.name), 'utf8')
            db.exec(file)
            db.exec(`insert into 'migrations' ('migrated_at', 'name') values (DATE('now'), '${name}')`)
        }
    }
}

console.log("Schema Created Successfully")