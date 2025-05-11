create table notes (
	-- 37 is the length of the uuid without the .md at the end
	filename char(37) primary key,
	title text
);
create table tags (
	tagid integer primary key,
	tagtext text
);
create table noteTagRelation (
	noteFilename char(37) NOT NULL,
	tagid integer NOT NULL,
	PRIMARY KEY (noteFilename, tagId)
);
-- full text search for notes and triggers
create virtual table notes_fts using fts5 (
	filename,
	title
);
create trigger insert_notes_fts
	after insert on notes
begin
	insert into notes_fts (filename, title)
	values (NEW.filename, NEW.title)
end;
create trigger update_notes_fts
	after update on notes
begin
	update notes_fts
	set
		filename = NEW.filename,
		title = NEW.title
	where filename = NEW.filename;
end;
create trigger delete_notes_fts
	after delete on notes
begin
	delete from notes_fts
	where filename = OLD.filename;
end;
-- full text search for tags and triggers
create virtual table tags_fts using fts5 (
	tagid,
 	tagtext
);
create trigger insert_tags_fts
	after insert on tags
begin
	insert into tags_fts (tagid, tagtext)
	values (NEW.tagid, NEW.tagtext)
end;
create trigger update_tags_fts
	after update on tags
begin
	update tags_fts
	set
		tagid = NEW.tagid,
		tagtext = NEW.tagtext
	where tagid = NEW.tagid;
end;
create trigger delete_tags_fts
	after delete on tags
begin
	delete from tags_fts
	where tagid = OLD.tagid;
end
