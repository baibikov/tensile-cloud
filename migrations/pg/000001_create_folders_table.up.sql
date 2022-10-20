create table if not exists files
(
    id        uuid         not null default uuid_in(md5(random()::text || random()::text)::cstring) primary key,
    folder_id uuid         not null,
    name      varchar(255) not null,
    type      varchar(50)  not null,
    format    varchar(50)  not null,

    constraint files_folder_id_folders_fkey foreign key (folder_id) references folders (id) on delete cascade on update cascade,
    constraint files_folder_id_name_type_format_uniq unique (folder_id, name, type, format)
);
