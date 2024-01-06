-- DROP SCHEMA IF EXISTS feed;
CREATE SCHEMA IF NOT EXISTS feed;
USE feed;

CREATE TABLE IF NOT EXISTS feed_site (
    id 		    INT 			NOT NULL AUTO_INCREMENT,
	name 	    VARCHAR(255) 	NOT NULL,
    url 	    VARCHAR(255) 	NOT NULL,
    type 	    VARCHAR(8) 		NOT NULL,
    updated     DATETIME        NULL,
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;


CREATE TABLE IF NOT EXISTS feed_content (
	id 			    INT             NOT NULL AUTO_INCREMENT,
    feed_site_id 	INT             NOT NULL,
    content_id	    VARCHAR(256)    NULL,
    title		    VARCHAR(512),
    link		    VARCHAR(512),
    pub_date	    DATETIME,
    description		TEXT,
    content		    TEXT,
    authors		    VARCHAR(256),
    hash            VARCHAR(512),
    PRIMARY KEY     (id),
    FOREIGN KEY     (feed_site_id) REFERENCES feed_site(id) ON DELETE CASCADE,
    UNIQUE INDEX    (hash)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;
