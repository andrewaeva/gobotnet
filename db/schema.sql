drop table if exists users;
create table users(
    uuid varchar(32) primary key ,
    ip text ,
    name text,
    command text,
    command_param text,
    status boolean,
    time datetime,
    information text
    groupid text 
);

drop table if exists output;
create table output(
    command_uuid varchar(32) primary key,
    uuid varchar(32),
    command text,
    output text
);

drop table if exists screenshots;
create table screenshots(
    screen_uuid varchar(32) primary key,
    uuid varchar(32),
    screen text
);

drop table if exists download;
create table download(
    download_uuid varchar(32) primary key,
    uuid varchar(32),
    download_base64_filename text,
    download_base64_pathfile text,
    download_base64_data text
);

drop table if exists upload;
create table upload(
    uuid varchar (32),
    upload_base64_filename text,
    upload_base64_data text,
    foreign key(uuid) references users(uuid)
);