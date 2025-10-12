CREATE TABLE languages (
    _id int AUTO_INCREMENT,
    name varchar(255),
    creators varchar(255),
    extensions varchar(255),
    firstAppeared TIMESTAMP,
    year integer,
    wiki  varchar(255),
    PRIMARY KEY (_id)
);

INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("Assembly", "Kathleen Booth", ".asm,.s,.inc,.wla,.SRC", null, 1947, "https://en.wikipedia.org/wiki/Assembly_language");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("Bash", "Brian Fox", ".sh", "1989-06-08 00:00:00", 1989, "https://en.wikipedia.org/wiki/Bash_(Unix_shell)");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("C", "Dennis Ritchie", ".c,.h", null, 1972, "https://en.wikipedia.org/wiki/C_(programming_language)");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("C++", "Bjarne Stroustrup", ".C,.cc,.cpp,.cxx,.c++,.h,.H,.hh,.hpp,.hxx,.h++,.cppm,.ixx", null, 1985, "https://en.wikipedia.org/wiki/C%2B%2B");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("C#", "Anders Hejlsberg", ".cs,.csx", null, 2000, "https://en.wikipedia.org/wiki/C_Sharp_(programming_language)");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("COBOL", "Howard Bromberg,Norman Discount,Vernon Reeves,Jean E. Sammet,William Selden,Gertrude Tierney", ".cbl,.cob,.cpy", null, 1959, "https://en.wikipedia.org/wiki/COBOL");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("Elixir", "José Valim", ".ex,exs", null, 2012, "https://en.wikipedia.org/wiki/Elixir_(programming_language)");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("Fortran", "John Backus", ".f90,.f,.for", null, 1957, "https://en.wikipedia.org/wiki/Fortran");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("Golang", "Robert Griesemer,Rob Pike,Ken Thompson", ".go", "2009-11-10 00:00:00", 2009, "https://en.wikipedia.org/wiki/Go_(programming_language)");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("HTML", "Tim Berners-Lee", ".html,.htm", null, 1993, "https://en.wikipedia.org/wiki/HTML");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("Java", "James Gosling", ".java,.class,.jar,.jmod", "1995-05-23 00:00:00", 1995, "https://en.wikipedia.org/wiki/Java_(programming_language)");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("JavaScript", "Brendan Eich", ".js,.cjs,.mjs", "1995-12-04 00:00:00", 1995, "https://en.wikipedia.org/wiki/JavaScript");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("Perl", "Larry Wall", ".plx,.pl,.pm,xs,.t,.pod,.cgi", "1987-12-18 00:00:00", 1987, "https://en.wikipedia.org/wiki/Perl");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("PHP", "Rasmus Lerdorf", ".php,.phar,.phtml,.pht,.phps", "1995-06-08 00:00:00", 1995, "https://en.wikipedia.org/wiki/PHP");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("Python", "Guido van Rossum", ".py,.pyw,.pyz,.pyi,.pyc,.pyd", "1991-02-20 00:00:00", 1991, "https://en.wikipedia.org/wiki/Python_(programming_language)");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("Ruby", "Yukihiro Matsumoto", ".rb,.ru", null, 1995, "https://en.wikipedia.org/wiki/Ruby_(programming_language)");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("Rust", "Graydon Hoare", ".rs,.rlib", "2015-05-15 00:00:00", 2015, "https://en.wikipedia.org/wiki/Rust_(programming_language)");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("Scala", "Martin Odersky", ".scala,.sc", "2004-01-20 00:00:00", 2004, "https://en.wikipedia.org/wiki/Scala_(programming_language)");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("SQL", "Donald D. Chamberlin,Raymond F. Boyce", ".sql", null, 1974, "https://en.wikipedia.org/wiki/SQL");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("Swift", "Chris Lattner,Doug Gregor,John McCall,Ted Kremenek,Joe Groff", ".swift,.SWIFT", "2014-06-02 00:00:00", 2014, "https://en.wikipedia.org/wiki/Swift_(programming_language)");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("TypeScript", "Anders Hejlsberg", ".ts,.tsx,.mts,.cts", "2012-10-01 00:00:00", 2012, "https://en.wikipedia.org/wiki/TypeScript");
INSERT INTO languages(name, creators, extensions, firstAppeared, year, wiki) VALUES ("XML", "Tim Bray,Jean Paoli,Michael Sperberg-McQueen,Eve Maler,François Yergeau,John W. Cowan", ".xml", "1998-02-10 00:00:00", 1998, "https://en.wikipedia.org/wiki/XML");
