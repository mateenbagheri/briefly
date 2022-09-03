USE briefly;

CREATE TABLE Users (
	userID INT NOT NULL AUTO_INCREMENT,
	name VARCHAR(50),
	familyName VARCHAR(50),
	password VARCHAR(100),
	salt VARCHAR(100),
	email VARCHAR(100),
	CONSTRAINT PK_Users PRIMARY KEY(userID),
	CONSTRAINT UC_UserEmail UNIQUE (email)
);

CREATE TABLE Collections (
	collectionID INT NOT NULL AUTO_INCREMENT,
	collectionName VARCHAR(50),
	userID INT,
	CONSTRAINT PK_Collections PRIMARY KEY(collectionID),
	CONSTRAINT FK_Collections_UserID FOREIGN KEY (userID)
		REFERENCES Users(userID)
		ON DELETE CASCADE
);

CREATE TABLE Links (
	linkID BIGINT NOT NULL AUTO_INCREMENT,
	link VARCHAR(600) NOT NULL,
	shortened VARCHAR(100) NOT NULL,
	createDate DATE NOT NULL,
	expDate DATE NOT NULL,
	CONSTRAINT PK_Links PRIMARY KEY(linkID),
	CONSTRAINT UC_MainLinks UNIQUE(link)
);

CREATE TABLE CollectionLinks (
	collectionLinkID BIGINT NOT NULL AUTO_INCREMENT,
    linkID BIGINT NOT NULL,
    collectionID INT NOT NULL,
    CONSTRAINT PK_CollectionLinks PRIMARY KEY(collectionLinkID),
    CONSTRAINT FK_CollectionLinks_collectionID FOREIGN KEY (collectionID)
		REFERENCES Collections(collectionID)
		ON DELETE CASCADE,
	CONSTRAINT FK_CollectionLinks_linkID FOREIGN KEY (linkID)
		REFERENCES Links(linkID)
        ON DELETE CASCADE
);

CREATE TABLE LinkHits (
	hitID BIGINT NOT NULL AUTO_INCREMENT,
	linkID BIGINT NOT NULL,
	hitDate DATE NOT NULL,
	CONSTRAINT PK_LinkHits PRIMARY KEY(hitID),
	CONSTRAINT FK_LinkHits_linkID FOREIGN KEY (linkID)
		REFERENCES Links(linkID)
		ON DELETE CASCADE
)
